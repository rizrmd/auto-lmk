package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/riz/auto-lmk/internal/llm"
	"github.com/riz/auto-lmk/internal/model"
	"github.com/riz/auto-lmk/internal/repository"
)

type CarHandler struct {
	repo          *repository.CarRepository
	analyticsRepo *repository.AnalyticsRepository
	llmProvider   llm.Provider
}

func NewCarHandler(repo *repository.CarRepository) *CarHandler {
	return &CarHandler{repo: repo}
}

func NewCarHandlerWithAnalytics(repo *repository.CarRepository, analyticsRepo *repository.AnalyticsRepository) *CarHandler {
	return &CarHandler{
		repo:          repo,
		analyticsRepo: analyticsRepo,
	}
}

func NewCarHandlerWithLLM(repo *repository.CarRepository, llmProvider llm.Provider) *CarHandler {
	return &CarHandler{
		repo:        repo,
		llmProvider: llmProvider,
	}
}

func NewCarHandlerWithAnalyticsAndLLM(repo *repository.CarRepository, analyticsRepo *repository.AnalyticsRepository, llmProvider llm.Provider) *CarHandler {
	return &CarHandler{
		repo:          repo,
		analyticsRepo: analyticsRepo,
		llmProvider:   llmProvider,
	}
}

// Create handles POST /api/cars
func (h *CarHandler) Create(w http.ResponseWriter, r *http.Request) {
	var car model.Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		slog.Error("failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Default values
	if car.Status == "" {
		car.Status = "available"
	}

	if err := h.repo.Create(r.Context(), &car); err != nil {
		slog.Error("failed to create car", "error", err)
		http.Error(w, "Failed to create car", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(car)
}

// Get handles GET /api/cars/:id
func (h *CarHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	car, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		slog.Error("failed to get car", "error", err, "id", id)
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	// Log car view for analytics
	if h.analyticsRepo != nil {
		go h.logCarView(r.Context(), id)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(car)
}

// List handles GET /api/cars
func (h *CarHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filters
	filters := make(map[string]interface{})

	if brand := r.URL.Query().Get("brand"); brand != "" {
		filters["brand"] = brand
	}

	if status := r.URL.Query().Get("status"); status != "" {
		filters["status"] = status
	} else {
		filters["status"] = "available" // Default to available
	}

	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseInt(maxPriceStr, 10, 64); err == nil {
			filters["max_price"] = maxPrice
		}
	}

	if transmission := r.URL.Query().Get("transmission"); transmission != "" {
		filters["transmission"] = transmission
	}

	cars, err := h.repo.List(r.Context(), filters)
	if err != nil {
		slog.Error("failed to list cars", "error", err)
		http.Error(w, "Failed to list cars", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  cars,
		"count": len(cars),
	})
}

// Search handles GET /api/cars/search
func (h *CarHandler) Search(w http.ResponseWriter, r *http.Request) {
	// Parse search parameters
	filters := make(map[string]interface{})

	if query := r.URL.Query().Get("q"); query != "" {
		filters["query"] = query
	}

	if transmission := r.URL.Query().Get("transmission"); transmission != "" {
		filters["transmission"] = transmission
	}

	if fuelType := r.URL.Query().Get("fuel_type"); fuelType != "" {
		filters["fuel_type"] = fuelType
	}

	if minPriceStr := r.URL.Query().Get("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseInt(minPriceStr, 10, 64); err == nil {
			filters["min_price"] = minPrice
		}
	}

	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseInt(maxPriceStr, 10, 64); err == nil {
			filters["max_price"] = maxPrice
		}
	}

	// Default status to available
	filters["status"] = "available"

	// Require at least one search parameter
	if len(filters) <= 1 { // only status
		http.Error(w, "At least one search parameter required", http.StatusBadRequest)
		return
	}

	cars, err := h.repo.SearchWithFilters(r.Context(), filters)
	if err != nil {
		slog.Error("failed to search cars", "error", err, "filters", filters)
		// Return HTML error for HTMX requests
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, `<div class="text-center py-8">
				<p class="text-red-600">Terjadi kesalahan saat mencari mobil. Silakan coba lagi.</p>
				<p class="text-sm text-gray-500 mt-2">Error: %s</p>
			</div>`, err.Error())
			return
		}
		http.Error(w, "Failed to search cars", http.StatusInternalServerError)
		return
	}

	// Save search analytics
	if query, ok := filters["query"].(string); ok && query != "" {
		go h.saveSearchAnalytics(r.Context(), query, len(cars))
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// Return HTML partial for HTMX
		h.renderSearchResults(w, cars, filters)
		return
	}

	// Return JSON for API calls
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    cars,
		"count":   len(cars),
		"filters": filters,
	})
}

// renderSearchResults renders search results as HTML for HTMX
func (h *CarHandler) renderSearchResults(w http.ResponseWriter, cars []*model.Car, filters map[string]interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if len(cars) == 0 {
		fmt.Fprintf(w, `<div class="text-center py-8">
			<p class="text-gray-500">Tidak ditemukan mobil yang sesuai dengan kriteria pencarian.</p>
			<p class="text-sm text-gray-400 mt-2">Coba ubah filter atau kata kunci pencarian.</p>
		</div>`)
		return
	}

	fmt.Fprintf(w, `<div class="mb-4">
		<h3 class="text-lg font-semibold text-gray-900">Hasil Pencarian (%d mobil)</h3>
	</div>`, len(cars))

	fmt.Fprintf(w, `<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">`)

	for _, car := range cars {
		formattedPrice := fmt.Sprintf("Rp %s", formatPrice(car.Price))

		fmt.Fprintf(w, `<div class="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-xl transition-shadow">
			<div class="h-48 bg-gray-200 relative">
				<div class="flex items-center justify-center h-full text-gray-400">
					<svg class="w-16 h-16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
					</svg>
				</div>
			</div>
			<div class="p-4">
				<h3 class="text-xl font-bold text-gray-900">%s %s</h3>
				<p class="text-gray-600 text-sm">%d • %s • %s</p>
				<p class="text-2xl font-bold text-blue-600 mt-2">%s</p>
				<div class="mt-4 flex space-x-2">
					<a href="/mobil/%d" class="flex-1 bg-blue-500 hover:bg-blue-600 text-white text-center py-2 rounded">
						Lihat Detail
					</a>
				</div>
			</div>
		</div>`, car.Brand, car.Model, car.Year, getTransmissionDisplay(car.Transmission), getFuelTypeDisplay(car.FuelType), formattedPrice, car.ID)
	}

	fmt.Fprintf(w, `</div>`)
}

