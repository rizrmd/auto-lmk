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

	// Whitelist of allowed column names to prevent SQL injection
	allowedColumns := map[string]bool{
		"brand": true, "model": true, "year": true, "price": true, "color": true,
		"engine_cc": true, "mileage": true, "transmission": true, "fuel_type": true,
		"seats": true, "description": true, "status": true, "photos": true,
	}

	for key, value := range updates {
		if !allowedColumns[key] {
			return fmt.Errorf("invalid column name: %s", key)
		}
		argCount++
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", key, argCount))
		args = append(args, value)
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Add updated_at (no parameter needed, use direct expression)
	setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")

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

	// Sanitize search term to prevent SQL injection and ensure proper pattern matching
searchTerm = strings.TrimSpace(searchTerm)
if len(searchTerm) == 0 {
		return []*model.Car{}, nil
	}
	if len(searchTerm) > 100 {
		searchTerm = searchTerm[:100] // Limit search term length
	}
	searchPattern := "%" + strings.ReplaceAll(searchTerm, "%", "\\%") + "%"
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

// GetCarPhotos retrieves all photos for a car (tenant-scoped)
func (r *CarRepository) GetCarPhotos(ctx context.Context, carID int) ([]*model.CarPhoto, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, fmt.Errorf("tenant ID required: %w", err)
	}

	query := `
		SELECT id, car_id, file_path, display_order, created_at
		FROM car_photos
		WHERE car_id = $1 AND EXISTS (
			SELECT 1 FROM cars WHERE id = $1 AND tenant_id = $2
		)
		ORDER BY display_order ASC, id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, carID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get car photos: %w", err)
	}
	defer rows.Close()

	var photos []*model.CarPhoto
	for rows.Next() {
		photo := &model.CarPhoto{}
		err := rows.Scan(
			&photo.ID, &photo.CarID, &photo.FilePath, &photo.DisplayOrder, &photo.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan car photo: %w", err)
		}
		photos = append(photos, photo)
	}

	return photos, nil
}

// CreateWithPhotos creates a car with photos in a transaction (tenant-scoped)
func (r *CarRepository) CreateWithPhotos(ctx context.Context, car *model.Car, photoURLs []string) (int, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return 0, fmt.Errorf("tenant ID required: %w", err)
	}

	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create the car
	query := `
		INSERT INTO cars (
			tenant_id, brand, model, year, price, mileage, transmission,
			fuel_type, engine_cc, seats, color, description, status, is_featured
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`

	err = tx.QueryRowContext(ctx, query,
		tenantID, car.Brand, car.Model, car.Year, car.Price, car.Mileage,
		car.Transmission, car.FuelType, car.EngineCC, car.Seats, car.Color,
		car.Description, car.Status, car.IsFeatured,
	).Scan(&car.ID)

	if err != nil {
		return 0, fmt.Errorf("failed to create car: %w", err)
	}

	// Create photos if provided
	if len(photoURLs) > 0 {
		for i, photoURL := range photoURLs {
			if photoURL == "" {
				continue
			}

			photoQuery := `
				INSERT INTO car_photos (car_id, file_path, display_order)
				VALUES ($1, $2, $3)
			`

			_, err = tx.ExecContext(ctx, photoQuery, car.ID, photoURL, i+1)
			if err != nil {
				return 0, fmt.Errorf("failed to create car photo: %w", err)
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return car.ID, nil
}

// SearchWithFilters searches cars with filters (alias for List method)
func (r *CarRepository) SearchWithFilters(ctx context.Context, filters map[string]interface{}) ([]*model.Car, error) {
	return r.List(ctx, filters)
}

// AddPhotos adds photos to an existing car (tenant-scoped)
func (r *CarRepository) AddPhotos(ctx context.Context, carID int, photoURLs []string) error {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return fmt.Errorf("tenant ID required: %w", err)
	}

	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Verify car exists and belongs to tenant
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM cars WHERE id = $1 AND tenant_id = $2)"
	err = tx.QueryRowContext(ctx, checkQuery, carID, tenantID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to verify car ownership: %w", err)
	}
	if !exists {
		return fmt.Errorf("car not found or no permission")
	}

	// Get current max display order for this car
	var maxOrder int
	orderQuery := "SELECT COALESCE(MAX(display_order), 0) FROM car_photos WHERE car_id = $1"
	err = tx.QueryRowContext(ctx, orderQuery, carID).Scan(&maxOrder)
	if err != nil {
		return fmt.Errorf("failed to get max display order: %w", err)
	}

	// Insert photos
	for i, photoURL := range photoURLs {
		if photoURL == "" {
			continue
		}

		photoQuery := `
			INSERT INTO car_photos (car_id, file_path, display_order)
			VALUES ($1, $2, $3)
		`

		_, err = tx.ExecContext(ctx, photoQuery, carID, photoURL, maxOrder+i+1)
		if err != nil {
			return fmt.Errorf("failed to add car photo: %w", err)
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetPhotoByID retrieves a specific photo by ID with ownership verification
func (r *CarRepository) GetPhotoByID(ctx context.Context, photoID int) (*model.CarPhoto, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, fmt.Errorf("tenant ID required: %w", err)
	}

	query := `
		SELECT cp.id, cp.car_id, cp.file_path, cp.display_order, cp.created_at
		FROM car_photos cp
		INNER JOIN cars c ON cp.car_id = c.id
		WHERE cp.id = $1 AND c.tenant_id = $2
	`

	photo := &model.CarPhoto{}
	err = r.db.QueryRowContext(ctx, query, photoID, tenantID).Scan(
		&photo.ID, &photo.CarID, &photo.FilePath, &photo.DisplayOrder, &photo.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("photo not found")
		}
		return nil, fmt.Errorf("failed to get photo: %w", err)
	}

	return photo, nil
}

// DeletePhoto deletes a photo by ID with ownership verification
func (r *CarRepository) DeletePhoto(ctx context.Context, photoID int) error {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return fmt.Errorf("tenant ID required: %w", err)
	}

	query := `
		DELETE FROM car_photos
		WHERE id = $1 AND EXISTS (
			SELECT 1 FROM cars c
			WHERE c.id = car_photos.car_id AND c.tenant_id = $2
		)
	`

	result, err := r.db.ExecContext(ctx, query, photoID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete photo: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("photo not found or no permission")
	}

	return nil
}

// GetBrandsFromDB gets all unique brands from database for current tenant (AI generation)
func (r *CarRepository) GetBrandsFromDB(ctx context.Context) ([]string, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, fmt.Errorf("tenant ID required: %w", err)
	}

	query := `SELECT DISTINCT brand FROM cars WHERE tenant_id = $1 ORDER BY COUNT(*) DESC`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query brands: %w", err)
	}
	defer rows.Close()

	var brands []string
	for rows.Next() {
		var brand string
		if err := rows.Scan(&brand); err == nil {
			brands = append(brands, brand)
		}
	}
	return brands, nil
}

// GetModelsByBrandFromDB gets all models for a specific brand for current tenant (AI generation)
func (r *CarRepository) GetModelsByBrandFromDB(ctx context.Context, brand string) ([]string, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, fmt.Errorf("tenant ID required: %w", err)
	}

	query := `SELECT DISTINCT model FROM cars WHERE brand = $1 AND tenant_id = $2 ORDER BY COUNT(*) DESC`
	rows, err := r.db.QueryContext(ctx, query, brand, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query models: %w", err)
	}
	defer rows.Close()

	var models []string
	for rows.Next() {
		var model string
		if err := rows.Scan(&model); err == nil {
			models = append(models, model)
		}
	}
	return models, nil
}

// GetAllBrandModelsFromDB gets all brand-model combinations for current tenant (AI generation)
func (r *CarRepository) GetAllBrandModelsFromDB(ctx context.Context) ([]struct {
	Brand string
	Model string
}, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, fmt.Errorf("tenant ID required: %w", err)
	}

	query := `SELECT DISTINCT brand, model FROM cars WHERE tenant_id = $1 ORDER BY COUNT(*) DESC`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query brand-models: %w", err)
	}
	defer rows.Close()

	var results []struct {
		Brand string
		Model string
	}
	for rows.Next() {
		var brand, model string
		if err := rows.Scan(&brand, &model); err == nil {
			results = append(results, struct {
				Brand string
				Model string
			}{Brand: brand, Model: model})
		}
	}
	return results, nil
}

// GetDatabaseDefaultsFromDB gets average/common values from existing cars for current tenant (AI generation)
func (r *CarRepository) GetDatabaseDefaultsFromDB(ctx context.Context) (avgEngineCC, avgSeats, avgPrice float64, commonTransmission, commonFuelType string, err error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return 1500, 5, 150000000, "MT", "Bensin", fmt.Errorf("tenant ID required: %w", err)
	}

	query := `
		SELECT
			COALESCE(AVG(engine_cc), 1500) as avg_engine_cc,
			COALESCE(AVG(seats), 5) as avg_seats,
			MODE() WITHIN GROUP (ORDER BY transmission) as common_transmission,
			MODE() WITHIN GROUP (ORDER BY fuel_type) as common_fuel_type,
			COALESCE(AVG(price), 150000000) as avg_price
		FROM cars
		WHERE tenant_id = $1 AND status = 'available'
	`

	err = r.db.QueryRowContext(ctx, query, tenantID).Scan(&avgEngineCC, &avgSeats, &commonTransmission, &commonFuelType, &avgPrice)
	if err != nil {
		return 1500, 5, 150000000, "MT", "Bensin", fmt.Errorf("failed to get database defaults: %w", err)
	}

	// Handle null values
	if commonTransmission == "" {
		commonTransmission = "MT"
	}
	if commonFuelType == "" {
		commonFuelType = "Bensin"
	}

	return avgEngineCC, avgSeats, avgPrice, commonTransmission, commonFuelType, nil
}

// GetCarSpecsFromDB gets specifications for similar cars from database for current tenant (AI generation)
func (r *CarRepository) GetCarSpecsFromDB(ctx context.Context, brand, carModel string) (map[string]interface{}, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return map[string]interface{}{}, fmt.Errorf("tenant ID required: %w", err)
	}

	var query string
	var args []interface{}

	if brand != "" && carModel != "" {
		query = `
			SELECT AVG(engine_cc), AVG(seats), MODE() WITHIN GROUP (ORDER BY transmission),
			       MODE() WITHIN GROUP (ORDER BY fuel_type), AVG(price)
			FROM cars
			WHERE brand = $1 AND model = $2 AND tenant_id = $3 AND status = 'available'
		`
		args = []interface{}{brand, carModel, tenantID}
	} else if brand != "" {
		query = `
			SELECT AVG(engine_cc), AVG(seats), MODE() WITHIN GROUP (ORDER BY transmission),
			       MODE() WITHIN GROUP (ORDER BY fuel_type), AVG(price)
			FROM cars
			WHERE brand = $1 AND tenant_id = $2 AND status = 'available'
		`
		args = []interface{}{brand, tenantID}
	} else {
		return map[string]interface{}{}, nil
	}

	var commonTransmission, commonFuelType string
	var avgEngineCCPtr, avgSeatsPtr, avgPricePtr *float64

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&avgEngineCCPtr, &avgSeatsPtr, &commonTransmission, &commonFuelType, &avgPricePtr)
	if err != nil {
		return map[string]interface{}{}, fmt.Errorf("failed to get car specs: %w", err)
	}

	result := make(map[string]interface{})

	if avgEngineCCPtr != nil {
		result["engine_cc"] = int(*avgEngineCCPtr)
	}
	if avgSeatsPtr != nil {
		result["seats"] = int(*avgSeatsPtr)
	}
	if commonTransmission != "" {
		result["transmission"] = commonTransmission
	}
	if commonFuelType != "" {
		result["fuel_type"] = commonFuelType
	}
	if avgPricePtr != nil {
		result["price"] = int64(*avgPricePtr)
	}

	return result, nil
}