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

type SalesHandler struct {
	repo *repository.SalesRepository
}

func NewSalesHandler(repo *repository.SalesRepository) *SalesHandler {
	return &SalesHandler{repo: repo}
}

// Create handles POST /api/sales
func (h *SalesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateSalesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	sales, err := h.repo.Create(r.Context(), &req)
	if err != nil {
		slog.Error("failed to create sales", "error", err)
		http.Error(w, "Failed to create sales", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sales)
}

// List handles GET /api/sales
func (h *SalesHandler) List(w http.ResponseWriter, r *http.Request) {
	salesList, err := h.repo.List(r.Context())
	if err != nil {
		slog.Error("failed to list sales", "error", err)
		http.Error(w, "Failed to list sales", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  salesList,
		"count": len(salesList),
	})
}

// Delete handles DELETE /api/sales/:id
func (h *SalesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid sales ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		slog.Error("failed to delete sales", "error", err, "id", id)
		http.Error(w, "Failed to delete sales", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sales deleted successfully"}`))
}
