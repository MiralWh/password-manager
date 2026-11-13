package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pwmanager/internal/testutil"
)

func TestMiddleware_NoAuth(t *testing.T) {
	ts := testutil.NewTestServer(t)
	rec := doJSON(t, ts, "GET", "/api/v1/passwords", "", false)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", rec.Code)
	}
}

func TestMiddleware_WrongToken(t *testing.T) {
	ts := testutil.NewTestServer(t)
	req := httptest.NewRequest("GET", "/api/v1/passwords", nil)
	req.Header.Set("Authorization", "Bearer garbage")
	rec := httptest.NewRecorder()
	ts.Router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", rec.Code)
	}
}

func TestMiddleware_ValidJWT(t *testing.T) {
	ts := testutil.NewTestServer(t)
	rec := doJSON(t, ts, "GET", "/api/v1/passwords", "", true)
	if rec.Code != http.StatusOK {
		t.Errorf("want 200, got %d", rec.Code)
	}
}

func TestMiddleware_WrongAPIKey(t *testing.T) {
	ts := testutil.NewTestServer(t)
	req := httptest.NewRequest("GET", "/api/v1/passwords", nil)
	req.Header.Set("X-API-Key", "pwm_invalidkey1234567890")
	rec := httptest.NewRecorder()
	ts.Router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", rec.Code)
	}
}

func TestMiddleware_APIKeyNotAllowedForKeysMgmt(t *testing.T) {
	ts := testutil.NewTestServer(t)
	rec := doJSON(t, ts, "POST", "/api/v1/auth/api-keys", `{"name":"desktop"}`, true)
	var k map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&k)

	// /auth/api-keys — только JWT, ключ не должен подходить
	req := httptest.NewRequest("GET", "/api/v1/auth/api-keys", nil)
	req.Header.Set("X-API-Key", k["key"].(string))
	rec = httptest.NewRecorder()
	ts.Router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("api-keys must be JWT-only, got %d", rec.Code)
	}
}
