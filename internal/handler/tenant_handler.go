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

type TenantHandler struct {
	repo *repository.TenantRepository
}

func NewTenantHandler(repo *repository.TenantRepository) *TenantHandler {
	return &TenantHandler{repo: repo}
}

// Create handles POST /api/admin/tenants
func (h *TenantHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tenant, err := h.repo.Create(&req)
	if err != nil {
		slog.Error("failed to create tenant", "error", err)
		http.Error(w, "Failed to create tenant", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tenant)
}

// Get handles GET /api/admin/tenants/:id
func (h *TenantHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid tenant ID", http.StatusBadRequest)
		return
	}

	tenant, err := h.repo.GetByID(id)
	if err != nil {
		slog.Error("failed to get tenant", "error", err, "id", id)
		http.Error(w, "Tenant not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenant)
}

// List handles GET /api/admin/tenants
func (h *TenantHandler) List(w http.ResponseWriter, r *http.Request) {
	tenants, err := h.repo.List()
	if err != nil {
		slog.Error("failed to list tenants", "error", err)
		http.Error(w, "Failed to list tenants", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  tenants,
		"count": len(tenants),
	})
}
