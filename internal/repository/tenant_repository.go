package repository

import (
	"database/sql"
	"fmt"

	"github.com/riz/auto-lmk/internal/model"
)

type TenantRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) Create(req *model.CreateTenantRequest) (*model.Tenant, error) {
	query := `
		INSERT INTO tenants (domain, name, whatsapp_number)
		VALUES ($1, $2, $3)
		RETURNING id, domain, name, whatsapp_number, pairing_status, status, created_at, updated_at
	`

	tenant := &model.Tenant{}
	err := r.db.QueryRow(query, req.Domain, req.Name, req.WhatsAppNumber).Scan(
		&tenant.ID,
		&tenant.Domain,
		&tenant.Name,
		&tenant.WhatsAppNumber,
		&tenant.PairingStatus,
		&tenant.Status,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	return tenant, nil
}

func (r *TenantRepository) GetByID(id int) (*model.Tenant, error) {
	query := `
		SELECT id, domain, name, whatsapp_number, pairing_status, status, created_at, updated_at
		FROM tenants
		WHERE id = $1
	`

	tenant := &model.Tenant{}
	err := r.db.QueryRow(query, id).Scan(
		&tenant.ID,
		&tenant.Domain,
		&tenant.Name,
		&tenant.WhatsAppNumber,
		&tenant.PairingStatus,
		&tenant.Status,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tenant not found")
		}
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	return tenant, nil
}

func (r *TenantRepository) GetByDomain(domain string) (*model.Tenant, error) {
	query := `
		SELECT id, domain, name, whatsapp_number, pairing_status, status, created_at, updated_at
		FROM tenants
		WHERE domain = $1
	`

	tenant := &model.Tenant{}
	err := r.db.QueryRow(query, domain).Scan(
		&tenant.ID,
		&tenant.Domain,
		&tenant.Name,
		&tenant.WhatsAppNumber,
		&tenant.PairingStatus,
		&tenant.Status,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tenant not found")
		}
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	return tenant, nil
}

func (r *TenantRepository) List() ([]*model.Tenant, error) {
	query := `
		SELECT id, domain, name, whatsapp_number, pairing_status, status, created_at, updated_at
		FROM tenants
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}
	defer rows.Close()

	var tenants []*model.Tenant
	for rows.Next() {
		tenant := &model.Tenant{}
		err := rows.Scan(
			&tenant.ID,
			&tenant.Domain,
			&tenant.Name,
			&tenant.WhatsAppNumber,
			&tenant.PairingStatus,
			&tenant.Status,
			&tenant.CreatedAt,
			&tenant.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tenant: %w", err)
		}
		tenants = append(tenants, tenant)
	}

	return tenants, nil
}

// UpdatePairingStatus updates the pairing_status field for a tenant
// Valid status values: "unpaired", "pairing_pending", "paired", "disconnected", "failed"
func (r *TenantRepository) UpdatePairingStatus(tenantID int, status string) error {
	query := `
		UPDATE tenants
		SET pairing_status = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.Exec(query, status, tenantID)
	if err != nil {
		return fmt.Errorf("failed to update pairing status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tenant not found with id %d", tenantID)
	}

	return nil
}

// UpdateWhatsAppNumber updates the whatsapp_number field for a tenant
func (r *TenantRepository) UpdateWhatsAppNumber(tenantID int, phoneNumber string) error {
	query := `
		UPDATE tenants
		SET whatsapp_number = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.Exec(query, phoneNumber, tenantID)
	if err != nil {
		return fmt.Errorf("failed to update whatsapp number: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tenant not found with id %d", tenantID)
	}

	return nil
}
