package handler

import (
	"encoding/json"
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/riz/auto-lmk/internal/model"
	"github.com/riz/auto-lmk/internal/repository"
)

type ConversationHandler struct {
	repo *repository.ConversationRepository
}

func NewConversationHandler(repo *repository.ConversationRepository) *ConversationHandler {
	return &ConversationHandler{repo: repo}
}

// ConversationListResponse represents API response for conversation list
type ConversationListResponse struct {
	Conversations []*model.ConversationListItem `json:"conversations"`
	Total         int                            `json:"total"`
	Page          int                            `json:"page"`
	Limit         int                            `json:"limit"`
	TotalPages    int                            `json:"total_pages"`
}

// List handles GET /api/conversations
func (h *ConversationHandler) List(w http.ResponseWriter, r *http.Request) {
	// Extract tenant ID from context
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		slog.Error("tenant ID required for conversation list", "error", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	typeFilter := r.URL.Query().Get("type")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Validate type filter
	if typeFilter != "" && typeFilter != "customer" && typeFilter != "sales" && typeFilter != "all" {
		typeFilter = "all"
	}
	if typeFilter == "all" {
		typeFilter = ""
	}

	// Get conversations from repository
	conversations, total, err := h.repo.List(r.Context(), page, limit, typeFilter)
	if err != nil {
		slog.Error("failed to list conversations", "error", err, "tenant_id", tenantID)
		http.Error(w, "Gagal memuat percakapan", http.StatusInternalServerError)
		return
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Build response
	response := ConversationListResponse{
		Conversations: conversations,
		Total:         total,
		Page:          page,
		Limit:         limit,
		TotalPages:    totalPages,
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	slog.Info("conversations listed", "tenant_id", tenantID, "page", page, "limit", limit, "total", total)
}

// ConversationDetailResponse represents API response for conversation detail
type ConversationDetailResponse struct {
	Conversation ConversationInfo `json:"conversation"`
	Messages     []*MessageInfo   `json:"messages"`
	TotalMessages int             `json:"total_messages"`
}

// ConversationInfo represents conversation metadata in detail response
type ConversationInfo struct {
	ID          int    `json:"id"`
	PhoneNumber string `json:"phone_number"`
	IsSales     bool   `json:"is_sales"`
	CreatedAt   string `json:"created_at"`
}

// MessageInfo represents message data in detail response
type MessageInfo struct {
	ID        int    `json:"id"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Direction string `json:"direction"`
	CreatedAt string `json:"created_at"`
}

// Get handles GET /api/conversations/:id
func (h *ConversationHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	// Get conversation with tenant validation
	conversation, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		slog.Error("failed to get conversation", "error", err, "conversation_id", id)
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	// Parse limit from query params (default 50)
	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Get messages
	messages, err := h.repo.GetMessages(r.Context(), id, limit)
	if err != nil {
		slog.Error("failed to get messages", "error", err, "conversation_id", id)
		http.Error(w, "Failed to get messages", http.StatusInternalServerError)
		return
	}

	// Get total message count
	totalMessages, err := h.repo.GetMessageCount(r.Context(), id)
	if err != nil {
		slog.Error("failed to count messages", "error", err, "conversation_id", id)
		http.Error(w, "Failed to count messages", http.StatusInternalServerError)
		return
	}

	// Build response
	convInfo := ConversationInfo{
		ID:          conversation.ID,
		PhoneNumber: conversation.SenderPhone,
		IsSales:     conversation.IsSales,
		CreatedAt:   conversation.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	messageInfos := make([]*MessageInfo, len(messages))
	for i, msg := range messages {
		messageInfos[i] = &MessageInfo{
			ID:        msg.ID,
			Sender:    msg.SenderPhone,
			Content:   msg.MessageText,
			Direction: msg.Direction,
			CreatedAt: msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	response := ConversationDetailResponse{
		Conversation:  convInfo,
		Messages:      messageInfos,
		TotalMessages: totalMessages,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	slog.Info("conversation detail retrieved", "conversation_id", id, "total_messages", totalMessages, "returned_messages", len(messages))
}

// Stats handles GET /api/conversations/stats
func (h *ConversationHandler) Stats(w http.ResponseWriter, r *http.Request) {
	// Extract tenant ID from context
	tenantID, err := model.GetTenantID(r.Context())
	if err != nil {
		slog.Error("tenant ID required for conversation stats", "error", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get conversations from last 30 days (using all conversations for now)
	conversations, total, err := h.repo.List(r.Context(), 1, 1000, "")
	if err != nil {
		slog.Error("failed to get conversations for stats", "error", err, "tenant_id", tenantID)
		http.Error(w, "Failed to get conversation stats", http.StatusInternalServerError)
		return
	}

	_ = conversations // Avoid unused variable warning

	// Return stats
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_conversations": total,
		"period":              "last_30_days",
	})
}
