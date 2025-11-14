package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/riz/auto-lmk/internal/model"
)

type CarRepository struct {
	db *sql.DB
}

func NewCarRepository(db *sql.DB) *CarRepository {
	return &CarRepository{db: db}
}

// Create creates a new car (tenant-scoped)
func (r *CarRepository) Create(ctx context.Context, car *model.Car) error {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return fmt.Errorf("tenant ID required: %w", err)
	}

	query := `
		INSERT INTO cars (
			tenant_id, brand, model, year, price, mileage, transmission,
			fuel_type, engine_cc, seats, color, description, status, is_featured
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at
	`

	err = r.db.QueryRowContext(ctx, query,
		tenantID, car.Brand, car.Model, car.Year, car.Price, car.Mileage,
		car.Transmission, car.FuelType, car.EngineCC, car.Seats, car.Color,
		car.Description, car.Status, car.IsFeatured,
	).Scan(&car.ID, &car.CreatedAt, &car.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create car: %w", err)
	}

	car.TenantID = tenantID
	return nil
}

// GetByID retrieves a car by ID (tenant-scoped)
func (r *CarRepository) GetByID(ctx context.Context, id int) (*model.Car, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, fmt.Errorf("tenant ID required: %w", err)
	}

	query := `
		SELECT id, tenant_id, brand, model, year, price, mileage, transmission,
			fuel_type, engine_cc, seats, color, description, status, is_featured,
			created_at, updated_at
		FROM cars
		WHERE id = $1 AND tenant_id = $2
	`

	car := &model.Car{}
	err = r.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&car.ID, &car.TenantID, &car.Brand, &car.Model, &car.Year, &car.Price,
		&car.Mileage, &car.Transmission, &car.FuelType, &car.EngineCC, &car.Seats,
		&car.Color, &car.Description, &car.Status, &car.IsFeatured,
		&car.CreatedAt, &car.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("car not found")
		}
		return nil, fmt.Errorf("failed to get car: %w", err)
	}

	return car, nil
}

// List retrieves all cars for tenant with filters
func (r *CarRepository) List(ctx context.Context, filters map[string]interface{}) ([]*model.Car, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, fmt.Errorf("tenant ID required: %w", err)
	}

	query := `
		SELECT id, tenant_id, brand, model, year, price, mileage, transmission,
			fuel_type, engine_cc, seats, color, description, status, is_featured,
			created_at, updated_at
		FROM cars
		WHERE tenant_id = $1
	`

	args := []interface{}{tenantID}
	argCount := 1

	// Add filters dynamically
	if brand, ok := filters["brand"].(string); ok && brand != "" {
		argCount++
		query += fmt.Sprintf(" AND brand = $%d", argCount)
		args = append(args, brand)
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
	}

	if maxPrice, ok := filters["max_price"].(int64); ok && maxPrice > 0 {
		argCount++
		query += fmt.Sprintf(" AND price <= $%d", argCount)
		args = append(args, maxPrice)
	}

	if transmission, ok := filters["transmission"].(string); ok && transmission != "" {
		argCount++
		query += fmt.Sprintf(" AND transmission = $%d", argCount)
		args = append(args, transmission)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list cars: %w", err)
	}
	defer rows.Close()

	var cars []*model.Car
	for rows.Next() {
		car := &model.Car{}
		err := rows.Scan(
			&car.ID, &car.TenantID, &car.Brand, &car.Model, &car.Year, &car.Price,
			&car.Mileage, &car.Transmission, &car.FuelType, &car.EngineCC, &car.Seats,
			&car.Color, &car.Description, &car.Status, &car.IsFeatured,
			&car.CreatedAt, &car.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan car: %w", err)
		}
		cars = append(cars, car)
	}

	return cars, nil
}

// Update updates a car (tenant-scoped)
func (r *CarRepository) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return fmt.Errorf("tenant ID required: %w", err)
	}

	// Build dynamic update query
	setClauses := []string{}
	args := []interface{}{}
	argCount := 0

	for key, value := range updates {
		argCount++
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", key, argCount))
		args = append(args, value)
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Add updated_at
	argCount++
	setClauses = append(setClauses, fmt.Sprintf("updated_at = CURRENT_TIMESTAMP"))

	// Add WHERE conditions
	argCount++
	args = append(args, id)
	idPlaceholder := argCount

	argCount++
	args = append(args, tenantID)
	tenantPlaceholder := argCount

	query := fmt.Sprintf(
		"UPDATE cars SET %s WHERE id = $%d AND tenant_id = $%d",
		strings.Join(setClauses, ", "),
		idPlaceholder,
		tenantPlaceholder,
	)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update car: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("car not found or no permission")
	}

	return nil
}

// Delete soft deletes a car (tenant-scoped)
func (r *CarRepository) Delete(ctx context.Context, id int) error {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return fmt.Errorf("tenant ID required: %w", err)
	}

	query := "DELETE FROM cars WHERE id = $1 AND tenant_id = $2"
	result, err := r.db.ExecContext(ctx, query, id, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete car: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("car not found or no permission")
	}

	return nil
}

// Search performs full-text search on cars (tenant-scoped)
func (r *CarRepository) Search(ctx context.Context, searchTerm string) ([]*model.Car, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, fmt.Errorf("tenant ID required: %w", err)
	}

	query := `
		SELECT id, tenant_id, brand, model, year, price, mileage, transmission,
			fuel_type, engine_cc, seats, color, description, status, is_featured,
			created_at, updated_at
		FROM cars
		WHERE tenant_id = $1
			AND status = 'available'
			AND (
				brand ILIKE $2
				OR model ILIKE $2
				OR description ILIKE $2
			)
		ORDER BY created_at DESC
	`

	searchPattern := "%" + searchTerm + "%"
	rows, err := r.db.QueryContext(ctx, query, tenantID, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search cars: %w", err)
	}
	defer rows.Close()

	var cars []*model.Car
	for rows.Next() {
		car := &model.Car{}
		err := rows.Scan(
			&car.ID, &car.TenantID, &car.Brand, &car.Model, &car.Year, &car.Price,
			&car.Mileage, &car.Transmission, &car.FuelType, &car.EngineCC, &car.Seats,
			&car.Color, &car.Description, &car.Status, &car.IsFeatured,
			&car.CreatedAt, &car.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan car: %w", err)
		}
		cars = append(cars, car)
	}

	return cars, nil
}
