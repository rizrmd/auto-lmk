package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/riz/auto-lmk/internal/model"
	"github.com/riz/auto-lmk/internal/repository"
)

type LeadHandler struct {
	repo *repository.LeadRepository
}

func NewLeadHandler(repo *repository.LeadRepository) *LeadHandler {
	return &LeadHandler{repo: repo}
}

// Create handles POST /api/leads
func (h *LeadHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	lead, err := h.repo.Create(r.Context(), &req)
	if err != nil {
		slog.Error("failed to create lead", "error", err)
		http.Error(w, "Failed to create lead", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(lead)
}

// List handles GET /api/leads
func (h *LeadHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	leads, err := h.repo.List(r.Context(), status)
	if err != nil {
		slog.Error("failed to list leads", "error", err)
		http.Error(w, "Failed to list leads", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  leads,
		"count": len(leads),
	})
}

// UpdateStatus handles PUT /api/leads/:id/status
func (h *LeadHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid lead ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateStatus(r.Context(), id, req.Status); err != nil {
		slog.Error("failed to update lead status", "error", err, "id", id)
		http.Error(w, "Failed to update lead status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Lead status updated successfully"}`))
}
