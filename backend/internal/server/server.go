package server

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"pwmanager/internal/auth"
	"pwmanager/internal/config"
	"pwmanager/internal/crypto"
	"pwmanager/internal/handler"
	"pwmanager/internal/repository"
	"pwmanager/internal/service"
)

//go:embed swagger.json
var swaggerSpec []byte

const swaggerHTML = `<!DOCTYPE html>
<html><head><title>API</title><meta charset="utf-8"/>
<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"/>
</head><body><div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script>window.onload=()=>{SwaggerUIBundle({url:'/swagger/doc.json',dom_id:'#swagger-ui',persistAuthorization:true})}</script>
</body></html>`

type Server struct {
	cfg    *config.Config
	db     *pgxpool.Pool
	router http.Handler
}

func New(cfg *config.Config, db *pgxpool.Pool) (*Server, error) {
	r, err := buildRouter(cfg, db)
	if err != nil {
		return nil, err
	}
	return &Server{cfg: cfg, db: db, router: r}, nil
}

func (s *Server) Handler() http.Handler { return s.router }

func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  s.cfg.ReadTimeout(),
		WriteTimeout: s.cfg.WriteTimeout(),
	}
	errCh := make(chan error, 1)
	go func() {
		log.Printf("listening on %s  |  swagger: http://localhost:%d/swagger/", addr, s.cfg.Server.Port)
		errCh <- srv.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

func buildRouter(cfg *config.Config, db *pgxpool.Pool) (http.Handler, error) {
	cipher, err := crypto.New(cfg.Server.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("crypto: %w", err)
	}
	jwtMgr := auth.NewJWTManager(cfg.Server.JWTSecret, cfg.JWTTTL())

	userRepo := repository.NewUserRepository(db)
	passwordRepo := repository.NewPasswordRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)

	authSvc := service.NewAuthService(userRepo, jwtMgr)
	passwordSvc := service.NewPasswordService(passwordRepo, cipher)
	groupSvc := service.NewGroupService(groupRepo)
	apiKeySvc := service.NewAPIKeyService(apiKeyRepo)

	authH := handler.NewAuthHandler(authSvc)
	pwH := handler.NewPasswordHandler(passwordSvc)
	grpH := handler.NewGroupHandler(groupSvc)
	akH := handler.NewAPIKeyHandler(apiKeySvc)

	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-API-Key"},
		AllowCredentials: false, MaxAge: 300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK); _, _ = w.Write([]byte("ok"))
	})
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})
	r.Get("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(swaggerHTML))
	})
	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(swaggerSpec)
	})

	register := func(prefix string) {
		r.Route(prefix, func(r chi.Router) {
			r.Post("/auth/login", authH.Login)
			r.Post("/auth/register", authH.Register)
			r.Group(func(r chi.Router) {
				r.Use(handler.JWTAuthMiddleware(jwtMgr))
				r.Get("/auth/api-keys", akH.List)
				r.Post("/auth/api-keys", akH.Create)
				r.Delete("/auth/api-keys/{id}", akH.Delete)
			})
			r.Group(func(r chi.Router) {
				r.Use(handler.AuthMiddleware(jwtMgr, apiKeySvc))
				r.Get("/passwords", pwH.List)
				r.Post("/passwords", pwH.Create)
				r.Get("/passwords/{id}", pwH.Get)
				r.Put("/passwords/{id}", pwH.Update)
				r.Delete("/passwords/{id}", pwH.Delete)
				r.Get("/groups", grpH.List)
				r.Post("/groups", grpH.Create)
				r.Put("/groups/{id}", grpH.Update)
				r.Delete("/groups/{id}", grpH.Delete)
			})
		})
	}
	register("/api")
	register("/api/v1")
	return r, nil
}
