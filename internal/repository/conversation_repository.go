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

// GetByID retrieves a conversation by ID with tenant isolation
func (r *ConversationRepository) GetByID(ctx context.Context, conversationID int) (*model.Conversation, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, fmt.Errorf("tenant ID required: %w", err)
	}

	query := `
		SELECT id, tenant_id, sender_phone, is_sales, created_at, updated_at
		FROM conversations
		WHERE id = $1 AND tenant_id = $2
	`

	conv := &model.Conversation{}
	err = r.db.QueryRowContext(ctx, query, conversationID, tenantID).Scan(
		&conv.ID, &conv.TenantID, &conv.SenderPhone, &conv.IsSales,
		&conv.CreatedAt, &conv.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return conv, nil
}

// GetMessageCount returns total count of messages in a conversation
func (r *ConversationRepository) GetMessageCount(ctx context.Context, conversationID int) (int, error) {
	query := `SELECT COUNT(*) FROM messages WHERE conversation_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, conversationID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count messages: %w", err)
	}

	return count, nil
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

// List retrieves paginated conversations with last message info
func (r *ConversationRepository) List(ctx context.Context, page, limit int, typeFilter string) ([]*model.ConversationListItem, int, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("tenant ID required: %w", err)
	}

	// Default pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Build WHERE clause
	whereClause := "WHERE c.tenant_id = $1"
	args := []interface{}{tenantID}
	argCount := 1

	if typeFilter == "customer" {
		argCount++
		whereClause += fmt.Sprintf(" AND c.is_sales = $%d", argCount)
		args = append(args, false)
	} else if typeFilter == "sales" {
		argCount++
		whereClause += fmt.Sprintf(" AND c.is_sales = $%d", argCount)
		args = append(args, true)
	}

	// Get total count
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM conversations c
		%s
	`, whereClause)

	var total int
	err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count conversations: %w", err)
	}

	// Get conversations with last message info
	query := fmt.Sprintf(`
		SELECT
			c.id,
			c.sender_phone,
			c.is_sales,
			c.created_at,
			COALESCE(last_msg.message_text, '') as last_message,
			COALESCE(last_msg.created_at, c.created_at) as last_message_at,
			COALESCE(msg_count.count, 0) as message_count
		FROM conversations c
		LEFT JOIN LATERAL (
			SELECT message_text, created_at
			FROM messages
			WHERE conversation_id = c.id
			ORDER BY created_at DESC
			LIMIT 1
		) last_msg ON true
		LEFT JOIN (
			SELECT conversation_id, COUNT(*) as count
			FROM messages
			GROUP BY conversation_id
		) msg_count ON msg_count.conversation_id = c.id
		%s
		ORDER BY last_message_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCount+1, argCount+2)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list conversations: %w", err)
	}
	defer rows.Close()

	var conversations []*model.ConversationListItem
	for rows.Next() {
		item := &model.ConversationListItem{}
		err := rows.Scan(
			&item.ID,
			&item.PhoneNumber,
			&item.IsSales,
			&item.CreatedAt,
			&item.LastMessage,
			&item.LastMessageAt,
			&item.MessageCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan conversation: %w", err)
		}
		conversations = append(conversations, item)
	}

	return conversations, total, nil
}
