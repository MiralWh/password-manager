package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"pwmanager/internal/auth"
	"pwmanager/internal/models"
	"pwmanager/internal/repository"
	"pwmanager/internal/service"
	"pwmanager/internal/testutil"
)

func newAuthSvc(t *testing.T) *service.AuthService {
	pool, cfg := testutil.NewTestDB(t)
	users := repository.NewUserRepository(pool)
	jwt := auth.NewJWTManager(cfg.Server.JWTSecret, time.Hour)
	return service.NewAuthService(users, jwt)
}

func TestAuthService_Login(t *testing.T) {
	svc := newAuthSvc(t)
	resp, err := svc.Login(context.Background(), models.LoginInput{Username: "admin", Password: "admin123"})
	if err != nil { t.Fatalf("Login: %v", err) }
	if resp.Token == "" || resp.Username != "admin" {
		t.Errorf("unexpected: %+v", resp)
	}
}

func TestAuthService_Login_Wrong(t *testing.T) {
	svc := newAuthSvc(t)
	cases := []models.LoginInput{
		{Username: "admin", Password: "wrong"},
		{Username: "ghost", Password: "admin123"},
	}
	for _, in := range cases {
		if _, err := svc.Login(context.Background(), in); !errors.Is(err, service.ErrInvalidCredentials) {
			t.Errorf("want ErrInvalidCredentials, got %v", err)
		}
	}
}

func TestAuthService_Register(t *testing.T) {
	svc := newAuthSvc(t)
	resp, err := svc.Register(context.Background(), models.RegisterInput{Username: "newuser", Password: "strongpw"})
	if err != nil { t.Fatalf("Register: %v", err) }
	if resp.Token == "" { t.Error("token empty") }

	if _, err := svc.Register(context.Background(), models.RegisterInput{Username: "newuser", Password: "strongpw"}); err == nil {
		t.Error("duplicate register should fail")
	}
}

func TestAuthService_Register_Validation(t *testing.T) {
	svc := newAuthSvc(t)
	for _, in := range []models.RegisterInput{
		{Username: "ab", Password: "longenough"},
		{Username: "alice", Password: "short"},
	} {
		if _, err := svc.Register(context.Background(), in); err == nil {
			t.Errorf("expected error for %+v", in)
		}
	}
}
