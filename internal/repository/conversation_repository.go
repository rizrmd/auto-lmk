package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/riz/auto-lmk/internal/model"
)

type ConversationRepository struct {
	db *sql.DB
}

func NewConversationRepository(db *sql.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

// GetOrCreate finds existing conversation or creates new one
func (r *ConversationRepository) GetOrCreate(ctx context.Context, senderPhone string, isSales bool) (*model.Conversation, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, fmt.Errorf("tenant ID required: %w", err)
	}

	// Try to get existing conversation
	query := `
		SELECT id, tenant_id, sender_phone, is_sales, created_at, updated_at
		FROM conversations
		WHERE tenant_id = $1 AND sender_phone = $2
		ORDER BY updated_at DESC
		LIMIT 1
	`

	conv := &model.Conversation{}
	err = r.db.QueryRowContext(ctx, query, tenantID, senderPhone).Scan(
		&conv.ID, &conv.TenantID, &conv.SenderPhone, &conv.IsSales,
		&conv.CreatedAt, &conv.UpdatedAt,
	)

	if err == nil {
		// Update timestamp
		r.db.ExecContext(ctx, "UPDATE conversations SET updated_at = CURRENT_TIMESTAMP WHERE id = $1", conv.ID)
		return conv, nil
	}

	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// Create new conversation
	insertQuery := `
		INSERT INTO conversations (tenant_id, sender_phone, is_sales)
		VALUES ($1, $2, $3)
		RETURNING id, tenant_id, sender_phone, is_sales, created_at, updated_at
	`

	err = r.db.QueryRowContext(ctx, insertQuery, tenantID, senderPhone, isSales).Scan(
		&conv.ID, &conv.TenantID, &conv.SenderPhone, &conv.IsSales,
		&conv.CreatedAt, &conv.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return conv, nil
}

// AddMessage adds a message to conversation
func (r *ConversationRepository) AddMessage(ctx context.Context, conversationID int, senderPhone, messageText, direction string) error {
	query := `
		INSERT INTO messages (conversation_id, sender_phone, message_text, direction)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query, conversationID, senderPhone, messageText, direction)
	if err != nil {
		return fmt.Errorf("failed to add message: %w", err)
	}

	return nil
}

// GetMessages retrieves messages for a conversation
func (r *ConversationRepository) GetMessages(ctx context.Context, conversationID int, limit int) ([]*model.Message, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
		SELECT id, conversation_id, sender_phone, message_text, direction, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, conversationID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		msg := &model.Message{}
		err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.SenderPhone, &msg.MessageText, &msg.Direction, &msg.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// ListConversations lists all conversations for tenant
func (r *ConversationRepository) ListConversations(ctx context.Context) ([]*model.Conversation, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, fmt.Errorf("tenant ID required: %w", err)
	}

	query := `
		SELECT id, tenant_id, sender_phone, is_sales, created_at, updated_at
		FROM conversations
		WHERE tenant_id = $1
		ORDER BY updated_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}
	defer rows.Close()

	var conversations []*model.Conversation
	for rows.Next() {
		conv := &model.Conversation{}
		err := rows.Scan(&conv.ID, &conv.TenantID, &conv.SenderPhone, &conv.IsSales, &conv.CreatedAt, &conv.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}
