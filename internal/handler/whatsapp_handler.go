package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/riz/auto-lmk/internal/model"
	"github.com/riz/auto-lmk/internal/repository"
	"github.com/riz/auto-lmk/internal/whatsapp"
)

// WhatsAppHandler handles WhatsApp admin operations
type WhatsAppHandler struct {
	waClient     *whatsapp.Client
	tenantRepo   *repository.TenantRepository
	settingsRepo *repository.WhatsAppSettingsRepository
}

func NewWhatsAppHandler(waClient *whatsapp.Client, tenantRepo *repository.TenantRepository) *WhatsAppHandler {
	return &WhatsAppHandler{
		waClient:   waClient,
		tenantRepo: tenantRepo,
	}
}

func NewWhatsAppHandlerWithSettings(waClient *whatsapp.Client, tenantRepo *repository.TenantRepository, settingsRepo *repository.WhatsAppSettingsRepository) *WhatsAppHandler {
	return &WhatsAppHandler{
		waClient:     waClient,
		tenantRepo:   tenantRepo,
		settingsRepo: settingsRepo,
	}
}

// PairingStatus represents WhatsApp pairing status
type PairingStatus struct {
	TenantID      int        `json:"tenant_id"`
	IsConnected   bool       `json:"is_connected"`
	PhoneNumber   string     `json:"phone_number,omitempty"`
	PairingStatus string     `json:"pairing_status"`
	LastConnected *time.Time `json:"last_connected,omitempty"`
	QRCode        string     `json:"qr_code,omitempty"`
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

	// If connected but phone number not in database, get it from WhatsApp client and save
	phoneNumber := ""
	if isConnected {
		detectedPhone := h.waClient.GetPhoneNumber(tenantID)
		if detectedPhone != "" {
			phoneNumber = detectedPhone

			// Save to database if not already saved
			if tenant.WhatsAppNumber == nil || *tenant.WhatsAppNumber != detectedPhone {
				if err := h.tenantRepo.UpdateWhatsAppNumber(tenantID, detectedPhone); err != nil {
					slog.Error("failed to update WhatsApp number", "error", err, "tenant_id", tenantID, "phone", detectedPhone)
				} else {
					slog.Info("WhatsApp number saved to database", "tenant_id", tenantID, "phone", detectedPhone)
				}
			}

			// Update pairing status to "paired" if still "pairing_pending"
			if tenant.PairingStatus != "paired" {
				if err := h.tenantRepo.UpdatePairingStatus(tenantID, "paired"); err != nil {
					slog.Error("failed to update pairing status", "error", err, "tenant_id", tenantID)
				}
				tenant.PairingStatus = "paired"
			}
		}
	} else if tenant.WhatsAppNumber != nil {
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
	slog.Info("[HANDLER-START] InitiatePairing request received", "method", r.Method, "url", r.URL.Path)

	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		slog.Error("[HANDLER-ERROR] Failed to get tenant ID from context", "error", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	slog.Info("[HANDLER-TENANT] Processing pairing request", "tenant_id", tenantID)

	// Check if already connected
	if h.waClient.IsConnected(tenantID) {
		slog.Warn("[HANDLER-ABORT] Tenant already connected", "tenant_id", tenantID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "already_connected",
			"message": "WhatsApp already connected",
		})
		return
	}

	// Start pairing - NO PHONE NUMBER REQUIRED!
	// WhatsApp will auto-detect the number from the device that scans the QR
	slog.Info("[HANDLER-PAIRING] Starting pairing process", "tenant_id", tenantID, "timeout", "60s")

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	slog.Info("[HANDLER-CALL] Calling waClient.PairTenant", "tenant_id", tenantID)
	result, err := h.waClient.PairTenant(ctx, tenantID)
	if err != nil {
		slog.Error("[HANDLER-ERROR] PairTenant failed", "tenant_id", tenantID, "error", err)
		http.Error(w, fmt.Sprintf("Failed to initiate pairing: %v", err), http.StatusInternalServerError)
		return
	}

	slog.Info("[HANDLER-RESULT] PairTenant returned", "tenant_id", tenantID, "status", result.Status, "has_qr", result.QRCode != "")

	// Update tenant pairing status in database
	slog.Info("[HANDLER-DB] Updating pairing status in database", "tenant_id", tenantID)
	pairingStatus := "pairing_pending"
	if result.Status == "paired" {
		pairingStatus = "paired"
		slog.Info("[HANDLER-DB] Pairing successful, status set to paired", "tenant_id", tenantID)

		// Save WhatsApp number to database
		if result.PhoneNumber != "" {
			slog.Info("[HANDLER-DB] Saving WhatsApp number", "tenant_id", tenantID, "phone", result.PhoneNumber)
			if err := h.tenantRepo.UpdateWhatsAppNumber(tenantID, result.PhoneNumber); err != nil {
				slog.Error("[HANDLER-ERROR] Failed to update WhatsApp number", "error", err, "tenant_id", tenantID, "phone", result.PhoneNumber)
				// Continue anyway - don't fail the request if DB update fails
			} else {
				slog.Info("[HANDLER-DB] WhatsApp number saved successfully", "tenant_id", tenantID, "phone", result.PhoneNumber)
			}
		}
	} else {
		slog.Info("[HANDLER-DB] QR generated, status set to pairing_pending", "tenant_id", tenantID)
	}

	slog.Info("[HANDLER-DB] Updating pairing status", "tenant_id", tenantID, "status", pairingStatus)
	if err := h.tenantRepo.UpdatePairingStatus(tenantID, pairingStatus); err != nil {
		slog.Error("[HANDLER-ERROR] Failed to update pairing status", "error", err, "tenant_id", tenantID, "status", pairingStatus)
		// Continue anyway - don't fail the request if DB update fails
	} else {
		slog.Info("[HANDLER-DB] Pairing status updated successfully", "tenant_id", tenantID, "status", pairingStatus)
	}

	// Build response
	slog.Info("[HANDLER-RESPONSE] Building response", "tenant_id", tenantID, "result_status", result.Status)
	response := map[string]interface{}{
		"status":  "pairing_initiated",
		"message": "Scan the QR code with WhatsApp to pair",
	}

	if result.Status == "paired" {
		response["status"] = "paired"
		response["message"] = "WhatsApp paired successfully"
		response["phone_number"] = result.PhoneNumber
		slog.Info("[HANDLER-RESPONSE] Sending paired response", "tenant_id", tenantID, "phone", result.PhoneNumber)
	} else {
		response["qr_code"] = result.QRCode
		qr_length := len(result.QRCode)
		slog.Info("[HANDLER-RESPONSE] Sending QR code response", "tenant_id", tenantID, "qr_length", qr_length)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("[HANDLER-ERROR] Failed to encode JSON response", "tenant_id", tenantID, "error", err)
	} else {
		slog.Info("[HANDLER-SUCCESS] Response sent successfully", "tenant_id", tenantID)
	}
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

	// Update tenant pairing status in database to "disconnected"
	if err := h.tenantRepo.UpdatePairingStatus(tenantID, "disconnected"); err != nil {
		slog.Error("failed to update pairing status after disconnect", "error", err, "tenant_id", tenantID)
		// Continue anyway - don't fail the request if DB update fails
	} else {
		slog.Info("pairing status updated to disconnected", "tenant_id", tenantID)
	}

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

	// Tambahkan parameter untuk force refresh
	refresh := r.URL.Query().Get("refresh") == "1"
	if refresh {
		slog.Info("QR image refresh requested", "tenant_id", tenantIDStr)
	}

	qrPath := fmt.Sprintf("/tmp/qr_%s.png", tenantIDStr)

	// Check if QR file exists
	if _, err := os.Stat(qrPath); os.IsNotExist(err) {
		http.Error(w, "QR code not found. Please initiate pairing first.", http.StatusNotFound)
		return
	}

	// Set proper headers for PNG dengan cache-busting
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Tambahkan timestamp untuk cache-busting
	timestamp := time.Now().Unix()
	w.Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))
	w.Header().Set("ETag", fmt.Sprintf(`"%d"`, timestamp))

	// Tambahkan header untuk mencegah cache
	w.Header().Set("Cache-Buster", fmt.Sprintf("%d", timestamp))

	http.ServeFile(w, r, qrPath)
}

