package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/riz/auto-lmk/internal/middleware"
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
	// Extract tenant ID from context
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		slog.Error("tenant ID required for sales creation", "error", err)
		middleware.Unauthorized(w, "Akses tidak sah")
		return
	}

	// Parse request
	var req model.CreateSalesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode sales create request", "error", err, "tenant_id", tenantID)
		middleware.BadRequest(w, "Format data tidak valid")
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		slog.Error("sales create validation failed", "error", err, "tenant_id", tenantID)
		middleware.BadRequest(w, err.Error())
		return
	}

	// Check for duplicate phone number within tenant
	existing, err := h.repo.GetByPhoneNumber(r.Context(), req.PhoneNumber)
	if err != nil && !strings.Contains(err.Error(), "sales not found") {
		slog.Error("failed to check existing sales", "error", err, "tenant_id", tenantID, "phone", req.PhoneNumber)
		middleware.InternalServerError(w, "Terjadi kesalahan server")
		return
	}
	if existing != nil {
		slog.Warn("duplicate phone number for sales", "tenant_id", tenantID, "phone", req.PhoneNumber)
		middleware.Conflict(w, "Nomor telepon sudah terdaftar sebagai sales", map[string]interface{}{
			"field": "phone_number",
			"value": req.PhoneNumber,
		})
		return
	}

	// Create sales member
	sales, err := h.repo.Create(r.Context(), &req)
	if err != nil {
		slog.Error("failed to create sales", "error", err, "tenant_id", tenantID, "phone", req.PhoneNumber)
		middleware.InternalServerError(w, "Gagal membuat sales")
		return
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sales)

	slog.Info("sales member created", "tenant_id", tenantID, "sales_id", sales.ID, "phone", req.PhoneNumber)
}

// List handles GET /api/sales
func (h *SalesHandler) List(w http.ResponseWriter, r *http.Request) {
	// Extract tenant ID from context
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		slog.Error("tenant ID required for sales list", "error", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get sales list
	sales, err := h.repo.List(r.Context())
	if err != nil {
		slog.Error("failed to list sales", "error", err, "tenant_id", tenantID)
		http.Error(w, "Gagal mengambil daftar sales", http.StatusInternalServerError)
		return
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  sales,
		"count": len(sales),
	})
}

// Delete handles DELETE /api/sales/{id}
func (h *SalesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Extract tenant ID from context
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		slog.Error("tenant ID required for sales delete", "error", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract sales ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error("invalid sales ID", "id_str", idStr, "error", err)
		http.Error(w, "ID sales tidak valid", http.StatusBadRequest)
		return
	}

	// Delete sales member
	err = h.repo.Delete(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "sales not found") {
			slog.Warn("sales not found for delete", "tenant_id", tenantID, "sales_id", id)
			http.Error(w, "Sales tidak ditemukan", http.StatusNotFound)
			return
		}
		slog.Error("failed to delete sales", "error", err, "tenant_id", tenantID, "sales_id", id)
		http.Error(w, "Gagal menghapus sales", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
	slog.Info("sales member deleted", "tenant_id", tenantID, "sales_id", id)
}

// Stats handles GET /api/sales/stats
func (h *SalesHandler) Stats(w http.ResponseWriter, r *http.Request) {
	// Extract tenant ID from context
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		slog.Error("tenant ID required for sales stats", "error", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get sales count
	sales, err := h.repo.List(r.Context())
	if err != nil {
		slog.Error("failed to get sales for stats", "error", err, "tenant_id", tenantID)
		http.Error(w, "Failed to get sales stats", http.StatusInternalServerError)
		return
	}

	// Return stats
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_members": len(sales),
	})
}
