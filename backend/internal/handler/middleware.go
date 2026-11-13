package handler

import (
	"context"
	"net/http"
	"strings"

	"pwmanager/internal/auth"
	"pwmanager/internal/service"
)

type ctxKey string

const (
	UserIDKey   ctxKey = "user_id"
	UsernameKey ctxKey = "username"
)

// JWTAuthMiddleware — только Bearer JWT (для управления API-ключами)
func JWTAuthMiddleware(jwtMgr *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if h == "" || !strings.HasPrefix(h, "Bearer ") {
				writeJSONError(w, http.StatusUnauthorized, "missing bearer token")
				return
			}
			claims, err := jwtMgr.Parse(strings.TrimPrefix(h, "Bearer "))
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, "invalid token")
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UsernameKey, claims.Username)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AuthMiddleware — JWT или X-API-Key
func AuthMiddleware(jwtMgr *auth.JWTManager, apiKeys *service.APIKeyService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if key := r.Header.Get("X-API-Key"); key != "" {
				uid, uname, err := apiKeys.Validate(r.Context(), key)
				if err != nil {
					writeJSONError(w, http.StatusUnauthorized, "invalid api key")
					return
				}
				ctx := context.WithValue(r.Context(), UserIDKey, uid)
				ctx = context.WithValue(ctx, UsernameKey, uname)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			h := r.Header.Get("Authorization")
			if h == "" || !strings.HasPrefix(h, "Bearer ") {
				writeJSONError(w, http.StatusUnauthorized, "missing credentials")
				return
			}
			claims, err := jwtMgr.Parse(strings.TrimPrefix(h, "Bearer "))
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, "invalid token")
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UsernameKey, claims.Username)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) int64 {
	if v, ok := ctx.Value(UserIDKey).(int64); ok {
		return v
	}
	return 0
}
