package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/riz/auto-lmk/internal/repository"
)

type AnalyticsHandler struct {
	analyticsRepo *repository.AnalyticsRepository
}

func NewAnalyticsHandler(analyticsRepo *repository.AnalyticsRepository) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsRepo: analyticsRepo}
}

// Dashboard shows analytics dashboard
func (h *AnalyticsHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	// Get top keywords
	topKeywords, err := h.analyticsRepo.GetTopSearchKeywords(r.Context(), 10)
	if err != nil {
		http.Error(w, "Failed to load analytics", http.StatusInternalServerError)
		return
	}

	// Get top viewed cars
	topCars, err := h.analyticsRepo.GetTopViewedCars(r.Context(), 10)
	if err != nil {
		http.Error(w, "Failed to load analytics", http.StatusInternalServerError)
		return
	}

	// Get search trends for last 30 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)
	trends, err := h.analyticsRepo.GetSearchTrends(r.Context(), startDate, endDate)
	if err != nil {
		http.Error(w, "Failed to load trends", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"TopKeywords": topKeywords,
		"TopCars":     topCars,
		"Trends":      trends,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetTopKeywords returns top search keywords as JSON
func (h *AnalyticsHandler) GetTopKeywords(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	keywords, err := h.analyticsRepo.GetTopSearchKeywords(r.Context(), limit)
	if err != nil {
		http.Error(w, "Failed to get keywords", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"keywords": keywords,
	})
}

// GetTopCars returns top viewed cars as JSON
func (h *AnalyticsHandler) GetTopCars(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	cars, err := h.analyticsRepo.GetTopViewedCars(r.Context(), limit)
	if err != nil {
		http.Error(w, "Failed to get cars", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"cars": cars,
	})
}

// GetTrends returns search trends as JSON
func (h *AnalyticsHandler) GetTrends(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	if startDateStr != "" {
		if d, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = d
		}
	}

	if endDateStr != "" {
		if d, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = d
		}
	}

	trends, err := h.analyticsRepo.GetSearchTrends(r.Context(), startDate, endDate)
	if err != nil {
		http.Error(w, "Failed to get trends", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"trends": trends,
	})
}

// ExportCSV exports analytics data as CSV file
func (h *AnalyticsHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	// Get top keywords (limit to 100)
	keywords, err := h.analyticsRepo.GetTopSearchKeywords(r.Context(), 100)
	if err != nil {
		http.Error(w, "Failed to get keywords", http.StatusInternalServerError)
		return
	}

	// Get top viewed cars (limit to 100)
	cars, err := h.analyticsRepo.GetTopViewedCars(r.Context(), 100)
	if err != nil {
		http.Error(w, "Failed to get cars", http.StatusInternalServerError)
		return
	}

	// Set CSV headers
	filename := "analytics_export_" + time.Now().Format("2006-01-02") + ".csv"
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

	// Write CSV content
	// Section 1: Search Keywords
	w.Write([]byte("Search Keywords Analytics\n"))
	w.Write([]byte("Keyword,Search Count,Last Searched At\n"))
	for _, keyword := range keywords {
		line := keyword.Keyword + "," +
			strconv.Itoa(keyword.SearchCount) + "," +
			keyword.LastSearchedAt.Format("2006-01-02 15:04:05") + "\n"
		w.Write([]byte(line))
	}

	// Section 2: Car Views
	w.Write([]byte("\n"))
	w.Write([]byte("Car Views Analytics\n"))
	w.Write([]byte("Car ID,Brand,Model,Year,View Count,Last Viewed At\n"))
	for _, car := range cars {
		line := strconv.Itoa(car.CarID) + "," +
			car.Brand + "," +
			car.Model + "," +
			strconv.Itoa(car.Year) + "," +
			strconv.Itoa(car.ViewCount) + "," +
			car.LastViewedAt.Format("2006-01-02 15:04:05") + "\n"
		w.Write([]byte(line))
	}
}
