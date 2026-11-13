package auth_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"pwmanager/internal/auth"
)

const secret = "test-secret-for-jwt-manager-1234567890"

func TestGenerateParse_RoundTrip(t *testing.T) {
	m := auth.NewJWTManager(secret, time.Hour)
	tok, exp, err := m.Generate(42, "admin")
	if err != nil { t.Fatalf("gen: %v", err) }
	if tok == "" { t.Fatal("empty token") }
	if delta := time.Until(exp); delta < 59*time.Minute || delta > 61*time.Minute {
		t.Errorf("exp offset: %v", delta)
	}
	claims, err := m.Parse(tok)
	if err != nil { t.Fatalf("parse: %v", err) }
	if claims.UserID != 42 || claims.Username != "admin" {
		t.Errorf("claims: %+v", claims)
	}
}

func TestParse_Invalid(t *testing.T) {
	m := auth.NewJWTManager(secret, time.Hour)
	for _, s := range []string{"", "garbage", "a.b.c"} {
		if _, err := m.Parse(s); err == nil {
			t.Errorf("expected error for %q", s)
		}
	}
}

func TestParse_WrongSecret(t *testing.T) {
	m1 := auth.NewJWTManager(secret, time.Hour)
	m2 := auth.NewJWTManager("different-secret-entirely-xxxxxxxxxx", time.Hour)
	tok, _, _ := m1.Generate(1, "u")
	if _, err := m2.Parse(tok); err == nil {
		t.Error("secret mismatch not detected")
	}
}

func TestParse_Expired(t *testing.T) {
	m := auth.NewJWTManager(secret, time.Hour)
	claims := &auth.Claims{UserID: 1, Username: "u",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		}}
	tok, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if _, err := m.Parse(tok); err == nil {
		t.Error("expired not detected")
	}
}

func TestParse_WrongMethod(t *testing.T) {
	m := auth.NewJWTManager(secret, time.Hour)
	claims := &auth.Claims{UserID: 1, Username: "u",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		}}
	tok, _ := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(secret))
	if _, err := m.Parse(tok); err == nil {
		t.Error("HS512 not rejected")
	}
}
