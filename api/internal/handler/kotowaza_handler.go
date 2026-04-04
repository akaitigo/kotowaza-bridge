package handler

import (
	"net/http"
	"strconv"

	"github.com/akaitigo/kotowaza-bridge/api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// KotowazaHandler handles HTTP requests for kotowaza endpoints.
type KotowazaHandler struct {
	svc *service.KotowazaService
}

// NewKotowazaHandler creates a new KotowazaHandler.
func NewKotowazaHandler(svc *service.KotowazaService) *KotowazaHandler {
	return &KotowazaHandler{svc: svc}
}

// List handles GET /api/v1/kotowaza
func (h *KotowazaHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	resp, err := h.svc.List(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list kotowaza")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// GetByID handles GET /api/v1/kotowaza/{id}
func (h *KotowazaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid kotowaza id")
		return
	}

	k, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get kotowaza")
		return
	}
	if k == nil {
		writeError(w, http.StatusNotFound, "kotowaza not found")
		return
	}

	writeJSON(w, http.StatusOK, k)
}

// Search handles GET /api/v1/kotowaza/search
func (h *KotowazaHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		writeError(w, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	resp, err := h.svc.Search(r.Context(), query, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to search kotowaza")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
