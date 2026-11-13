package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"pwmanager/internal/testutil"
)

func TestGroupHandler_CRUD(t *testing.T) {
	ts := testutil.NewTestServer(t)

	rec := doJSON(t, ts, "POST", "/api/v1/groups", `{"name":"Work"}`, true)
	if rec.Code != http.StatusCreated { t.Fatalf("create %d", rec.Code) }
	var g map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&g)
	idStr := fmt.Sprintf("%d", int64(g["id"].(float64)))

	rec = doJSON(t, ts, "GET", "/api/v1/groups", "", true)
	var list []map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&list)
	if len(list) != 1 { t.Errorf("list: %d", len(list)) }

	rec = doJSON(t, ts, "GET", "/api/v1/groups?q=Work", "", true)
	json.NewDecoder(rec.Body).Decode(&list)
	if len(list) != 1 { t.Errorf("search: %d", len(list)) }

	rec = doJSON(t, ts, "PUT", "/api/v1/groups/"+idStr, `{"name":"Job"}`, true)
	if rec.Code != http.StatusOK { t.Fatalf("update %d", rec.Code) }

	rec = doJSON(t, ts, "DELETE", "/api/v1/groups/"+idStr, "", true)
	if rec.Code != http.StatusNoContent { t.Fatalf("delete %d", rec.Code) }
}

func TestGroupHandler_Validation(t *testing.T) {
	ts := testutil.NewTestServer(t)
	rec := doJSON(t, ts, "POST", "/api/v1/groups", `{"name":""}`, true)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}
