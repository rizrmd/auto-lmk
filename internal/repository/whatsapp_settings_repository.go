package repository

import (
	"context"
	"database/sql"

	"github.com/riz/auto-lmk/internal/model"
)

type WhatsAppSettingsRepository struct {
	db *sql.DB
}

func NewWhatsAppSettingsRepository(db *sql.DB) *WhatsAppSettingsRepository {
	return &WhatsAppSettingsRepository{db: db}
}

// GetSettings retrieves WhatsApp settings for a tenant
func (r *WhatsAppSettingsRepository) GetSettings(ctx context.Context, tenantID int) (*model.WhatsAppSettings, error) {
	query := `
		SELECT tenant_id, bot_number, fallback_number, created_at, updated_at
		FROM whatsapp_settings
		WHERE tenant_id = $1
	`

	settings := &model.WhatsAppSettings{}
	err := r.db.QueryRowContext(ctx, query, tenantID).Scan(
		&settings.TenantID,
		&settings.BotNumber,
		&settings.FallbackNumber,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Return default settings if not found
			return &model.WhatsAppSettings{
				TenantID:       tenantID,
				FallbackNumber: "+6281234567890", // Default fallback
			}, nil
		}
		return nil, err
	}

	return settings, nil
}

// SaveSettings saves or updates WhatsApp settings for a tenant
func (r *WhatsAppSettingsRepository) SaveSettings(ctx context.Context, settings *model.WhatsAppSettings) error {
	query := `
		INSERT INTO whatsapp_settings (tenant_id, bot_number, fallback_number, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT (tenant_id)
		DO UPDATE SET
			bot_number = EXCLUDED.bot_number,
			fallback_number = EXCLUDED.fallback_number,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := r.db.ExecContext(ctx, query,
		settings.TenantID,
		settings.BotNumber,
		settings.FallbackNumber,
	)

	return err
}
