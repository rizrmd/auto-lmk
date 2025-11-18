package repository

import (
	"context"
	"database/sql"

	"github.com/riz/auto-lmk/internal/model"
)

type ShowroomRepository struct {
	db *sql.DB
}

func NewShowroomRepository(db *sql.DB) *ShowroomRepository {
	return &ShowroomRepository{db: db}
}

// GetByTenantID retrieves showroom settings for a tenant
// Returns default empty settings if not found
func (r *ShowroomRepository) GetByTenantID(ctx context.Context, tenantID int) (*model.ShowroomSettings, error) {
	var showroom model.ShowroomSettings

	query := `
		SELECT
			tenant_id,
			address,
			phone,
			email,
			business_hours,
			latitude,
			longitude,
			map_embed,
			created_at,
			updated_at
		FROM showroom_settings
		WHERE tenant_id = $1
	`

	err := r.db.QueryRowContext(ctx, query, tenantID).Scan(
		&showroom.TenantID,
		&showroom.Address,
		&showroom.Phone,
		&showroom.Email,
		&showroom.BusinessHours,
		&showroom.Latitude,
		&showroom.Longitude,
		&showroom.MapEmbed,
		&showroom.CreatedAt,
		&showroom.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Return default empty settings
		return &model.ShowroomSettings{
			TenantID: tenantID,
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return &showroom, nil
}

// CreateOrUpdate creates or updates showroom settings using UPSERT pattern
func (r *ShowroomRepository) CreateOrUpdate(ctx context.Context, showroom *model.ShowroomSettings) error {
	query := `
		INSERT INTO showroom_settings (
			tenant_id,
			address,
			phone,
			email,
			business_hours,
			latitude,
			longitude,
			map_embed,
			created_at,
			updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()
		)
		ON CONFLICT (tenant_id)
		DO UPDATE SET
			address = EXCLUDED.address,
			phone = EXCLUDED.phone,
			email = EXCLUDED.email,
			business_hours = EXCLUDED.business_hours,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			map_embed = EXCLUDED.map_embed,
			updated_at = NOW()
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		showroom.TenantID,
		showroom.Address,
		showroom.Phone,
		showroom.Email,
		showroom.BusinessHours,
		showroom.Latitude,
		showroom.Longitude,
		showroom.MapEmbed,
	)

	return err
}
