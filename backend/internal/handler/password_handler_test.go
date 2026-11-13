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

func doJSON(t *testing.T, ts *testutil.TestServer, method, path, body string, auth bool) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if auth { ts.SetAuth(req) }
	rec := httptest.NewRecorder()
	ts.Router.ServeHTTP(rec, req)
	return rec
}

func TestPasswordHandler_CRUD(t *testing.T) {
	ts := testutil.NewTestServer(t)

	// Create
	rec := doJSON(t, ts, "POST", "/api/v1/passwords",
		`{"title":"GitHub","username":"octocat","password":"s3cr3t","url":"https://github.com"}`, true)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create %d: %s", rec.Code, rec.Body.String())
	}
	var p map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&p)
	id := int64(p["id"].(float64))
	idStr := fmt.Sprintf("%d", id)

	// Get
	rec = doJSON(t, ts, "GET", "/api/v1/passwords/"+idStr, "", true)
	if rec.Code != http.StatusOK { t.Fatalf("get %d", rec.Code) }

	// Update
	rec = doJSON(t, ts, "PUT", "/api/v1/passwords/"+idStr,
		`{"title":"GH2","username":"u2","password":"p2"}`, true)
	if rec.Code != http.StatusOK {
		t.Fatalf("update %d: %s", rec.Code, rec.Body.String())
	}

	// List
	rec = doJSON(t, ts, "GET", "/api/v1/passwords", "", true)
	var list []map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&list)
	if len(list) != 1 || list[0]["title"] != "GH2" {
		t.Errorf("list: %+v", list)
	}

	// Search
	rec = doJSON(t, ts, "GET", "/api/v1/passwords?q=GH2", "", true)
	json.NewDecoder(rec.Body).Decode(&list)
	if len(list) != 1 { t.Errorf("search: %d", len(list)) }

	rec = doJSON(t, ts, "GET", "/api/v1/passwords?q=nothing", "", true)
	json.NewDecoder(rec.Body).Decode(&list)
	if len(list) != 0 { t.Errorf("no-match search: %d", len(list)) }

	// Delete
	rec = doJSON(t, ts, "DELETE", "/api/v1/passwords/"+idStr, "", true)
	if rec.Code != http.StatusNoContent { t.Fatalf("delete %d", rec.Code) }

	// 404
	rec = doJSON(t, ts, "GET", "/api/v1/passwords/"+idStr, "", true)
	if rec.Code != http.StatusNotFound { t.Errorf("want 404, got %d", rec.Code) }
}

func TestPasswordHandler_Validation(t *testing.T) {
	ts := testutil.NewTestServer(t)
	rec := doJSON(t, ts, "POST", "/api/v1/passwords", `{"title":"x"}`, true)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}

func TestPasswordHandler_BadID(t *testing.T) {
	ts := testutil.NewTestServer(t)
	rec := doJSON(t, ts, "GET", "/api/v1/passwords/abc", "", true)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}
