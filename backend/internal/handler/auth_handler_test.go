package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pwmanager/internal/testutil"
)

func TestAuthHandler_Login_Success(t *testing.T) {
	ts := testutil.NewTestServer(t)
	req := httptest.NewRequest("POST", "/api/v1/auth/login",
		strings.NewReader(`{"username":"admin","password":"admin123"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ts.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Token    string `json:"token"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Token == "" || resp.Username != "admin" {
		t.Errorf("unexpected: %+v", resp)
	}
}

func TestAuthHandler_Login_BadPassword(t *testing.T) {
	ts := testutil.NewTestServer(t)
	req := httptest.NewRequest("POST", "/api/v1/auth/login",
		strings.NewReader(`{"username":"admin","password":"wrong"}`))
	rec := httptest.NewRecorder()
	ts.Router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", rec.Code)
	}
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	ts := testutil.NewTestServer(t)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader("not-json"))
	rec := httptest.NewRecorder()
	ts.Router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}

func TestAuthHandler_Register(t *testing.T) {
	ts := testutil.NewTestServer(t)
	req := httptest.NewRequest("POST", "/api/v1/auth/register",
		strings.NewReader(`{"username":"newbie","password":"strongpass"}`))
	rec := httptest.NewRecorder()
	ts.Router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestAuthHandler_Health(t *testing.T) {
	ts := testutil.NewTestServer(t)
	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	ts.Router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("want 200, got %d", rec.Code)
	}
}

func TestAuthHandler_SwaggerDoc(t *testing.T) {
	ts := testutil.NewTestServer(t)
	req := httptest.NewRequest("GET", "/swagger/doc.json", nil)
	rec := httptest.NewRecorder()
	ts.Router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK { t.Fatalf("status %d", rec.Code) }
	if !strings.Contains(rec.Body.String(), `"openapi"`) {
		t.Error("not swagger json")
	}
}
