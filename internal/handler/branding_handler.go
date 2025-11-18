package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/riz/auto-lmk/internal/model"
	"github.com/riz/auto-lmk/internal/repository"
)

type BrandingHandler struct {
	brandingRepo *repository.BrandingRepository
}

func NewBrandingHandler(brandingRepo *repository.BrandingRepository) *BrandingHandler {
	return &BrandingHandler{
		brandingRepo: brandingRepo,
	}
}

// GetSettings retrieves branding settings for the current tenant
func (h *BrandingHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		http.Error(w, "Tenant not found", http.StatusBadRequest)
		return
	}

	branding, err := h.brandingRepo.GetByTenantID(r.Context(), tenantID)
	if err != nil {
		slog.Error("failed to get branding settings", "error", err, "tenant_id", tenantID)
		http.Error(w, "Failed to load branding settings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(branding)
}

// UpdateSettings updates branding settings for the current tenant
func (h *BrandingHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		http.Error(w, "Tenant not found", http.StatusBadRequest)
		return
	}

	var req model.BrandingUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get existing branding to preserve file paths
	existing, err := h.brandingRepo.GetByTenantID(r.Context(), tenantID)
	if err != nil {
		slog.Error("failed to get existing branding", "error", err, "tenant_id", tenantID)
		http.Error(w, "Failed to load existing settings", http.StatusInternalServerError)
		return
	}

	// Update only the text fields and header style
	branding := &model.BrandingSettings{
		TenantID:       tenantID,
		LogoPath:       existing.LogoPath,
		FaviconPath:    existing.FaviconPath,
		CustomTitle:    req.CustomTitle,
		CustomSubtitle: req.CustomSubtitle,
		PromoText:      req.PromoText,
		HeaderStyle:    req.HeaderStyle,
	}

	if err := h.brandingRepo.CreateOrUpdate(r.Context(), branding); err != nil {
		slog.Error("failed to update branding settings", "error", err, "tenant_id", tenantID)
		http.Error(w, "Failed to save branding settings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Branding settings updated successfully"})
}

// UploadLogo handles logo file upload for the current tenant
func (h *BrandingHandler) UploadLogo(w http.ResponseWriter, r *http.Request) {
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		http.Error(w, "Tenant not found", http.StatusBadRequest)
		return
	}

	// Parse multipart form (max 5MB)
	err = r.ParseMultipartForm(5 << 20)
	if err != nil {
		slog.Error("failed to parse multipart form", "error", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, fileHeader, err := r.FormFile("logo")
	if err != nil {
		http.Error(w, "No logo file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file size (2MB max for logo)
	if fileHeader.Size > 2<<20 {
		http.Error(w, "Logo file is too large (max 2MB)", http.StatusBadRequest)
		return
	}

	// Validate file type
	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		http.Error(w, "File must be an image", http.StatusBadRequest)
		return
	}

	// Create upload directory
	uploadDir := filepath.Join("uploads", "branding", strconv.Itoa(tenantID))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		slog.Error("failed to create upload directory", "error", err, "dir", uploadDir)
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// Generate unique filename
	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("logo_%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		slog.Error("failed to create destination file", "error", err, "path", filePath)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		slog.Error("failed to copy uploaded file", "error", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Update logo path in database (relative path from uploads/)
	logoPath := filepath.Join("/uploads", "branding", strconv.Itoa(tenantID), filename)
	if err := h.brandingRepo.UpdateLogo(r.Context(), tenantID, logoPath); err != nil {
		slog.Error("failed to update logo path in database", "error", err, "tenant_id", tenantID)
		http.Error(w, "Failed to update logo", http.StatusInternalServerError)
		return
	}

	slog.Info("logo uploaded successfully", "tenant_id", tenantID, "path", logoPath)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Logo uploaded successfully",
		"logo_path": logoPath,
	})
}

// UploadFavicon handles favicon file upload for the current tenant
func (h *BrandingHandler) UploadFavicon(w http.ResponseWriter, r *http.Request) {
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		http.Error(w, "Tenant not found", http.StatusBadRequest)
		return
	}

	// Parse multipart form (max 1MB)
	err = r.ParseMultipartForm(1 << 20)
	if err != nil {
		slog.Error("failed to parse multipart form", "error", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, fileHeader, err := r.FormFile("favicon")
	if err != nil {
		http.Error(w, "No favicon file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file size (512KB max for favicon)
	if fileHeader.Size > 512<<10 {
		http.Error(w, "Favicon file is too large (max 512KB)", http.StatusBadRequest)
		return
	}

	// Validate file type (accept image/* and .ico)
	contentType := fileHeader.Header.Get("Content-Type")
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !strings.HasPrefix(contentType, "image/") && ext != ".ico" {
		http.Error(w, "File must be an image or .ico file", http.StatusBadRequest)
		return
	}

	// Create upload directory
	uploadDir := filepath.Join("uploads", "branding", strconv.Itoa(tenantID))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		slog.Error("failed to create upload directory", "error", err, "dir", uploadDir)
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// Use fixed filename for favicon
	filename := "favicon" + ext
	filePath := filepath.Join(uploadDir, filename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		slog.Error("failed to create destination file", "error", err, "path", filePath)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		slog.Error("failed to copy uploaded file", "error", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Update favicon path in database
	faviconPath := filepath.Join("/uploads", "branding", strconv.Itoa(tenantID), filename)
	if err := h.brandingRepo.UpdateFavicon(r.Context(), tenantID, faviconPath); err != nil {
		slog.Error("failed to update favicon path in database", "error", err, "tenant_id", tenantID)
		http.Error(w, "Failed to update favicon", http.StatusInternalServerError)
		return
	}

	slog.Info("favicon uploaded successfully", "tenant_id", tenantID, "path", faviconPath)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":      "Favicon uploaded successfully",
		"favicon_path": faviconPath,
	})
}
