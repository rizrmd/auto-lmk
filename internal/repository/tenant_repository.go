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
