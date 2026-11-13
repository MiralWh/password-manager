package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"pwmanager/internal/models"
	"pwmanager/internal/service"
)

type APIKeyHandler struct {
	svc *service.APIKeyService
}

func NewAPIKeyHandler(svc *service.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{svc: svc}
}

func (h *APIKeyHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context(), UserIDFromContext(r.Context()))
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []models.APIKey{}
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *APIKeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in models.APIKeyInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid body")
		return
	}
	k, err := h.svc.Create(r.Context(), UserIDFromContext(r.Context()), in)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, k)
}

func (h *APIKeyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := h.svc.Delete(r.Context(), UserIDFromContext(r.Context()), id); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