// GetSettings returns WhatsApp settings for current tenant
func (h *WhatsAppHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if h.settingsRepo == nil {
		http.Error(w, "Settings repository not available", http.StatusServiceUnavailable)
		return
	}

	settings, err := h.settingsRepo.GetSettings(r.Context(), tenantID)
	if err != nil {
		slog.Error("failed to get WhatsApp settings", "error", err, "tenant_id", tenantID)
		http.Error(w, "Failed to get settings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

// UpdateSettings updates WhatsApp settings for current tenant
func (h *WhatsAppHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if h.settingsRepo == nil {
		http.Error(w, "Settings repository not available", http.StatusServiceUnavailable)
		return
	}

	var settings model.WhatsAppSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate settings
	if settings.FallbackNumber == "" {
		http.Error(w, "fallback_number is required", http.StatusBadRequest)
		return
	}

	// Set tenant ID
	settings.TenantID = tenantID

	if err := h.settingsRepo.SaveSettings(r.Context(), &settings); err != nil {
		slog.Error("failed to save WhatsApp settings", "error", err, "tenant_id", tenantID)
		http.Error(w, "Failed to save settings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Settings updated successfully",
	})
}

// GetEffectiveNumber returns the effective WhatsApp number (bot if available, fallback otherwise)
func (h *WhatsAppHandler) GetEffectiveNumber(w http.ResponseWriter, r *http.Request) {
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if h.settingsRepo == nil {
		http.Error(w, "Settings repository not available", http.StatusServiceUnavailable)
		return
	}

	settings, err := h.settingsRepo.GetSettings(r.Context(), tenantID)
	if err != nil {
		slog.Error("failed to get WhatsApp settings", "error", err, "tenant_id", tenantID)
		http.Error(w, "Failed to get settings", http.StatusInternalServerError)
		return
	}

	// Check if bot is connected and has a number
	effectiveNumber := settings.FallbackNumber // Default to fallback
	isBotAvailable := false

	if h.waClient != nil && h.waClient.IsConnected(tenantID) && settings.BotNumber != nil && *settings.BotNumber != "" {
		effectiveNumber = *settings.BotNumber
		isBotAvailable = true
	}

	response := map[string]interface{}{
		"effective_number": effectiveNumber,
		"is_bot_available": isBotAvailable,
		"fallback_number":  settings.FallbackNumber,
		"bot_number":       settings.BotNumber,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
