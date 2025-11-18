package handler

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/riz/auto-lmk/internal/model"
	"github.com/riz/auto-lmk/internal/repository"
)

type PageHandler struct {
	carRepo          *repository.CarRepository
	salesRepo        *repository.SalesRepository
	tenantRepo       *repository.TenantRepository
	conversationRepo *repository.ConversationRepository
	blogRepo         *repository.BlogRepository
	brandingRepo     *repository.BrandingRepository
	showroomRepo     *repository.ShowroomRepository
	templates        *template.Template
	funcMap          template.FuncMap
}

func NewPageHandler(carRepo *repository.CarRepository, salesRepo *repository.SalesRepository, tenantRepo *repository.TenantRepository, conversationRepo *repository.ConversationRepository, blogRepo *repository.BlogRepository, brandingRepo *repository.BrandingRepository, showroomRepo *repository.ShowroomRepository) *PageHandler {
	// Define template functions
	funcMap := template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
		"add": func(a, b int) int {
			return a + b
		},
		"formatPrice": func(price int64) string {
			// Format: 250000000 -> "Rp 250.000.000"
			priceStr := fmt.Sprintf("%d", price)
			var result []rune
			for i, digit := range priceStr {
				if i > 0 && (len(priceStr)-i)%3 == 0 {
					result = append(result, '.')
				}
				result = append(result, digit)
			}
			return "Rp " + string(result)
		},
		"formatMileage": func(km int) string {
			// Format: 50000 -> "50.000 km"
			kmStr := fmt.Sprintf("%d", km)
			var result []rune
			for i, digit := range kmStr {
				if i > 0 && (len(kmStr)-i)%3 == 0 {
					result = append(result, '.')
				}
				result = append(result, digit)
			}
			return string(result) + " km"
		},
		"whatsappURL": func(phone, carName string) string {
			message := fmt.Sprintf("Halo, saya tertarik dengan %s", carName)
			return fmt.Sprintf("https://wa.me/%s?text=%s", phone, url.QueryEscape(message))
		},
		"default": func(defaultVal interface{}, val interface{}) interface{} {
			if val == nil || val == "" {
				return defaultVal
			}
			return val
		},
		"dict": func(values ...interface{}) map[string]interface{} {
			if len(values)%2 != 0 {
				return nil
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil
				}
				dict[key] = values[i+1]
			}
			return dict
		},
		"seq": func(start, end int) []int {
			// Generate sequence of integers from start to end (inclusive)
			if start > end {
				return []int{}
			}
			seq := make([]int, end-start+1)
			for i := range seq {
				seq[i] = start + i
			}
			return seq
		},
		"len": func(arr interface{}) int {
			// Get length of array/slice or string
			if arr == nil {
				return 0
			}
			switch v := arr.(type) {
			case []interface{}:
				return len(v)
			case []string:
				return len(v)
			case []int:
				return len(v)
			case string:
				return len(v)
			default:
				return 0
			}
		},
		"slice": func(s string, start, end int) string {
			// Slice string safely
			if start < 0 {
				start = 0
			}
			if end > len(s) {
				end = len(s)
			}
			if start >= end || start >= len(s) {
				return ""
			}
			return s[start:end]
		},
		"js": func(s interface{}) string {
			// JavaScript string escaping
			if s == nil {
				return ""
			}
			str := fmt.Sprintf("%v", s)
			str = strings.ReplaceAll(str, "\\", "\\\\")
			str = strings.ReplaceAll(str, "'", "\\'")
			str = strings.ReplaceAll(str, "\"", "\\\"")
			str = strings.ReplaceAll(str, "\n", "\\n")
			str = strings.ReplaceAll(str, "\r", "\\r")
			str = strings.ReplaceAll(str, "\t", "\\t")
			return str
		},
	}

	// Parse all templates with custom functions
	tmpl, err := template.New("").Funcs(funcMap).ParseFiles(
		// Layouts
		"templates/layouts/base.html",
		// Components
		"templates/components/button.html",
		"templates/components/card.html",
		"templates/components/input.html",
		"templates/components/nav.html",
		"templates/components/hero.html",
		"templates/components/footer.html",
		"templates/components/gallery.html",
		"templates/components/pagination.html",
		"templates/components/whatsapp-button.html",
		// Pages
		"templates/pages/home.html",
		"templates/pages/cars.html",
		"templates/pages/car-detail.html",
		"templates/pages/contact.html",
		"templates/pages/blog.html",
		"templates/pages/blog-detail.html",
		// Admin standalone pages (content templates parsed separately in renderAdminPage)
		"templates/admin/whatsapp.html",
	)
	if err != nil {
		fmt.Printf("FATAL: Failed to parse templates: %v\n", err)
		panic(fmt.Sprintf("Template parsing failed: %v", err))
	}

	return &PageHandler{
		carRepo:          carRepo,
		salesRepo:        salesRepo,
		tenantRepo:       tenantRepo,
		conversationRepo: conversationRepo,
		blogRepo:         blogRepo,
		brandingRepo:     brandingRepo,
		showroomRepo:     showroomRepo,
		templates:        tmpl,
		funcMap:          funcMap,
	}
}

