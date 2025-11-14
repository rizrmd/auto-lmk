package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/riz/auto-lmk/internal/model"
)

type LeadRepository struct {
	db *sql.DB
}

func NewLeadRepository(db *sql.DB) *LeadRepository {
	return &LeadRepository{db: db}
}

// Create creates a new lead (tenant-scoped)
func (r *LeadRepository) Create(ctx context.Context, req *model.CreateLeadRequest) (*model.Lead, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, fmt.Errorf("tenant ID required: %w", err)
	}

	query := `
		INSERT INTO leads (tenant_id, phone_number, name, interested_car_id, conversation_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, tenant_id, phone_number, name, interested_car_id, conversation_id, status, created_at, updated_at
	`

	lead := &model.Lead{}
	err = r.db.QueryRowContext(ctx, query, tenantID, req.PhoneNumber, req.Name, req.InterestedCarID, req.ConversationID).Scan(
		&lead.ID, &lead.TenantID, &lead.PhoneNumber, &lead.Name,
		&lead.InterestedCarID, &lead.ConversationID, &lead.Status,
		&lead.CreatedAt, &lead.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create lead: %w", err)
	}

	return lead, nil
}

// List retrieves all leads for tenant
func (r *LeadRepository) List(ctx context.Context, status string) ([]*model.Lead, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, fmt.Errorf("tenant ID required: %w", err)
	}

	query := `
		SELECT id, tenant_id, phone_number, name, interested_car_id, conversation_id, status, created_at, updated_at
		FROM leads
		WHERE tenant_id = $1
	`

	args := []interface{}{tenantID}

	if status != "" {
		query += " AND status = $2"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list leads: %w", err)
	}
	defer rows.Close()

	var leads []*model.Lead
	for rows.Next() {
		lead := &model.Lead{}
		err := rows.Scan(
			&lead.ID, &lead.TenantID, &lead.PhoneNumber, &lead.Name,
			&lead.InterestedCarID, &lead.ConversationID, &lead.Status,
			&lead.CreatedAt, &lead.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lead: %w", err)
		}
		leads = append(leads, lead)
	}

	return leads, nil
}

// UpdateStatus updates lead status
func (r *LeadRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return fmt.Errorf("tenant ID required: %w", err)
	}

	query := "UPDATE leads SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND tenant_id = $3"
	result, err := r.db.ExecContext(ctx, query, status, id, tenantID)
	if err != nil {
		return fmt.Errorf("failed to update lead status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("lead not found or no permission")
	}

	return nil
}
