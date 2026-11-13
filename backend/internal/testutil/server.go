package testutil

import (
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"pwmanager/internal/auth"
	"pwmanager/internal/config"
	"pwmanager/internal/server"
)

type TestServer struct {
	Pool     *pgxpool.Pool
	Cfg      *config.Config
	UserID   int64
	Username string
	Password string
	Router   http.Handler
	JWT      *auth.JWTManager
	Token    string
}

func NewTestServer(t *testing.T) *TestServer {
	t.Helper()
	pool, cfg := NewTestDB(t)
	userID := FindUserID(t, pool, "admin")

	jwtMgr := auth.NewJWTManager(cfg.Server.JWTSecret, cfg.JWTTTL())
	token, _, err := jwtMgr.Generate(userID, "admin")
	if err != nil { t.Fatalf("jwt: %v", err) }

	srv, err := server.New(cfg, pool)
	if err != nil { t.Fatalf("server: %v", err) }

	return &TestServer{
		Pool: pool, Cfg: cfg, UserID: userID,
		Username: "admin", Password: "admin123",
		Router: srv.Handler(), JWT: jwtMgr, Token: token,
	}
}

func (ts *TestServer) SetAuth(r *http.Request) *http.Request {
	r.Header.Set("Authorization", "Bearer "+ts.Token)
	return r
}
