package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/riz/auto-lmk/internal/repository"
)

type ConversationHandler struct {
	repo *repository.ConversationRepository
}

func NewConversationHandler(repo *repository.ConversationRepository) *ConversationHandler {
	return &ConversationHandler{repo: repo}
}

// List handles GET /api/conversations
func (h *ConversationHandler) List(w http.ResponseWriter, r *http.Request) {
	conversations, err := h.repo.ListConversations(r.Context())
	if err != nil {
		slog.Error("failed to list conversations", "error", err)
		http.Error(w, "Failed to list conversations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  conversations,
		"count": len(conversations),
	})
}

// Get handles GET /api/conversations/:id
func (h *ConversationHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	limit := 50 // Default limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	messages, err := h.repo.GetMessages(r.Context(), id, limit)
	if err != nil {
		slog.Error("failed to get messages", "error", err, "conversation_id", id)
		http.Error(w, "Failed to get messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"conversation_id": id,
		"messages":        messages,
		"count":           len(messages),
	})
}
