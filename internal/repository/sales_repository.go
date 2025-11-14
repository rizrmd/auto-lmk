package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/riz/auto-lmk/internal/model"
)

type SalesRepository struct {
	db *sql.DB
}

func NewSalesRepository(db *sql.DB) *SalesRepository {
	return &SalesRepository{db: db}
}

// Create registers a new sales person (tenant-scoped)
func (r *SalesRepository) Create(ctx context.Context, req *model.CreateSalesRequest) (*model.Sales, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, fmt.Errorf("tenant ID required: %w", err)
	}

	query := `
		INSERT INTO sales (tenant_id, phone_number, name)
		VALUES ($1, $2, $3)
		RETURNING id, tenant_id, phone_number, name, status, registered_at
	`

	sales := &model.Sales{}
	err = r.db.QueryRowContext(ctx, query, tenantID, req.PhoneNumber, req.Name).Scan(
		&sales.ID,
		&sales.TenantID,
		&sales.PhoneNumber,
		&sales.Name,
		&sales.Status,
		&sales.RegisteredAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create sales: %w", err)
	}

	return sales, nil
}

// GetByPhoneNumber retrieves sales by phone number (tenant-scoped)
func (r *SalesRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*model.Sales, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, fmt.Errorf("tenant ID required: %w", err)
	}

	query := `
		SELECT id, tenant_id, phone_number, name, status, registered_at
		FROM sales
		WHERE tenant_id = $1 AND phone_number = $2
	`

	sales := &model.Sales{}
	err = r.db.QueryRowContext(ctx, query, tenantID, phoneNumber).Scan(
		&sales.ID,
		&sales.TenantID,
		&sales.PhoneNumber,
		&sales.Name,
		&sales.Status,
		&sales.RegisteredAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sales not found")
		}
		return nil, fmt.Errorf("failed to get sales: %w", err)
	}

	return sales, nil
}

// IsSales checks if a phone number is registered as sales for a tenant
func (r *SalesRepository) IsSales(tenantID int, phoneNumber string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM sales
			WHERE tenant_id = $1 AND phone_number = $2 AND status = 'active'
		)
	`

	var exists bool
	err := r.db.QueryRow(query, tenantID, phoneNumber).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check sales status: %w", err)
	}

	return exists, nil
}

// List retrieves all sales for tenant
func (r *SalesRepository) List(ctx context.Context) ([]*model.Sales, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, fmt.Errorf("tenant ID required: %w", err)
	}

	query := `
		SELECT id, tenant_id, phone_number, name, status, registered_at
		FROM sales
		WHERE tenant_id = $1
		ORDER BY registered_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list sales: %w", err)
	}
	defer rows.Close()

	var salesList []*model.Sales
	for rows.Next() {
		sales := &model.Sales{}
		err := rows.Scan(
			&sales.ID,
			&sales.TenantID,
			&sales.PhoneNumber,
			&sales.Name,
			&sales.Status,
			&sales.RegisteredAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sales: %w", err)
		}
		salesList = append(salesList, sales)
	}

	return salesList, nil
}

// Delete removes a sales person (tenant-scoped)
func (r *SalesRepository) Delete(ctx context.Context, id int) error {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return fmt.Errorf("tenant ID required: %w", err)
	}

	query := "DELETE FROM sales WHERE id = $1 AND tenant_id = $2"
	result, err := r.db.ExecContext(ctx, query, id, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete sales: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("sales not found or no permission")
	}

	return nil
}
