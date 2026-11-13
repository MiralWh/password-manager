package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pwmanager/internal/testutil"
)

func TestAPIKeyHandler_FullFlow(t *testing.T) {
	ts := testutil.NewTestServer(t)

	// Create
	rec := doJSON(t, ts, "POST", "/api/v1/auth/api-keys", `{"name":"desktop"}`, true)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create %d: %s", rec.Code, rec.Body.String())
	}
	var k map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&k)
	key := k["key"].(string)
	idStr := fmt.Sprintf("%d", int64(k["id"].(float64)))

	if !strings.HasPrefix(key, "pwm_") {
		t.Errorf("bad key: %q", key)
	}

	// List (без plaintext)
	rec = doJSON(t, ts, "GET", "/api/v1/auth/api-keys", "", true)
	var list []map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&list)
	if len(list) != 1 {
		t.Errorf("list: %d", len(list))
	}
	if _, ok := list[0]["key"]; ok {
		t.Error("plaintext leaked in list response")
	}

	// Использование ключа для доступа к /passwords
	req := httptest.NewRequest("GET", "/api/v1/passwords", nil)
	req.Header.Set("X-API-Key", key)
	rec = httptest.NewRecorder()
	ts.Router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("API key auth failed: %d", rec.Code)
	}

	// Delete
	rec = doJSON(t, ts, "DELETE", "/api/v1/auth/api-keys/"+idStr, "", true)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete %d", rec.Code)
	}

	// Ключ больше не работает
	req = httptest.NewRequest("GET", "/api/v1/passwords", nil)
	req.Header.Set("X-API-Key", key)
	rec = httptest.NewRecorder()
	ts.Router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("deleted key still works: %d", rec.Code)
	}
}
