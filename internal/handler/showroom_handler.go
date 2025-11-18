package handler

import (
	"encoding/json"
	"net/http"

	"github.com/riz/auto-lmk/internal/model"
	"github.com/riz/auto-lmk/internal/repository"
)

type ShowroomHandler struct {
	showroomRepo *repository.ShowroomRepository
}

func NewShowroomHandler(showroomRepo *repository.ShowroomRepository) *ShowroomHandler {
	return &ShowroomHandler{
		showroomRepo: showroomRepo,
	}
}

// GetSettings retrieves showroom settings for public display
func (h *ShowroomHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		http.Error(w, "Tenant not found", http.StatusBadRequest)
		return
	}

	showroom, err := h.showroomRepo.GetByTenantID(r.Context(), tenantID)
	if err != nil {
		http.Error(w, "Failed to get showroom settings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(showroom)
}

// GetAdminSettings retrieves showroom settings for admin view
func (h *ShowroomHandler) GetAdminSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		http.Error(w, "Tenant not found", http.StatusBadRequest)
		return
	}

	showroom, err := h.showroomRepo.GetByTenantID(r.Context(), tenantID)
	if err != nil {
		http.Error(w, "Failed to get showroom settings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(showroom)
}

// UpdateSettings updates showroom settings
func (h *ShowroomHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		http.Error(w, "Tenant not found", http.StatusBadRequest)
		return
	}

	var req model.ShowroomUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create showroom settings object
	showroom := &model.ShowroomSettings{
		TenantID:      tenantID,
		Address:       req.Address,
		Phone:         req.Phone,
		Email:         req.Email,
		BusinessHours: req.BusinessHours,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		MapEmbed:      req.MapEmbed,
	}

	// Update in database
	if err := h.showroomRepo.CreateOrUpdate(r.Context(), showroom); err != nil {
		http.Error(w, "Failed to update showroom settings", http.StatusInternalServerError)
		return
	}

	// Return updated settings
	updatedShowroom, err := h.showroomRepo.GetByTenantID(r.Context(), tenantID)
	if err != nil {
		http.Error(w, "Failed to get updated settings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedShowroom)
}
