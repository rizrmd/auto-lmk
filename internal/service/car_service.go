package service

import (
	"context"
	"fmt"

	"github.com/riz/auto-lmk/internal/model"
	"github.com/riz/auto-lmk/internal/repository"
)

// CarService handles business logic for cars
type CarService struct {
	carRepo   *repository.CarRepository
	photoRepo *repository.CarRepository // TODO: Create separate photo repository
}

func NewCarService(carRepo *repository.CarRepository) *CarService {
	return &CarService{
		carRepo: carRepo,
	}
}

// CreateCar creates a new car with validation
func (s *CarService) CreateCar(ctx context.Context, car *model.Car) error {
	// Validation
	if err := s.validateCar(car); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Set defaults
	if car.Status == "" {
		car.Status = "available"
	}

	return s.carRepo.Create(ctx, car)
}

// GetCarWithDetails gets car with photos and specs
func (s *CarService) GetCarWithDetails(ctx context.Context, id int) (*model.Car, error) {
	car, err := s.carRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// TODO: Load photos and specs
	// photos, _ := s.photoRepo.GetByCarID(ctx, id)
	// specs, _ := s.specRepo.GetByCarID(ctx, id)

	return car, nil
}

// SearchCarsForBot searches cars for WhatsApp bot with smart filtering
func (s *CarService) SearchCarsForBot(ctx context.Context, filters map[string]interface{}) ([]*model.Car, error) {
	// Smart filter parsing for natural language
	// e.g., "matic" -> transmission: automatic
	if trans, ok := filters["transmission"].(string); ok {
		if trans == "matic" {
			filters["transmission"] = "automatic"
		}
	}

	return s.carRepo.List(ctx, filters)
}

func (s *CarService) validateCar(car *model.Car) error {
	if car.Brand == "" {
		return fmt.Errorf("brand is required")
	}
	if car.Model == "" {
		return fmt.Errorf("model is required")
	}
	if car.Year < 1900 || car.Year > 2100 {
		return fmt.Errorf("invalid year")
	}
	if car.Price <= 0 {
		return fmt.Errorf("price must be positive")
	}
	return nil
}
