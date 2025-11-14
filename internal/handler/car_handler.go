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

type CarHandler struct {
	repo *repository.CarRepository
}

func NewCarHandler(repo *repository.CarRepository) *CarHandler {
	return &CarHandler{repo: repo}
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
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query required", http.StatusBadRequest)
		return
	}

	cars, err := h.repo.Search(r.Context(), query)
	if err != nil {
		slog.Error("failed to search cars", "error", err)
		http.Error(w, "Failed to search cars", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  cars,
		"count": len(cars),
		"query": query,
	})
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