// Helper functions for display
func formatPrice(price int64) string {
	// Simple formatting for Indonesian Rupiah
	s := fmt.Sprintf("%d", price)
	n := len(s)
	if n <= 3 {
		return s
	}

	result := ""
	for i, digit := range s {
		if i > 0 && (n-i)%3 == 0 {
			result += "."
		}
		result += string(digit)
	}
	return result
}

func getTransmissionDisplay(transmission *string) string {
	if transmission == nil {
		return "Manual"
	}
	switch *transmission {
	case "AT":
		return "Automatic"
	case "MT":
		return "Manual"
	default:
		return *transmission
	}
}

func getFuelTypeDisplay(fuelType *string) string {
	if fuelType == nil {
		return "Bensin"
	}
	switch *fuelType {
	case "Bensin":
		return "Bensin"
	case "Diesel":
		return "Diesel"
	default:
		return *fuelType
	}
}

// Update handles PUT /api/cars/:id
func (h *CarHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		slog.Error("failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.Update(r.Context(), id, updates); err != nil {
		slog.Error("failed to update car", "error", err, "id", id)
		http.Error(w, "Failed to update car", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Car updated successfully"}`))
}

// Delete handles DELETE /api/cars/:id
func (h *CarHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		slog.Error("failed to delete car", "error", err, "id", id)
		http.Error(w, "Failed to delete car", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Car deleted successfully"}`))
}

// saveSearchAnalytics saves search keywords for analytics
func (h *CarHandler) saveSearchAnalytics(ctx context.Context, query string, resultCount int) {
	if h.analyticsRepo == nil {
		slog.Warn("analytics repository not available, skipping search analytics")
		return
	}

	// Only log if there's an actual search query
	if query != "" {
		if err := h.analyticsRepo.LogSearchEvent(ctx, query, resultCount); err != nil {
			slog.Error("failed to save search analytics", "error", err, "query", query)
		}
	}
}

// logCarView logs a car view event for analytics
func (h *CarHandler) logCarView(ctx context.Context, carID int) {
	if h.analyticsRepo == nil {
		slog.Warn("analytics repository not available, skipping car view analytics")
		return
	}

	if err := h.analyticsRepo.LogCarView(ctx, carID); err != nil {
		slog.Error("failed to log car view", "error", err, "car_id", carID)
	}
}

// AIGenerate handles POST /api/cars/ai-generate
// Generates car details from a text description using AI
func (h *CarHandler) AIGenerate(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var request struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.Error("failed to decode AI generate request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(request.Description) == "" {
		http.Error(w, "Description is required", http.StatusBadRequest)
		return
	}

	// Use database-driven AI generation (always available)
	slog.Info("Using database-driven AI generation", "description", request.Description)
	carData := h.GenerateCarDataFromDatabase(r.Context(), request.Description)
	slog.Info("Database-driven AI generation result", "data", carData)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(carData)
}

// generateMockCarData creates dynamic mock response based on description
func (h *CarHandler) generateMockCarData(description string) map[string]interface{} {
	desc := strings.ToLower(description)

	// Default values
	result := map[string]interface{}{
		"brand":        "Unknown",
		"model":        "Unknown",
		"year":         2020,
		"price":        150000000,
		"mileage":      50000,
		"transmission": "manual",
		"fuel_type":    "gasoline",
		"engine_cc":    1500,
		"seats":        5,
		"color":        "Unknown",
		"description":  description,
	}

	// Parse brand
	if strings.Contains(desc, "toyota") {
		result["brand"] = "Toyota"
		if strings.Contains(desc, "avanza") {
			result["model"] = "Avanza"
			result["engine_cc"] = 1496
			result["seats"] = 7
			result["price"] = 150000000
		} else if strings.Contains(desc, "yaris") {
			result["model"] = "Yaris"
			result["engine_cc"] = 1496
			result["seats"] = 5
			result["price"] = 180000000
		} else if strings.Contains(desc, "fortuner") {
			result["model"] = "Fortuner"
			result["engine_cc"] = 2393
			result["seats"] = 7
			result["price"] = 450000000
		} else if strings.Contains(desc, "rush") {
			result["model"] = "Rush"
			result["engine_cc"] = 1496
			result["seats"] = 7
			result["price"] = 200000000
		}
	} else if strings.Contains(strings.ToLower(desc), "honda") {
		result["brand"] = "Honda"
		if strings.Contains(strings.ToLower(desc), "br-v") || strings.Contains(strings.ToLower(desc), "brv") {
			result["model"] = "BR-V"
			result["engine_cc"] = 1497
			result["seats"] = 7
			result["price"] = 250000000
		} else if strings.Contains(strings.ToLower(desc), "jazz") {
			result["model"] = "Jazz"
			result["engine_cc"] = 1497
			result["seats"] = 5
			result["price"] = 180000000
		} else if strings.Contains(strings.ToLower(desc), "civic") {
			result["model"] = "Civic"
			result["engine_cc"] = 1498
			result["seats"] = 5
			result["price"] = 300000000
		} else if strings.Contains(strings.ToLower(desc), "hr-v") || strings.Contains(strings.ToLower(desc), "hrv") {
			result["model"] = "HR-V"
			result["engine_cc"] = 1497
			result["seats"] = 5
			result["price"] = 250000000
		} else if strings.Contains(strings.ToLower(desc), "cr-v") || strings.Contains(strings.ToLower(desc), "crv") {
			result["model"] = "CR-V"
			result["engine_cc"] = 1498
			result["seats"] = 7
			result["price"] = 350000000
		}
	} else if strings.Contains(desc, "suzuki") {
		result["brand"] = "Suzuki"
		if strings.Contains(desc, "ertiga") {
			result["model"] = "Ertiga"
			result["engine_cc"] = 1462
			result["seats"] = 7
			result["price"] = 160000000
		} else if strings.Contains(desc, "swift") {
			result["model"] = "Swift"
			result["engine_cc"] = 1197
			result["seats"] = 5
			result["price"] = 140000000
		} else if strings.Contains(desc, "baleno") {
			result["model"] = "Baleno"
			result["engine_cc"] = 1373
			result["seats"] = 5
			result["price"] = 150000000
		}
	} else if strings.Contains(desc, "daihatsu") {
		result["brand"] = "Daihatsu"
		if strings.Contains(desc, "xenia") {
			result["model"] = "Xenia"
			result["engine_cc"] = 1496
			result["seats"] = 7
			result["price"] = 140000000
		} else if strings.Contains(desc, "sigra") {
			result["model"] = "Sigra"
			result["engine_cc"] = 1197
			result["seats"] = 7
			result["price"] = 120000000
		} else if strings.Contains(desc, "ayla") {
			result["model"] = "Ayla"
			result["engine_cc"] = 998
			result["seats"] = 5
			result["price"] = 100000000
		}
	} else if strings.Contains(desc, "mitsubishi") {
		result["brand"] = "Mitsubishi"
		if strings.Contains(desc, "xpander") {
			result["model"] = "Xpander"
			result["engine_cc"] = 1499
			result["seats"] = 7
			result["price"] = 200000000
		} else if strings.Contains(desc, "pajero") {
			result["model"] = "Pajero"
			result["engine_cc"] = 2442
			result["seats"] = 7
			result["price"] = 500000000
		}
	} else {
		// Try to infer brand from model if brand not explicitly mentioned
		if strings.Contains(desc, "br-v") || strings.Contains(desc, "brv") {
			result["brand"] = "Honda"
			result["model"] = "BR-V"
			result["engine_cc"] = 1497
			result["seats"] = 7
			result["price"] = 250000000
		} else if strings.Contains(desc, "jazz") {
			result["brand"] = "Honda"
			result["model"] = "Jazz"
			result["engine_cc"] = 1497
			result["seats"] = 5
			result["price"] = 180000000
		} else if strings.Contains(desc, "avanza") {
			result["brand"] = "Toyota"
			result["model"] = "Avanza"
			result["engine_cc"] = 1496
			result["seats"] = 7
			result["price"] = 150000000
		}
	}

	// Parse year (look for 4-digit numbers between 1990-2030)
	reYear := regexp.MustCompile(`(19|20)\d{2}`)
	if matches := reYear.FindAllString(desc, -1); len(matches) > 0 {
		for _, match := range matches {
			if year, err := strconv.Atoi(match); err == nil && year >= 1990 && year <= 2030 {
				result["year"] = year
				break
			}
		}
	}

	// Parse price (look for patterns like "150 juta", "500 juta", etc.)
	rePrice := regexp.MustCompile(`(\d+)\s*(juta|jutaan|jt)`)
	if matches := rePrice.FindStringSubmatch(desc); len(matches) >= 3 {
		if amount, err := strconv.Atoi(matches[1]); err == nil {
			result["price"] = amount * 100000000
		}
	} else {
		// Try to parse "Rp 225.000.000" format
		rePriceRp := regexp.MustCompile(`rp\s*([\d\.]+)`)
		if matches := rePriceRp.FindStringSubmatch(desc); len(matches) >= 2 {
			priceStr := strings.ReplaceAll(matches[1], ".", "")
			if price, err := strconv.Atoi(priceStr); err == nil && price >= 1000000 && price <= 2000000000 {
				result["price"] = price
			}
		} else {
			// Alternative price parsing for any 3-9 digit numbers
			rePrice2 := regexp.MustCompile(`(\d{3,9})`)
			if matches := rePrice2.FindAllString(desc, -1); len(matches) > 0 {
				for _, match := range matches {
					if price, err := strconv.Atoi(match); err == nil && price >= 100000000 && price <= 2000000000 {
						result["price"] = price
						break
					}
				}
			}
		}
	}

	// Parse mileage (look for "km 50rb", "50000 km", "10.000 km", etc.)
	if strings.Contains(desc, "km") {
		// Try "15.000 km" format first (number then km)
		reMileageDots := regexp.MustCompile(`(\d[\d\.]*)\s*km`)
		if matches := reMileageDots.FindStringSubmatch(desc); len(matches) >= 2 {
			mileageStr := strings.ReplaceAll(matches[1], ".", "")
			if mileage, err := strconv.Atoi(mileageStr); err == nil && mileage >= 0 && mileage <= 500000 {
				result["mileage"] = mileage
			}
		} else {
			// Try "km 15.000" format (km then number)
			reMileageDotsReverse := regexp.MustCompile(`km\s*(\d[\d\.]*)`)
			if matches := reMileageDotsReverse.FindStringSubmatch(desc); len(matches) >= 2 {
				mileageStr := strings.ReplaceAll(matches[1], ".", "")
				if mileage, err := strconv.Atoi(mileageStr); err == nil && mileage >= 0 && mileage <= 500000 {
					result["mileage"] = mileage
				}
			} else {
				// Try "10rb km" or "10 ribu km" format
				reMileage := regexp.MustCompile(`(\d+)\s*(?:rb|ribu)?\s*km`)
				if matches := reMileage.FindStringSubmatch(desc); len(matches) >= 2 {
					if mileage, err := strconv.Atoi(matches[1]); err == nil {
						if strings.Contains(desc, "rb") || strings.Contains(desc, "ribu") {
							result["mileage"] = mileage * 1000
						} else {
							result["mileage"] = mileage
						}
					}
				} else {
					// Try "km 10rb" or "km 10 ribu" format
					reMileageReverse := regexp.MustCompile(`km\s*(\d+)\s*(?:rb|ribu)?`)
					if matches := reMileageReverse.FindStringSubmatch(desc); len(matches) >= 2 {
						if mileage, err := strconv.Atoi(matches[1]); err == nil {
							if strings.Contains(desc, "rb") || strings.Contains(desc, "ribu") {
								result["mileage"] = mileage * 1000
							} else {
								result["mileage"] = mileage
							}
						}
					}
				}
			}
		}
	}

	// Additional mileage parsing: try patterns without "km" but with "rb" or "ribu"
	if result["mileage"] == 50000 { // If still default, try other patterns
		if strings.Contains(desc, "rb") || strings.Contains(desc, "ribu") {
			reMileageAlt := regexp.MustCompile(`(\d+)\s*(?:rb|ribu)`)
			if matches := reMileageAlt.FindStringSubmatch(desc); len(matches) >= 2 {
				if mileage, err := strconv.Atoi(matches[1]); err == nil {
					result["mileage"] = mileage * 1000
				}
			}
		}
	}

	// Parse transmission
	if strings.Contains(desc, "matic") || strings.Contains(desc, "automatic") || strings.Contains(desc, "at") {
		result["transmission"] = "automatic"
	} else if strings.Contains(desc, "cvt") {
		result["transmission"] = "cvt"
	} else if strings.Contains(desc, "manual") || strings.Contains(desc, "mt") {
		result["transmission"] = "manual"
	}

	// Parse fuel type
	if strings.Contains(desc, "diesel") {
		result["fuel_type"] = "diesel"
	} else if strings.Contains(desc, "listrik") || strings.Contains(desc, "electric") {
		result["fuel_type"] = "electric"
	} else if strings.Contains(desc, "hybrid") {
		result["fuel_type"] = "hybrid"
	}

	// Parse color
	colors := []string{"hitam", "putih", "merah", "biru", "hijau", "kuning", "abu", "silver", "gray", "coklat", "orange", "ungu"}
	for _, color := range colors {
		if strings.Contains(desc, color) {
			result["color"] = strings.Title(color)
			break
		}
	}

	// Generate enhanced description
	brand := result["brand"].(string)
	model := result["model"].(string)
	year := result["year"].(int)
	price := result["price"].(int)
	mileage := result["mileage"].(int)
	transmission := result["transmission"].(string)

	priceInJuta := price / 100000000
	mileageStr := fmt.Sprintf("%d", mileage)
	if mileage >= 1000 {
		mileageStr = fmt.Sprintf("%d ribu", mileage/1000)
	}

	result["description"] = fmt.Sprintf("%s %s %d dengan kondisi sangat bagus. Mobil ini memiliki kilometer %s km, transmisi %s, dan harga %d juta rupiah. Fitur lengkap termasuk AC, power window, central lock, dan audio system. Kondisi mesin terawat dengan baik, tidak pernah mengalami tabrakan serius. Dokumen lengkap dan pajak panjang.",
		brand, model, year, mileageStr, transmission, priceInJuta)

	return result
}

// GenerateCarDataFromDatabase uses database data to drive car generation (no hardcode)
func (h *CarHandler) GenerateCarDataFromDatabase(ctx context.Context, description string) map[string]interface{} {
	slog.Info("Starting database-driven car data generation", "description", description)

	desc := strings.ToLower(description)

	// Start with database defaults
	result := h.getDatabaseDefaults(ctx)
	result["description"] = description

	// Parse brand and model from database patterns
	brand, model := h.parseBrandModelFromDatabase(ctx, desc)
	if brand != "" {
		result["brand"] = brand
	}
	if model != "" {
		result["model"] = model
	}

	// Get car specs from database if we found similar cars
	if brand != "" || model != "" {
		similarCarSpecs := h.getCarSpecsFromDatabase(ctx, brand, model)
		for key, value := range similarCarSpecs {
			// Only use database specs if user hasn't provided explicit values
			if result[key] == nil || result[key] == "Unknown" || result[key] == 0 {
				result[key] = value
			}
		}
	}

	// Parse explicit numeric values from user description (takes priority over database)
	h.parseExplicitValues(description, result)

	slog.Info("Database-driven car data generation completed", "result", result)
	return result
}

// getDatabaseDefaults gets reasonable defaults from existing database data
func (h *CarHandler) getDatabaseDefaults(ctx context.Context) map[string]interface{} {
	// Use repository method to get database defaults
	avgEngineCC, avgSeats, avgPrice, commonTransmission, commonFuelType, err := h.repo.GetDatabaseDefaultsFromDB(ctx)
	if err != nil {
		slog.Warn("Failed to get database defaults, using fallback values", "error", err)
		return map[string]interface{}{
			"brand":        "Unknown",
			"model":        "Unknown",
			"year":         2020,
			"price":        150000000,
			"mileage":      50000,
			"transmission": "MT",
			"fuel_type":    "Bensin",
			"engine_cc":    1500,
			"seats":        5,
			"color":        "Unknown",
		}
	}

	result := map[string]interface{}{
		"brand":        "Unknown",
		"model":        "Unknown",
		"year":         2020,
		"price":        int64(avgPrice),
		"mileage":      50000,
		"transmission": commonTransmission,
		"fuel_type":    commonFuelType,
		"engine_cc":    int(avgEngineCC),
		"seats":        int(avgSeats),
		"color":        "Unknown",
	}

	slog.Info("Database defaults retrieved", "result", result)
	return result
}

// normalizeBrandModel normalizes brand/model strings for better matching
func (h *CarHandler) normalizeBrandModel(input string) string {
	// Convert to lowercase
	normalized := strings.ToLower(input)

	// Replace hyphens with spaces
	normalized = strings.ReplaceAll(normalized, "-", " ")

	// Replace underscores with spaces
	normalized = strings.ReplaceAll(normalized, "_", " ")

	// Remove extra spaces
	words := strings.Fields(normalized)
	normalized = strings.Join(words, " ")

	return normalized
}

// containsNormalized checks if target contains pattern with normalization
func (h *CarHandler) containsNormalized(target, pattern string) bool {
	normalizedTarget := h.normalizeBrandModel(target)
	normalizedPattern := h.normalizeBrandModel(pattern)

	return strings.Contains(normalizedTarget, normalizedPattern)
}

// calculateSimilarity calculates a simple similarity score between two strings
func (h *CarHandler) calculateSimilarity(str1, str2 string) float64 {
	// Use Levenshtein-like approach with word-based comparison
	words1 := strings.Fields(h.normalizeBrandModel(str1))
	words2 := strings.Fields(h.normalizeBrandModel(str2))

	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}

	matches := 0
	for _, w1 := range words1 {
		for _, w2 := range words2 {
			if w1 == w2 || (len(w1) >= 3 && len(w2) >= 3 &&
				(strings.HasPrefix(w1, w2) || strings.HasPrefix(w2, w1))) {
				matches++
				break
			}
		}
	}

	// Calculate similarity as ratio of matches to total words
	maxWords := len(words1)
	if len(words2) > maxWords {
		maxWords = len(words2)
	}

	return float64(matches) / float64(maxWords)
}

// parseBrandModelFromDatabase searches for brand/model patterns in database with improved matching
func (h *CarHandler) parseBrandModelFromDatabase(ctx context.Context, desc string) (string, string) {
	// Normalize the description once for consistent processing
	normalizedDesc := h.normalizeBrandModel(desc)
	slog.Info("Normalized description for matching", "original", desc, "normalized", normalizedDesc)

	// Get all unique brands from database using repository
	brands, err := h.repo.GetBrandsFromDB(ctx)
	if err != nil {
		slog.Warn("Failed to query brands from database", "error", err)
		return "", ""
	}

	// BRAND-FIRST SEARCH: Search for brand mentions with improved matching
	var bestBrand string
	var bestBrandScore float64

	for _, brand := range brands {
		// Try exact normalized match first
		if h.containsNormalized(normalizedDesc, brand) {
			// Calculate similarity score
			score := h.calculateSimilarity(normalizedDesc, brand)
			slog.Info("Brand match found", "brand", brand, "score", score)

			if score > bestBrandScore {
				bestBrand = brand
				bestBrandScore = score
			}
		}
	}

	// If we found a good brand match (score >= 0.5 or exact match)
	if bestBrand != "" && bestBrandScore >= 0.5 {
		slog.Info("Best brand match selected", "brand", bestBrand, "score", bestBrandScore)

		// Now search for models within this brand
		models, err := h.repo.GetModelsByBrandFromDB(ctx, bestBrand)
		if err == nil {
			var bestModel string
			var bestModelScore float64

			// Search for model mentions with improved matching
			for _, model := range models {
				if h.containsNormalized(normalizedDesc, model) {
					score := h.calculateSimilarity(normalizedDesc, model)
					slog.Info("Model match found", "brand", bestBrand, "model", model, "score", score)

					if score > bestModelScore {
						bestModel = model
						bestModelScore = score
					}
				}
			}

			// Return model if we found a good match (score >= 0.4 for models, more lenient)
			if bestModel != "" && bestModelScore >= 0.4 {
				slog.Info("Best brand/model match found", "brand", bestBrand, "model", bestModel, "model_score", bestModelScore)
				return bestBrand, bestModel
			}

			// Return brand only if no good model match
			slog.Info("Brand match found without suitable model", "brand", bestBrand)
			return bestBrand, ""
		}
	}

	// MODEL-FIRST SEARCH: If no good brand match, try model-first search
	brandModels, err := h.repo.GetAllBrandModelsFromDB(ctx)
	if err != nil {
		slog.Warn("Failed to query brand-models from database", "error", err)
		return "", ""
	}

	var bestBrandModel string
	var bestModelScore float64
	var associatedBrand string

	for _, bm := range brandModels {
		if h.containsNormalized(normalizedDesc, bm.Model) {
			score := h.calculateSimilarity(normalizedDesc, bm.Model)
			slog.Info("Model-first search match", "brand", bm.Brand, "model", bm.Model, "score", score)

			if score > bestModelScore {
				bestBrandModel = bm.Model
				bestModelScore = score
				associatedBrand = bm.Brand
			}
		}
	}

	// Return the best model-first match if score is decent
	if bestBrandModel != "" && bestModelScore >= 0.4 {
		slog.Info("Best model-first match found", "brand", associatedBrand, "model", bestBrandModel, "score", bestModelScore)
		return associatedBrand, bestBrandModel
	}

	// PARTIAL MATCH FALLBACK: Try to find partial word matches for better typo handling
	descWords := strings.Fields(normalizedDesc)

	// Try each word as potential brand/model
	for _, word := range descWords {
		if len(word) < 3 {
			continue // Skip very short words
		}

		// Check if word matches any brand partially
		for _, brand := range brands {
			normalizedBrand := h.normalizeBrandModel(brand)
			brandWords := strings.Fields(normalizedBrand)

			for _, brandWord := range brandWords {
				if len(brandWord) >= 3 && (word == brandWord ||
					strings.HasPrefix(word, brandWord) ||
					strings.HasPrefix(brandWord, word)) {

					slog.Info("Partial brand match found", "word", word, "brand", brand)
					return brand, ""
				}
			}
		}
	}

	slog.Info("No brand/model matches found", "description", desc)
	return "", ""
}

// getCarSpecsFromDatabase gets specifications for similar cars from database
func (h *CarHandler) getCarSpecsFromDatabase(ctx context.Context, brand, model string) map[string]interface{} {
	slog.Info("Getting car specs from database", "brand", brand, "model", model)

	// Use repository method to get car specs
	result, err := h.repo.GetCarSpecsFromDB(ctx, brand, model)
	if err != nil {
		slog.Warn("Failed to get car specs from database", "error", err)
		return map[string]interface{}{}
	}

	slog.Info("Car specs retrieved from database", "result", result)
	return result
}

// IndonesianSynonymMap provides comprehensive synonym mapping for Indonesian car terms
var IndonesianSynonymMap = map[string][]string{
	"AT":      {"matic", "automatic", "automatik", "otomatis", "matik", "auto"},
	"MT":      {"manual", "mt"},
	"Bensin":  {"bensin", "premium", "pertamax", "pertalite"},
	"Diesel":  {"diesel", "solar"},
	"Listrik": {"listrik", "electric", "ev", "hybrid"},
	"engine_cc": {"mesin", "engine", "kapasitas", "cc"},
	"color": {"warna", "color", "cat"},
	"doors": {"pintu"},
	"seats": {"kursi", "seat", "bangku"},
}

// IndonesianNumberMultipliers handles Indonesian number formats
var IndonesianNumberMultipliers = map[string]int64{
	"ribu": 1000,
	"rb":   1000,
	"k":    1000,
	"juta": 1000000,
	"jt":   1000000,
	"m":    1000000,
	"miliar": 1000000000,
}

// Extended color names for Indonesian car market
var IndonesianColors = []string{
	"hitam", "putih", "merah", "biru", "hijau", "kuning",
	"abu-abu", "abu", "silver", "gray", "coklat", "orange",
	"ungu", "pink", "emas", "gold", "perak", "maroon", "beige",
	"cream", "navy", "toscana", "bronze", "champagne", "coklat muda",
	"biru dongker", "biru tua", "hijau tua", "merah marun",
}

// Enhanced transmission variants
var transmissionVariants = map[string]string{
	"matic":        "AT",
	"automatic":    "AT",
	"automatik":    "AT",
	"otomatis":     "AT",
	"matik":        "AT",
	"auto":         "AT",
	"at":           "AT",
	"cvt":          "CVT",
	"dsg":          "DSG",
	"tiptronic":    "AT",
	"manual":       "MT",
	"mt":           "MT",
	"manual transmission": "MT",
}

// Enhanced fuel type patterns
var fuelTypePatterns = map[string][]string{
	"Bensin":  {"bensin", "premium", "pertamax", "pertalite", "gasoline", "petrol"},
	"Diesel":  {"diesel", "solar", "diesel fuel"},
	"Listrik": {"listrik", "electric", "ev", "hybrid", "plug-in hybrid", "phev"},
}

// parseExplicitValues parses explicit numeric values from user description with Phase 2 enhancements
func (h *CarHandler) parseExplicitValues(description string, result map[string]interface{}) {
	desc := strings.ToLower(description)

	// Apply synonym normalization to the description
	desc = h.normalizeIndonesianSynonyms(desc)

	// Parse year with enhanced pattern recognition
	h.parseYear(desc, result)

	// Parse price with comprehensive Indonesian formats
	h.parsePrice(desc, result)

	// Parse mileage with Indonesian number formats
	h.parseMileage(desc, result)

	// Parse transmission with Indonesian variants
	h.parseTransmission(desc, result)

	// Parse fuel type with enhanced patterns
	h.parseFuelType(desc, result)

	// Parse engine CC with better recognition
	h.parseEngineCC(desc, result)

	// Parse color with extended color names
	h.parseColor(desc, result)

	// Parse seats and doors if mentioned
	h.parseSeatsAndDoors(desc, result)

	// Generate enhanced description
	h.generateEnhancedDescription(description, result)
}

// normalizeIndonesianSynonyms normalizes Indonesian car terms to standard form
func (h *CarHandler) normalizeIndonesianSynonyms(text string) string {
	normalized := text

	// Apply transmission normalization
	for variant, normalizedValue := range transmissionVariants {
		if strings.Contains(normalized, variant) {
			normalized = strings.ReplaceAll(normalized, variant, normalizedValue)
		}
	}

	// Apply other synonym mappings
	for standardValue, synonyms := range IndonesianSynonymMap {
		for _, synonym := range synonyms {
			if strings.Contains(normalized, synonym) && !strings.Contains(normalized, standardValue) {
				normalized = strings.ReplaceAll(normalized, synonym, standardValue)
			}
		}
	}

	return normalized
}

// parseYear extracts year information with fuzzy matching
func (h *CarHandler) parseYear(desc string, result map[string]interface{}) {
	// Enhanced year pattern with typo tolerance
	reYear := regexp.MustCompile(`(19|20)\d{2}`)
	if matches := reYear.FindAllString(desc, -1); len(matches) > 0 {
		for _, match := range matches {
			if year, err := strconv.Atoi(match); err == nil {
				// More flexible year range for edge cases
				if year >= 1985 && year <= 2035 {
					result["year"] = year
					break
				}
			}
		}
	}

	// Try to extract year from partial matches or typos
	if result["year"] == nil {
		// Handle common typos like "201o" (letter o instead of 0)
		reYearFuzzy := regexp.MustCompile(`20[1-2][0-9o]`)
		if matches := reYearFuzzy.FindAllString(desc, -1); len(matches) > 0 {
			for _, match := range matches {
				// Replace letter 'o' with '0'
				corrected := strings.ReplaceAll(match, "o", "0")
				if year, err := strconv.Atoi(corrected); err == nil && year >= 2010 && year <= 2029 {
					result["year"] = year
					break
				}
			}
		}
	}
}

// parsePrice handles comprehensive Indonesian price formats
func (h *CarHandler) parsePrice(desc string, result map[string]interface{}) {
	// Try juta/jutaan/jt formats first (most common in Indonesian)
	rePrice := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*(juta|jutaan|jt)`)
	if matches := rePrice.FindStringSubmatch(desc); len(matches) >= 3 {
		if amount, err := strconv.ParseFloat(strings.ReplaceAll(matches[1], ".", ""), 64); err == nil {
			result["price"] = int64(amount * 100000000)
			return
		}
	}

	// Try miliar format for luxury cars
	rePriceMiliar := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*(miliar)`)
	if matches := rePriceMiliar.FindStringSubmatch(desc); len(matches) >= 3 {
		if amount, err := strconv.ParseFloat(strings.ReplaceAll(matches[1], ".", ""), 64); err == nil {
			result["price"] = int64(amount * 1000000000)
			return
		}
	}

	// Try Rp format with dots
	rePriceRp := regexp.MustCompile(`rp\s*([\d\.]+(?:\.\d+)*)`)
	if matches := rePriceRp.FindStringSubmatch(desc); len(matches) >= 2 {
		priceStr := strings.ReplaceAll(matches[1], ".", "")
		if price, err := strconv.Atoi(priceStr); err == nil && price >= 1000000 && price <= 10000000000 {
			result["price"] = price
			return
		}
	}

	// Try prices with "ribu" for very low prices (unlikely but possible)
	rePriceRibu := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*(ribu|rb)`)
	if matches := rePriceRibu.FindStringSubmatch(desc); len(matches) >= 3 {
		if amount, err := strconv.ParseFloat(strings.ReplaceAll(matches[1], ".", ""), 64); err == nil {
			result["price"] = int64(amount * 1000000) // Convert ribu to base price
			return
		}
	}

	// Try any large numbers that could be prices
	rePriceNumbers := regexp.MustCompile(`(\d{7,12})`)
	if matches := rePriceNumbers.FindAllString(desc, -1); len(matches) > 0 {
		for _, match := range matches {
			if price, err := strconv.Atoi(match); err == nil && price >= 50000000 && price <= 10000000000 {
				result["price"] = price
				return
			}
		}
	}
}

// parseMileage handles Indonesian mileage formats with multipliers
func (h *CarHandler) parseMileage(desc string, result map[string]interface{}) {
	// Try "number km" format with dots
	reMileageDots := regexp.MustCompile(`(\d[\d\.]*)\s*km`)
	if matches := reMileageDots.FindStringSubmatch(desc); len(matches) >= 2 {
		mileageStr := strings.ReplaceAll(matches[1], ".", "")
		if mileage, err := strconv.Atoi(mileageStr); err == nil && mileage >= 0 && mileage <= 1000000 {
			result["mileage"] = mileage
			return
		}
	}

	// Try "km number" format
	reMileageReverse := regexp.MustCompile(`km\s*(\d[\d\.]*)`)
	if matches := reMileageReverse.FindStringSubmatch(desc); len(matches) >= 2 {
		mileageStr := strings.ReplaceAll(matches[1], ".", "")
		if mileage, err := strconv.Atoi(mileageStr); err == nil && mileage >= 0 && mileage <= 1000000 {
			result["mileage"] = mileage
			return
		}
	}

	// Try with Indonesian multipliers
	for multiplier, value := range IndonesianNumberMultipliers {
		// Try format like "15 ribu km" or "15rb km"
		reMultiplier := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*` + multiplier + `\s*km`)
		if matches := reMultiplier.FindStringSubmatch(desc); len(matches) >= 2 {
			if amount, err := strconv.ParseFloat(strings.ReplaceAll(matches[1], ".", ""), 64); err == nil {
				mileage := int64(amount * float64(value))
				if mileage >= 0 && mileage <= 1000000 {
					result["mileage"] = int(mileage)
					return
				}
			}
		}

		// Try format like "km 15 ribu" or "km 15rb"
		reMultiplierReverse := regexp.MustCompile(`km\s*(\d+(?:\.\d+)?)\s*` + multiplier)
		if matches := reMultiplierReverse.FindStringSubmatch(desc); len(matches) >= 2 {
			if amount, err := strconv.ParseFloat(strings.ReplaceAll(matches[1], ".", ""), 64); err == nil {
				mileage := int64(amount * float64(value))
				if mileage >= 0 && mileage <= 1000000 {
					result["mileage"] = int(mileage)
					return
				}
			}
		}
	}

	// Fallback: try to find mileage without explicit "km" but with multipliers
	for multiplier, value := range IndonesianNumberMultipliers {
		if value >= 1000 && value <= 1000000 { // Only use reasonable multipliers for mileage
			reMultiplierOnly := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*` + multiplier + `\b`)
			if matches := reMultiplierOnly.FindStringSubmatch(desc); len(matches) >= 2 {
				if amount, err := strconv.ParseFloat(strings.ReplaceAll(matches[1], ".", ""), 64); err == nil {
					mileage := int64(amount * float64(value))
					if mileage >= 0 && mileage <= 500000 { // Reasonable mileage range
						result["mileage"] = int(mileage)
						return
					}
				}
			}
		}
	}
}

// parseTransmission handles comprehensive transmission recognition
func (h *CarHandler) parseTransmission(desc string, result map[string]interface{}) {
	// Check each transmission variant
	for variant, standardValue := range transmissionVariants {
		if strings.Contains(desc, variant) {
			result["transmission"] = standardValue
			return
		}
	}

	// Handle combined patterns like "transmisi otomatis"
	if strings.Contains(desc, "transmisi") || strings.Contains(desc, "trans") {
		if strings.Contains(desc, "otomatis") || strings.Contains(desc, "automatic") {
			result["transmission"] = "AT"
		} else if strings.Contains(desc, "manual") {
			result["transmission"] = "MT"
		}
	}
}

// parseFuelType handles enhanced fuel type recognition
func (h *CarHandler) parseFuelType(desc string, result map[string]interface{}) {
	for fuelType, patterns := range fuelTypePatterns {
		for _, pattern := range patterns {
			if strings.Contains(desc, pattern) {
				result["fuel_type"] = fuelType
				return
			}
		}
	}

	// Handle hybrid/electric special cases
	if strings.Contains(desc, "hybrid") && strings.Contains(desc, "plug-in") {
		result["fuel_type"] = "Listrik" // PHEV treated as electric
	}
}

// parseEngineCC extracts engine displacement information
func (h *CarHandler) parseEngineCC(desc string, result map[string]interface{}) {
	// Look for CC patterns
	reCC := regexp.MustCompile(`(\d{3,4})\s*cc`)
	if matches := reCC.FindStringSubmatch(desc); len(matches) >= 2 {
		if cc, err := strconv.Atoi(matches[1]); err == nil && cc >= 600 && cc <= 8000 {
			result["engine_cc"] = cc
			return
		}
	}

	// Try patterns without explicit "cc" but in engine context
	if strings.Contains(desc, "mesin") || strings.Contains(desc, "engine") || strings.Contains(desc, "kapasitas") {
		reEngineSize := regexp.MustCompile(`(\d{3,4})`)
		if matches := reEngineSize.FindAllString(desc, -1); len(matches) > 0 {
			for _, match := range matches {
				if cc, err := strconv.Atoi(match); err == nil && cc >= 800 && cc <= 6000 {
					result["engine_cc"] = cc
					return
				}
			}
		}
	}
}

// parseColor handles comprehensive color recognition with Indonesian names
func (h *CarHandler) parseColor(desc string, result map[string]interface{}) {
	// Check for color context keywords first
	colorContext := strings.Contains(desc, "warna") || strings.Contains(desc, "color") || strings.Contains(desc, "cat")

	for _, color := range IndonesianColors {
		if strings.Contains(desc, color) {
			// For common multi-word colors, handle them properly
			formattedColor := strings.Title(color)
			if strings.Contains(color, " ") {
				words := strings.Split(color, " ")
				for i, word := range words {
					words[i] = strings.Title(word)
				}
				formattedColor = strings.Join(words, " ")
			}
			result["color"] = formattedColor
			return
		}
	}

	// If no color found but color context exists, try to infer from other words
	if colorContext {
		// Common adjectives that might describe color
		colorAdjectives := []string{"terang", "gelap", "muda", "tua", "netral"}
		for _, adj := range colorAdjectives {
			if strings.Contains(desc, adj) {
				result["color"] = strings.Title(adj)
				return
			}
		}
	}
}

// parseSeatsAndDoors extracts seats and doors information
func (h *CarHandler) parseSeatsAndDoors(desc string, result map[string]interface{}) {
	// Parse seats
	if strings.Contains(desc, "kursi") || strings.Contains(desc, "seat") || strings.Contains(desc, "bangku") {
		reSeats := regexp.MustCompile(`(\d+)\s*(?:kursi|seat|bangku)`)
		if matches := reSeats.FindStringSubmatch(desc); len(matches) >= 2 {
			if seats, err := strconv.Atoi(matches[1]); err == nil && seats >= 2 && seats <= 8 {
				result["seats"] = seats
			}
		}
	}

	// Parse doors
	if strings.Contains(desc, "pintu") {
		reDoors := regexp.MustCompile(`(\d+)\s*pintu`)
		if matches := reDoors.FindStringSubmatch(desc); len(matches) >= 2 {
			if doors, err := strconv.Atoi(matches[1]); err == nil && doors >= 2 && doors <= 5 {
				result["doors"] = doors
			}
		}
	}
}

// generateEnhancedDescription creates an enhanced Indonesian description
func (h *CarHandler) generateEnhancedDescription(originalDesc string, result map[string]interface{}) {
	brand, _ := result["brand"].(string)
	model, _ := result["model"].(string)
	year, _ := result["year"].(int)
	price, _ := result["price"].(int64)
	mileage, _ := result["mileage"].(int)
	transmission, _ := result["transmission"].(string)
	fuelType, _ := result["fuel_type"].(string)
	color, _ := result["color"].(string)
	engineCC, _ := result["engine_cc"].(int)
	seats, _ := result["seats"].(int)

	// Fallbacks for missing values
	if brand == "" { brand = "Unknown" }
	if model == "" { model = "Mobil" }
	if transmission == "" { transmission = "MT" }
	if fuelType == "" { fuelType = "Bensin" }
	if color == "" { color = "Warna Standar" }

	// Format price for Indonesian
	var priceStr string
	if price >= 1000000000 {
		priceInMiliar := float64(price) / 1000000000
		priceStr = fmt.Sprintf("%.1f miliar", priceInMiliar)
	} else {
		priceInJuta := float64(price) / 100000000
		priceStr = fmt.Sprintf("%.0f juta", priceInJuta)
	}

	// Format mileage
	var mileageStr string
	if mileage >= 1000 {
		mileageInRibu := mileage / 1000
		mileageStr = fmt.Sprintf("%d ribu", mileageInRibu)
	} else {
		mileageStr = fmt.Sprintf("%d", mileage)
	}

	// Build enhanced description with Indonesian formatting
	description := fmt.Sprintf("%s %s tahun %d %s. ", brand, model, year, color)
	description += fmt.Sprintf("Mobil dengan transmisi %s dan bahan bakar %s", transmission, fuelType)

	if engineCC > 0 {
		description += fmt.Sprintf(", kapasitas mesin %d cc", engineCC)
	}

	description += fmt.Sprintf(". Kilometer %s km, harga %s rupiah.", mileageStr, priceStr)

	if seats > 0 {
		description += fmt.Sprintf(" Kapasitas %d kursi.", seats)
	}

	description += " Kondisi sangat terawat, tidak pernah mengalami tabrakan serius. "
	description += "Fitur lengkap termasuk AC, power window, central lock, dan audio system. "
	description += "Dokumen lengkap dan pajak panjang."

	result["description"] = description
}

// UploadPhotos handles POST /api/cars/{id}/photos
// Uploads multiple photos for a car
func (h *CarHandler) UploadPhotos(w http.ResponseWriter, r *http.Request) {
	// Get car ID from URL
	carIDStr := chi.URLParam(r, "id")
	carID, err := strconv.Atoi(carIDStr)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	// Parse multipart form (max 32MB)
	err = r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		slog.Error("failed to parse multipart form", "error", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get files from form
	files := r.MultipartForm.File["photos"]
	if len(files) == 0 {
		http.Error(w, "No photos uploaded", http.StatusBadRequest)
		return
	}

	// Get display orders if provided
	displayOrders := r.MultipartForm.Value["display_orders"]

	// Process each photo
	var uploadedPhotos []model.CarPhoto
	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			slog.Error("failed to open uploaded file", "error", err, "filename", fileHeader.Filename)
			continue
		}
		defer file.Close()

		// Validate file type
		if !isValidImageType(fileHeader.Header.Get("Content-Type")) {
			slog.Warn("invalid file type", "filename", fileHeader.Filename, "content_type", fileHeader.Header.Get("Content-Type"))
			continue
		}

		// Generate unique filename
		ext := filepath.Ext(fileHeader.Filename)
		newFilename := fmt.Sprintf("car_%d_%d_%d%s", carID, time.Now().UnixNano(), i, ext)

		// Create uploads directory if it doesn't exist
		uploadsDir := "uploads/cars"
		if err := os.MkdirAll(uploadsDir, 0755); err != nil {
			slog.Error("failed to create uploads directory", "error", err)
			continue
		}

		// Save file
		filePath := filepath.Join(uploadsDir, newFilename)
		dst, err := os.Create(filePath)
		if err != nil {
			slog.Error("failed to create file", "error", err, "filepath", filePath)
			continue
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			slog.Error("failed to save file", "error", err, "filepath", filePath)
			continue
		}

		// Determine display order
		displayOrder := i
		if i < len(displayOrders) {
			if order, err := strconv.Atoi(displayOrders[i]); err == nil {
				displayOrder = order
			}
		}

		// Create photo record
		photo := model.CarPhoto{
			CarID:        carID,
			FilePath:     "/uploads/cars/" + newFilename,
			DisplayOrder: displayOrder,
		}
		uploadedPhotos = append(uploadedPhotos, photo)
	}

	// Extract file paths and save to database
	var photoURLs []string
	for _, photo := range uploadedPhotos {
		photoURLs = append(photoURLs, photo.FilePath)
	}

	// Save photos to database using existing AddPhotos method
	if err := h.repo.AddPhotos(r.Context(), carID, photoURLs); err != nil {
		slog.Error("failed to save photos to database", "error", err, "car_id", carID)
		http.Error(w, "Failed to save photos", http.StatusInternalServerError)
		return
	}

	slog.Info("photos uploaded successfully", "car_id", carID, "count", len(uploadedPhotos))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Photos uploaded successfully",
		"count":   len(uploadedPhotos),
	})
}

// DeletePhoto deletes a car photo
func (h *CarHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	photoIDStr := chi.URLParam(r, "photoId")
	photoID, err := strconv.Atoi(photoIDStr)
	if err != nil {
		http.Error(w, "Invalid photo ID", http.StatusBadRequest)
		return
	}

	// Get photo info before deletion
	photo, err := h.repo.GetPhotoByID(r.Context(), photoID)
	if err != nil {
		http.Error(w, "Photo not found", http.StatusNotFound)
		return
	}

	// Delete photo from database
	if err := h.repo.DeletePhoto(r.Context(), photoID); err != nil {
		slog.Error("failed to delete photo", "error", err, "photo_id", photoID)
		http.Error(w, "Failed to delete photo", http.StatusInternalServerError)
		return
	}

	// Delete physical file
	if photo.FilePath != "" {
		fullPath := filepath.Join("uploads", photo.FilePath)
		if err := os.Remove(fullPath); err != nil {
			// Log error but don't fail the request
			slog.Warn("failed to delete photo file", "error", err, "path", fullPath)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Photo deleted successfully",
	})
}

// isValidImageType checks if the content type is a valid image
func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/webp",
		"image/gif",
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}