// renderAdminPage renders an admin page with specific content template
func (h *PageHandler) renderAdminPage(w http.ResponseWriter, contentFile string, data map[string]interface{}) error {
	fmt.Printf("DEBUG: renderAdminPage called with contentFile=%s\n", contentFile)

	// Parse layout and specific content template
	tmpl, err := template.New("").Funcs(h.funcMap).ParseFiles(
		"templates/admin/layout.html",
		contentFile,
	)
	if err != nil {
		fmt.Printf("DEBUG: ParseFiles error: %v\n", err)
		return err
	}

	fmt.Printf("DEBUG: Template names: %v\n", tmpl.DefinedTemplates())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.ExecuteTemplate(w, "admin/layout.html", data)
}

// Home renders the homepage
func (h *PageHandler) Home(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "Beranda"

	// Load featured cars for homepage
	filters := map[string]interface{}{
		"limit": 6,
	}
	cars, err := h.carRepo.List(r.Context(), filters) // Get 6 featured cars
	if err == nil {
		data["FeaturedCars"] = cars
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "home.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Cars renders the cars listing page
func (h *PageHandler) Cars(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "Daftar Mobil"

	// Load all cars for listing (will be filtered client-side via HTMX)
	filters := map[string]interface{}{
		"limit": 100,
	}
	cars, err := h.carRepo.List(r.Context(), filters)
	if err != nil {
		slog.Error("failed to load cars", "error", err)
		cars = []*model.Car{}
	}
	data["Cars"] = cars

	// Get unique brands for filter
	brandMap := make(map[string]bool)
	for _, car := range cars {
		brandMap[car.Brand] = true
	}
	brands := make([]string, 0, len(brandMap))
	for brand := range brandMap {
		brands = append(brands, brand)
	}
	data["Brands"] = brands

	// Default filter values (can be overridden by query params)
	data["SearchQuery"] = r.URL.Query().Get("search")
	data["SelectedBrand"] = r.URL.Query().Get("brand")
	data["SelectedTransmission"] = r.URL.Query().Get("transmission")
	data["SelectedPriceRange"] = r.URL.Query().Get("price_range")
	data["SelectedFuelType"] = r.URL.Query().Get("fuel_type")
	data["SortBy"] = r.URL.Query().Get("sort")
	if data["SortBy"] == "" {
		data["SortBy"] = "newest"
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "cars.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CarDetail renders the car detail page
func (h *PageHandler) CarDetail(w http.ResponseWriter, r *http.Request) {
	carIDStr := chi.URLParam(r, "id")

	data := h.getDefaultData(r)
	data["Title"] = "Detail Mobil"

	// Convert car ID to int
	carID, err := strconv.Atoi(carIDStr)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	// Load car data
	car, err := h.carRepo.GetByID(r.Context(), carID)
	if err != nil {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	// Create a safe copy with dereferenced pointers for template
	safeCar := struct {
		ID           int       `json:"id"`
		Brand        string    `json:"brand"`
		Model        string    `json:"model"`
		Year         int       `json:"year"`
		Price        int64     `json:"price"`
		Mileage      int       `json:"mileage"`
		Transmission string    `json:"transmission"`
		FuelType     string    `json:"fuel_type"`
		EngineCC     int       `json:"engine_cc"`
		Seats        int       `json:"seats"`
		Color        string    `json:"color"`
		Description  string    `json:"description"`
		Status       string    `json:"status"`
		IsFeatured   bool      `json:"is_featured"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}{
		ID:           car.ID,
		Brand:        car.Brand,
		Model:        car.Model,
		Year:         car.Year,
		Price:        car.Price,
		Status:       car.Status,
		IsFeatured:   car.IsFeatured,
		CreatedAt:    car.CreatedAt,
		UpdatedAt:    car.UpdatedAt,
	}

	// Dereference pointers with safe defaults
	if car.Mileage != nil {
		safeCar.Mileage = *car.Mileage
	}
	if car.Transmission != nil {
		safeCar.Transmission = *car.Transmission
	}
	if car.FuelType != nil {
		safeCar.FuelType = *car.FuelType
	}
	if car.EngineCC != nil {
		safeCar.EngineCC = *car.EngineCC
	}
	if car.Seats != nil {
		safeCar.Seats = *car.Seats
	}
	if car.Color != nil {
		safeCar.Color = *car.Color
	}
	if car.Description != nil {
		safeCar.Description = *car.Description
	}

	data["Car"] = safeCar

	// Load car photos
	photos, err := h.carRepo.GetCarPhotos(r.Context(), carID)
	if err != nil {
		photos = []*model.CarPhoto{} // Empty array if error
	}
	data["Photos"] = photos

	// Load similar cars (same brand, exclude current car)
	filters := map[string]interface{}{
		"brand": car.Brand,
		"limit": 3,
	}
	similarCars, err := h.carRepo.List(r.Context(), filters)
	if err == nil {
		// Filter out the current car
		filtered := make([]*model.Car, 0)
		for _, c := range similarCars {
			if c.ID != car.ID {
				filtered = append(filtered, c)
			}
		}
		data["SimilarCars"] = filtered
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "car-detail.html", data); err != nil {
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
	data["ActiveMenu"] = "dashboard"

	// Fetch stats for dashboard cards
	// 1. Sales Team Stats
	sales, err := h.salesRepo.List(r.Context())
	if err != nil {
		slog.Error("failed to get sales for dashboard", "error", err)
		data["SalesCount"] = 0
	} else {
		data["SalesCount"] = len(sales)
	}

	// 2. Conversation Stats (using all conversations as approximation of last 30 days)
	_, conversationsTotal, err := h.conversationRepo.List(r.Context(), 1, 1000, "")
	if err != nil {
		slog.Error("failed to get conversations for dashboard", "error", err)
		data["ConversationsCount"] = 0
	} else {
		data["ConversationsCount"] = conversationsTotal
	}

	// 3. Cars Stats (TODO: implement when car stats are available)
	data["CarsCount"] = 0

	// 4. WhatsApp Status (will be loaded via HTMX or set default)
	data["WhatsAppConnected"] = false
	data["WhatsAppPhone"] = ""

	// Fetch recent conversations for widget (Story 4.3)
	recentConvs, _, err := h.conversationRepo.List(r.Context(), 1, 5, "")
	if err != nil {
		slog.Error("failed to get recent conversations", "error", err)
		data["RecentConversations"] = []interface{}{}
	} else {
		data["RecentConversations"] = recentConvs
	}

	if err := h.renderAdminPage(w, "templates/admin/dashboard.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminCars renders the admin cars management page
func (h *PageHandler) AdminCars(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "Kelola Mobil"

	// Parse query parameters
	search := r.URL.Query().Get("search")
	brand := r.URL.Query().Get("brand")
	status := r.URL.Query().Get("status")
	sortBy := r.URL.Query().Get("sort")

	// Get cars from repository
	cars, err := h.carRepo.List(r.Context(), map[string]interface{}{})
	if err != nil {
		slog.Error("failed to get cars", "error", err)
		http.Error(w, "Failed to get cars", http.StatusInternalServerError)
		return
	}

	// Filter cars based on query parameters
	var filteredCars []model.Car
	for _, car := range cars {
		// Brand filter
		if brand != "" && car.Brand != brand {
			continue
		}
		// Status filter
		if status != "" && car.Status != status {
			continue
		}
		// Search filter
		if search != "" && !strings.Contains(strings.ToLower(car.Brand+" "+car.Model), strings.ToLower(search)) {
			continue
		}
		filteredCars = append(filteredCars, *car)
	}

	// Sort cars
	switch sortBy {
	case "price_asc":
		// Sort by price low to high
		for i := 0; i < len(filteredCars)-1; i++ {
			for j := i + 1; j < len(filteredCars); j++ {
				if filteredCars[i].Price < filteredCars[j].Price {
					filteredCars[i], filteredCars[j] = filteredCars[j], filteredCars[i]
				}
			}
		}
	case "price_desc":
		// Sort by price high to low
		for i := 0; i < len(filteredCars)-1; i++ {
			for j := i + 1; j < len(filteredCars); j++ {
				if filteredCars[i].Price > filteredCars[j].Price {
					filteredCars[i], filteredCars[j] = filteredCars[j], filteredCars[i]
				}
			}
		}
	case "year_desc":
		// Sort by year newest to oldest
		for i := 0; i < len(filteredCars)-1; i++ {
			for j := i + 1; j < len(filteredCars); j++ {
				if filteredCars[i].Year > filteredCars[j].Year {
					filteredCars[i], filteredCars[j] = filteredCars[j], filteredCars[i]
				}
			}
		}
	default: // "newest" or empty
		// Sort by created_at newest to oldest
		for i := 0; i < len(filteredCars)-1; i++ {
			for j := i + 1; j < len(filteredCars); j++ {
				if filteredCars[i].CreatedAt.After(filteredCars[j].CreatedAt) {
					filteredCars[i], filteredCars[j] = filteredCars[j], filteredCars[i]
				}
			}
		}
	}

	// Get unique brands for filter dropdown
	brandMap := make(map[string]bool)
	var brands []string
	for _, car := range cars {
		if !brandMap[car.Brand] {
			brands = append(brands, car.Brand)
			brandMap[car.Brand] = true
		}
	}

	// Load photos for each car
	carsWithPhotos := make([]map[string]interface{}, len(filteredCars))
	for i, car := range filteredCars {
		photos, err := h.carRepo.GetCarPhotos(r.Context(), car.ID)
		if err != nil {
			slog.Error("failed to get car photos", "error", err, "car_id", car.ID)
			photos = []*model.CarPhoto{}
		}

		carsWithPhotos[i] = map[string]interface{}{
			"ID":           car.ID,
			"TenantID":     car.TenantID,
			"Brand":        car.Brand,
			"Model":        car.Model,
			"Year":         car.Year,
			"Price":        car.Price,
			"Mileage":      car.Mileage,
			"Transmission": car.Transmission,
			"FuelType":     car.FuelType,
			"EngineCC":     car.EngineCC,
			"Seats":        car.Seats,
			"Color":        car.Color,
			"Description":  car.Description,
			"Status":       car.Status,
			"IsFeatured":   car.IsFeatured,
			"CreatedAt":    car.CreatedAt,
			"UpdatedAt":    car.UpdatedAt,
			"Photos":       photos,
		}
	}

	// Prepare data for template
	data["Cars"] = carsWithPhotos
	data["Brands"] = brands
	data["SelectedBrand"] = brand
	data["SelectedStatus"] = status
	data["SortBy"] = sortBy
	data["SearchQuery"] = search

	if err := h.renderAdminPage(w, "templates/admin/cars.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminCarsNew renders the admin car creation form
func (h *PageHandler) AdminCarsNew(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "Tambah Mobil Baru"
	data["ActiveMenu"] = "cars"

	if err := h.renderAdminPage(w, "templates/admin/cars_new.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminCarsEdit renders the admin car edit form
func (h *PageHandler) AdminCarsEdit(w http.ResponseWriter, r *http.Request) {
	carIDStr := chi.URLParam(r, "id")

	// Convert car ID to int
	carID, err := strconv.Atoi(carIDStr)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	// Load car data
	car, err := h.carRepo.GetByID(r.Context(), carID)
	if err != nil {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	// Load car photos
	photos, err := h.carRepo.GetCarPhotos(r.Context(), carID)
	if err != nil {
		slog.Error("failed to get car photos", "error", err, "car_id", carID)
		photos = []*model.CarPhoto{}
	}

	data := h.getDefaultData(r)
	data["Title"] = "Edit Mobil"
	data["ActiveMenu"] = "cars"
	data["Car"] = car
	data["Photos"] = photos

	if err := h.renderAdminPage(w, "templates/admin/cars_edit.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminSales renders the admin sales management page
func (h *PageHandler) AdminSales(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "Manajemen Tim Sales"
	data["ActiveMenu"] = "sales"

	// Fetch sales list
	sales, err := h.salesRepo.List(r.Context())
	if err != nil {
		slog.Error("failed to list sales for admin page", "error", err)
		http.Error(w, "Gagal memuat data sales", http.StatusInternalServerError)
		return
	}
	data["Sales"] = sales

	if err := h.renderAdminPage(w, "templates/admin/sales.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminSalesTable renders just the sales table for HTMX updates
func (h *PageHandler) AdminSalesTable(w http.ResponseWriter, r *http.Request) {
	// Fetch sales list
	sales, err := h.salesRepo.List(r.Context())
	if err != nil {
		slog.Error("failed to list sales for table update", "error", err)
		http.Error(w, "Gagal memuat data sales", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Sales": sales,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "sales_table", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminWhatsApp renders the admin WhatsApp management page
func (h *PageHandler) AdminWhatsApp(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "WhatsApp Management"
	data["ActiveMenu"] = "whatsapp"

	// Get tenant ID from context
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		slog.Warn("tenant ID not found for WhatsApp page", "error", err)
		tenantID = 0 // Default, will be populated by middleware
	}
	data["TenantID"] = tenantID

	if err := h.renderAdminPage(w, "templates/admin/whatsapp_new.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminConversations renders the admin conversations list page
func (h *PageHandler) AdminConversations(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "Conversations"
	data["ActiveMenu"] = "conversations"

	// Initial empty state - HTMX will load actual data
	data["Conversations"] = nil
	data["Total"] = 0
	data["Page"] = 1
	data["Limit"] = 20
	data["TotalPages"] = 0

	if err := h.renderAdminPage(w, "templates/admin/conversations.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminConversationDetail renders the conversation detail page with message thread
func (h *PageHandler) AdminConversationDetail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := fmt.Sscanf(idStr, "%d", new(int))
	if err != nil || id != 1 {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	var conversationID int
	fmt.Sscanf(idStr, "%d", &conversationID)

	// Get conversation details
	conversation, err := h.conversationRepo.GetByID(r.Context(), conversationID)
	if err != nil {
		slog.Error("failed to get conversation", "error", err, "conversation_id", conversationID)
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	// Get messages (limit 50 for initial load)
	messages, err := h.conversationRepo.GetMessages(r.Context(), conversationID, 50)
	if err != nil {
		slog.Error("failed to get messages", "error", err, "conversation_id", conversationID)
		http.Error(w, "Failed to load messages", http.StatusInternalServerError)
		return
	}

	// Get total message count
	totalMessages, err := h.conversationRepo.GetMessageCount(r.Context(), conversationID)
	if err != nil {
		slog.Error("failed to count messages", "error", err, "conversation_id", conversationID)
		http.Error(w, "Failed to count messages", http.StatusInternalServerError)
		return
	}

	// Format conversation data
	convData := map[string]interface{}{
		"ID":          conversation.ID,
		"PhoneNumber": conversation.SenderPhone,
		"IsSales":     conversation.IsSales,
		"CreatedAt":   conversation.CreatedAt.Format("02 Jan 2006, 15:04"),
	}

	// Format messages data
	messageData := make([]map[string]interface{}, len(messages))
	for i, msg := range messages {
		messageData[i] = map[string]interface{}{
			"ID":        msg.ID,
			"Sender":    msg.SenderPhone,
			"Content":   msg.MessageText,
			"Direction": msg.Direction,
			"CreatedAt": msg.CreatedAt.Format("02 Jan 2006, 15:04"),
		}
	}

	data := h.getDefaultData(r)
	data["Title"] = "Detail Conversation"
	data["ActiveMenu"] = "conversations"
	data["Conversation"] = convData
	data["Messages"] = messageData
	data["TotalMessages"] = totalMessages

	if err := h.renderAdminPage(w, "templates/admin/conversation_detail.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminConversationsTable renders just the conversations table for HTMX updates
func (h *PageHandler) AdminConversationsTable(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page := 1
	limit := 20
	typeFilter := r.URL.Query().Get("type")

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := fmt.Sscanf(pageStr, "%d", &page); err == nil && p == 1 && page > 0 {
			// page is valid
		} else {
			page = 1
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := fmt.Sscanf(limitStr, "%d", &limit); err == nil && l == 1 && limit > 0 && limit <= 100 {
			// limit is valid
		} else {
			limit = 20
		}
	}

	// Validate type filter
	if typeFilter != "" && typeFilter != "customer" && typeFilter != "sales" && typeFilter != "all" {
		typeFilter = "all"
	}
	if typeFilter == "all" {
		typeFilter = ""
	}

	// Get conversations from repository
	conversations, total, err := h.conversationRepo.List(r.Context(), page, limit, typeFilter)
	if err != nil {
		slog.Error("failed to list conversations for table", "error", err)
		http.Error(w, "Gagal memuat conversations", http.StatusInternalServerError)
		return
	}

	// Calculate total pages
	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	data := map[string]interface{}{
		"Conversations": conversations,
		"Total":         total,
		"Page":          page,
		"Limit":         limit,
		"TotalPages":    totalPages,
	}

	// Parse and render just the table template
	tmpl, err := template.New("").Funcs(h.funcMap).ParseFiles("templates/admin/conversations.html")
	if err != nil {
		slog.Error("failed to parse conversations template", "error", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "conversations_table", data); err != nil {
		slog.Error("failed to render conversations table", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminAnalytics renders the admin analytics page
func (h *PageHandler) AdminAnalytics(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "Analytics"
	data["ActiveMenu"] = "analytics"

	if err := h.renderAdminPage(w, "templates/admin/analytics.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminSettings renders the admin settings page
func (h *PageHandler) AdminSettings(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "Pengaturan"
	data["ActiveMenu"] = "settings"

	if err := h.renderAdminPage(w, "templates/admin/settings.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// BlogList renders the public blog listing page
func (h *PageHandler) BlogList(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "Blog"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// BlogDetail renders the public blog detail page
func (h *PageHandler) BlogDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	data := h.getDefaultData(r)
	data["Title"] = "Blog Post"
	data["Slug"] = slug

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminBlog renders the admin blog management page
func (h *PageHandler) AdminBlog(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "Kelola Blog"
	data["ActiveMenu"] = "blog"

	if err := h.renderAdminPage(w, "templates/admin/blog.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminBlogNew renders the admin blog creation form
func (h *PageHandler) AdminBlogNew(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "Buat Post Blog Baru"
	data["ActiveMenu"] = "blog"
	data["IsEdit"] = false

	if err := h.renderAdminPage(w, "templates/admin/blog_form.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminBlogEdit renders the admin blog edit form
func (h *PageHandler) AdminBlogEdit(w http.ResponseWriter, r *http.Request) {
	blogID := chi.URLParam(r, "id")

	data := h.getDefaultData(r)
	data["Title"] = "Edit Post Blog"
	data["ActiveMenu"] = "blog"
	data["IsEdit"] = true
	data["BlogID"] = blogID

	if err := h.renderAdminPage(w, "templates/admin/blog_form.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminBranding renders the admin branding settings page
func (h *PageHandler) AdminBranding(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "Pengaturan Branding"
	data["ActiveMenu"] = "branding"

	if err := h.renderAdminPage(w, "templates/admin/branding.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminShowroom renders the admin showroom settings page
func (h *PageHandler) AdminShowroom(w http.ResponseWriter, r *http.Request) {
	data := h.getDefaultData(r)
	data["Title"] = "Informasi Showroom"
	data["ActiveMenu"] = "showroom"

	if err := h.renderAdminPage(w, "templates/admin/showroom.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// getDefaultData returns default template data with tenant info
func (h *PageHandler) getDefaultData(r *http.Request) map[string]interface{} {
	// Default values
	data := map[string]interface{}{
		"TenantName":     "Auto LMK",
		"WhatsAppNumber": "6281234567890",
		"Description":    "Platform penjualan mobil terpercaya",
		"Keywords":       "mobil, jual mobil, beli mobil",
		"FeaturedCars":   []interface{}{},
		"UserName":       "Admin",
		"LogoPath":       nil,
		"FaviconPath":    nil,
		"CustomTitle":    nil,
		"CustomSubtitle": nil,
		"PromoText":      nil,
	}

	// Try to load branding if tenant context is available
	tenantID, err := model.GetTenantID(r.Context())
	if err == nil && tenantID > 0 && h.brandingRepo != nil {
		branding, err := h.brandingRepo.GetByTenantID(r.Context(), tenantID)
		if err == nil && branding != nil {
			// Override with branding settings if available
			if branding.LogoPath != nil {
				data["LogoPath"] = *branding.LogoPath
			}
			if branding.FaviconPath != nil {
				data["FaviconPath"] = *branding.FaviconPath
			}
			if branding.CustomTitle != nil {
				data["CustomTitle"] = *branding.CustomTitle
			}
			if branding.CustomSubtitle != nil {
				data["CustomSubtitle"] = *branding.CustomSubtitle
			}
			if branding.PromoText != nil {
				data["PromoText"] = *branding.PromoText
			}
		}
	}

	// Try to load showroom settings if tenant context is available
	if tenantID > 0 && h.showroomRepo != nil {
		showroom, err := h.showroomRepo.GetByTenantID(r.Context(), tenantID)
		if err == nil && showroom != nil {
			// Add showroom data to template
			if showroom.Address != nil {
				data["ShowroomAddress"] = *showroom.Address
			}
			if showroom.Phone != nil {
				data["ShowroomPhone"] = *showroom.Phone
			}
			if showroom.Email != nil {
				data["ShowroomEmail"] = *showroom.Email
			}
			if showroom.BusinessHours != nil {
				data["ShowroomBusinessHours"] = *showroom.BusinessHours
			}
			if showroom.Latitude != nil {
				data["ShowroomLatitude"] = *showroom.Latitude
			}
			if showroom.Longitude != nil {
				data["ShowroomLongitude"] = *showroom.Longitude
			}
			if showroom.MapEmbed != nil {
				data["ShowroomMapEmbed"] = *showroom.MapEmbed
			}
		}
	}

	return data
}
