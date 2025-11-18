package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/riz/auto-lmk/internal/model"
)

type BrandingRepository struct {
	db *sql.DB
}

func NewBrandingRepository(db *sql.DB) *BrandingRepository {
	return &BrandingRepository{db: db}
}

// GetByTenantID retrieves branding settings for a specific tenant
func (r *BrandingRepository) GetByTenantID(ctx context.Context, tenantID int) (*model.BrandingSettings, error) {
	query := `
		SELECT tenant_id, logo_path, favicon_path, custom_title, custom_subtitle,
		       promo_text, header_style, created_at, updated_at
		FROM tenant_branding
		WHERE tenant_id = $1
	`

	var branding model.BrandingSettings
	err := r.db.QueryRowContext(ctx, query, tenantID).Scan(
		&branding.TenantID,
		&branding.LogoPath,
		&branding.FaviconPath,
		&branding.CustomTitle,
		&branding.CustomSubtitle,
		&branding.PromoText,
		&branding.HeaderStyle,
		&branding.CreatedAt,
		&branding.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Return default branding if not found
		return &model.BrandingSettings{
			TenantID:    tenantID,
			HeaderStyle: "default",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return &branding, nil
}

// CreateOrUpdate creates or updates branding settings for a tenant
func (r *BrandingRepository) CreateOrUpdate(ctx context.Context, branding *model.BrandingSettings) error {
	query := `
		INSERT INTO tenant_branding (
			tenant_id, logo_path, favicon_path, custom_title, custom_subtitle,
			promo_text, header_style, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (tenant_id) DO UPDATE SET
			logo_path = EXCLUDED.logo_path,
			favicon_path = EXCLUDED.favicon_path,
			custom_title = EXCLUDED.custom_title,
			custom_subtitle = EXCLUDED.custom_subtitle,
			promo_text = EXCLUDED.promo_text,
			header_style = EXCLUDED.header_style,
			updated_at = EXCLUDED.updated_at
	`

	branding.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		branding.TenantID,
		branding.LogoPath,
		branding.FaviconPath,
		branding.CustomTitle,
		branding.CustomSubtitle,
		branding.PromoText,
		branding.HeaderStyle,
		branding.UpdatedAt,
	)

	return err
}

// UpdateLogo updates only the logo path for a tenant
func (r *BrandingRepository) UpdateLogo(ctx context.Context, tenantID int, logoPath string) error {
	query := `
		INSERT INTO tenant_branding (tenant_id, logo_path, header_style, updated_at)
		VALUES ($1, $2, 'default', $3)
		ON CONFLICT (tenant_id) DO UPDATE SET
			logo_path = EXCLUDED.logo_path,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query, tenantID, logoPath, time.Now())
	return err
}

// UpdateFavicon updates only the favicon path for a tenant
func (r *BrandingRepository) UpdateFavicon(ctx context.Context, tenantID int, faviconPath string) error {
	query := `
		INSERT INTO tenant_branding (tenant_id, favicon_path, header_style, updated_at)
		VALUES ($1, $2, 'default', $3)
		ON CONFLICT (tenant_id) DO UPDATE SET
			favicon_path = EXCLUDED.favicon_path,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query, tenantID, faviconPath, time.Now())
	return err
}
