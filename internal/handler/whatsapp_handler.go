package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/riz/auto-lmk/internal/model"
	"github.com/riz/auto-lmk/internal/repository"
	"github.com/riz/auto-lmk/internal/whatsapp"
)

// WhatsAppHandler handles WhatsApp admin operations
type WhatsAppHandler struct {
	waClient   *whatsapp.Client
	tenantRepo *repository.TenantRepository
}

func NewWhatsAppHandler(waClient *whatsapp.Client, tenantRepo *repository.TenantRepository) *WhatsAppHandler {
	return &WhatsAppHandler{
		waClient:   waClient,
		tenantRepo: tenantRepo,
	}
}

// PairingStatus represents WhatsApp pairing status
type PairingStatus struct {
	TenantID       int       `json:"tenant_id"`
	IsConnected    bool      `json:"is_connected"`
	PhoneNumber    string    `json:"phone_number,omitempty"`
	PairingStatus  string    `json:"pairing_status"`
	LastConnected  *time.Time `json:"last_connected,omitempty"`
	QRCode         string    `json:"qr_code,omitempty"`
}

// GetStatus returns WhatsApp connection status for current tenant
func (h *WhatsAppHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get tenant info
	tenant, err := h.tenantRepo.GetByID(tenantID)
	if err != nil {
		http.Error(w, "Tenant not found", http.StatusNotFound)
		return
	}

	// Check connection status
	isConnected := h.waClient.IsConnected(tenantID)

	phoneNumber := ""
	if tenant.WhatsAppNumber != nil {
		phoneNumber = *tenant.WhatsAppNumber
	}

	status := PairingStatus{
		TenantID:      tenantID,
		IsConnected:   isConnected,
		PhoneNumber:   phoneNumber,
		PairingStatus: tenant.PairingStatus,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// InitiatePairing starts WhatsApp pairing process
func (h *WhatsAppHandler) InitiatePairing(w http.ResponseWriter, r *http.Request) {
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get tenant info
	tenant, err := h.tenantRepo.GetByID(tenantID)
	if err != nil {
		http.Error(w, "Tenant not found", http.StatusNotFound)
		return
	}

	if tenant.WhatsAppNumber == nil || *tenant.WhatsAppNumber == "" {
		http.Error(w, "WhatsApp number not configured", http.StatusBadRequest)
		return
	}

	// Check if already connected
	if h.waClient.IsConnected(tenantID) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "already_connected",
			"message": "WhatsApp already connected",
		})
		return
	}

	// Start pairing
	phoneNumber := *tenant.WhatsAppNumber
	slog.Info("initiating WhatsApp pairing", "tenant_id", tenantID, "phone", phoneNumber)

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	qrCode, err := h.waClient.PairTenant(ctx, tenantID, phoneNumber)
	if err != nil {
		slog.Error("failed to initiate pairing", "error", err)
		http.Error(w, fmt.Sprintf("Failed to initiate pairing: %v", err), http.StatusInternalServerError)
		return
	}

	// Update tenant pairing status
	// TODO: Update tenant.pairing_status in database

	response := map[string]interface{}{
		"status":  "pairing_initiated",
		"qr_code": qrCode,
		"message": "Scan the QR code with WhatsApp to pair",
	}

	if qrCode == "PAIRED_SUCCESSFULLY" {
		response["status"] = "paired"
		response["message"] = "WhatsApp paired successfully"
		delete(response, "qr_code")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Disconnect disconnects WhatsApp for current tenant
func (h *WhatsAppHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.waClient.Disconnect(tenantID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to disconnect: %v", err), http.StatusInternalServerError)
		return
	}

	// TODO: Update tenant pairing_status in database

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "disconnected",
		"message": "WhatsApp disconnected successfully",
	})
}

// SendTestMessage sends a test message to verify bot is working
func (h *WhatsAppHandler) SendTestMessage(w http.ResponseWriter, r *http.Request) {
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		PhoneNumber string `json:"phone_number"`
		Message     string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.PhoneNumber == "" {
		http.Error(w, "phone_number is required", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		req.Message = "Halo! Ini adalah pesan test dari Auto LMK Bot. Bot WhatsApp Anda sudah berhasil terhubung! 🎉"
	}

	// Send test message
	err = h.waClient.SendMessage(tenantID, req.PhoneNumber, req.Message)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send message: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "sent",
		"message": "Test message sent successfully",
	})
}

// GetQRCodeImage serves the QR code image for pairing
func (h *WhatsAppHandler) GetQRCodeImage(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := chi.URLParam(r, "tenant_id")
	if tenantIDStr == "" {
		http.Error(w, "tenant_id required", http.StatusBadRequest)
		return
	}

	qrPath := fmt.Sprintf("/tmp/qr_%s.png", tenantIDStr)

	http.ServeFile(w, r, qrPath)
}
