package handler

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/riz/auto-lmk/internal/repository"
)

type PageHandler struct {
	carRepo    *repository.CarRepository
	tenantRepo *repository.TenantRepository
	templates  *template.Template
}

func NewPageHandler(carRepo *repository.CarRepository, tenantRepo *repository.TenantRepository) *PageHandler {
	// Parse all templates
	tmpl := template.Must(template.ParseFiles(
		"templates/layouts/base.html",
		"templates/pages/home.html",
		"templates/pages/cars.html",
		"templates/pages/car-detail.html",
		"templates/pages/contact.html",
		"templates/admin/layout.html",
		"templates/admin/dashboard.html",
		"templates/admin/cars.html",
		"templates/admin/leads.html",
		"templates/admin/whatsapp.html",
	))

	return &PageHandler{
		carRepo:    carRepo,
		tenantRepo: tenantRepo,
		templates:  tmpl,
	}
}

// Home renders the homepage
func (h *PageHandler) Home(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "Beranda"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Cars renders the cars listing page
func (h *PageHandler) Cars(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "Daftar Mobil"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CarDetail renders the car detail page
func (h *PageHandler) CarDetail(w http.ResponseWriter, r *http.Request) {
	carID := chi.URLParam(r, "id")

	data := h.getDefaultData(r)
	data["Title"] = "Detail Mobil"
	data["CarID"] = carID

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Contact renders the contact page
func (h *PageHandler) Contact(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "Kontak Kami"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminDashboard renders the admin dashboard
func (h *PageHandler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "Dashboard Admin"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "admin/layout.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminCars renders the admin cars management page
func (h *PageHandler) AdminCars(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "Kelola Mobil"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "admin/layout.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminLeads renders the admin leads management page
func (h *PageHandler) AdminLeads(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "Kelola Lead"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "admin/layout.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminWhatsApp renders the admin WhatsApp management page
func (h *PageHandler) AdminWhatsApp(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "WhatsApp Bot"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "admin/layout.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// getDefaultData returns default template data with tenant info
func (h *PageHandler) getDefaultData(r *http.Request) map[string]interface{} {
	// For now, use default values
	// TODO: Get tenant from context when tenant middleware is implemented for frontend routes
	return map[string]interface{}{
		"TenantName":     "Auto LMK",
		"WhatsAppNumber": "6281234567890",
		"Description":    "Platform penjualan mobil terpercaya",
		"Keywords":       "mobil, jual mobil, beli mobil",
		"FeaturedCars":   []interface{}{},
	}
}
