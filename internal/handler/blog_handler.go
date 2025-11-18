package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/riz/auto-lmk/internal/llm"
	"github.com/riz/auto-lmk/internal/model"
	"github.com/riz/auto-lmk/internal/repository"
)

type BlogHandler struct {
	repo        *repository.BlogRepository
	llmProvider llm.Provider
}

func NewBlogHandler(repo *repository.BlogRepository) *BlogHandler {
	return &BlogHandler{repo: repo}
}

func NewBlogHandlerWithLLM(repo *repository.BlogRepository, llmProvider llm.Provider) *BlogHandler {
	return &BlogHandler{
		repo:        repo,
		llmProvider: llmProvider,
	}
}

// List handles GET /api/admin/blog - lists all blog posts for admin
func (h *BlogHandler) List(w http.ResponseWriter, r *http.Request) {
	statusFilter := r.URL.Query().Get("status")

	posts, err := h.repo.List(r.Context(), statusFilter)
	if err != nil {
		slog.Error("failed to list blog posts", "error", err)
		http.Error(w, "Failed to list blog posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"posts": posts,
	})
}

// Get handles GET /api/admin/blog/:id - gets a single blog post
func (h *BlogHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid blog post ID", http.StatusBadRequest)
		return
	}

	post, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		slog.Error("failed to get blog post", "error", err, "id", id)
		http.Error(w, "Blog post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// Create handles POST /api/admin/blog - creates a new blog post
func (h *BlogHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.BlogPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	// Generate slug from title
	slug := h.repo.GenerateSlug(req.Title)

	// Ensure slug is unique
	isUnique, err := h.repo.IsSlugUnique(r.Context(), slug, 0)
	if err != nil {
		slog.Error("failed to check slug uniqueness", "error", err)
		http.Error(w, "Failed to validate slug", http.StatusInternalServerError)
		return
	}

	if !isUnique {
		// Append number to make it unique
		counter := 1
		for {
			testSlug := fmt.Sprintf("%s-%d", slug, counter)
			isUnique, err = h.repo.IsSlugUnique(r.Context(), testSlug, 0)
			if err != nil {
				slog.Error("failed to check slug uniqueness", "error", err)
				http.Error(w, "Failed to validate slug", http.StatusInternalServerError)
				return
			}
			if isUnique {
				slug = testSlug
				break
			}
			counter++
		}
	}

	// Validate status
	if req.Status != "draft" && req.Status != "published" {
		req.Status = "draft"
	}

	post := &model.BlogPost{
		Title:   req.Title,
		Slug:    slug,
		Content: req.Content,
		Excerpt: req.Excerpt,
		Status:  req.Status,
	}

	if err := h.repo.Create(r.Context(), post); err != nil {
		slog.Error("failed to create blog post", "error", err)
		http.Error(w, "Failed to create blog post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// Update handles PUT /api/admin/blog/:id - updates a blog post
func (h *BlogHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid blog post ID", http.StatusBadRequest)
		return
	}

	var req model.BlogPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	// Generate new slug if title changed
	existingPost, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		slog.Error("failed to get existing blog post", "error", err, "id", id)
		http.Error(w, "Blog post not found", http.StatusNotFound)
		return
	}

	slug := existingPost.Slug
	if req.Title != existingPost.Title {
		slug = h.repo.GenerateSlug(req.Title)

		// Ensure slug is unique (excluding current post)
		isUnique, err := h.repo.IsSlugUnique(r.Context(), slug, id)
		if err != nil {
			slog.Error("failed to check slug uniqueness", "error", err)
			http.Error(w, "Failed to validate slug", http.StatusInternalServerError)
			return
		}

		if !isUnique {
			// Append number to make it unique
			counter := 1
			for {
				testSlug := fmt.Sprintf("%s-%d", slug, counter)
				isUnique, err = h.repo.IsSlugUnique(r.Context(), testSlug, id)
				if err != nil {
					slog.Error("failed to check slug uniqueness", "error", err)
					http.Error(w, "Failed to validate slug", http.StatusInternalServerError)
					return
				}
				if isUnique {
					slug = testSlug
					break
				}
				counter++
			}
		}
	}

	// Validate status
	if req.Status != "draft" && req.Status != "published" {
		req.Status = "draft"
	}

	updates := map[string]interface{}{
		"title":   req.Title,
		"slug":    slug,
		"content": req.Content,
		"excerpt": req.Excerpt,
		"status":  req.Status,
	}

	if err := h.repo.Update(r.Context(), id, updates); err != nil {
		slog.Error("failed to update blog post", "error", err, "id", id)
		http.Error(w, "Failed to update blog post", http.StatusInternalServerError)
		return
	}

	// Get updated post
	post, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		slog.Error("failed to get updated blog post", "error", err, "id", id)
		http.Error(w, "Failed to get updated blog post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// Delete handles DELETE /api/admin/blog/:id - deletes a blog post
func (h *BlogHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid blog post ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		slog.Error("failed to delete blog post", "error", err, "id", id)
		http.Error(w, "Failed to delete blog post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Blog post deleted successfully"}`))
}

// GenerateAI handles POST /api/admin/blog/generate-ai - generates blog content using AI
func (h *BlogHandler) GenerateAI(w http.ResponseWriter, r *http.Request) {
	if h.llmProvider == nil {
		http.Error(w, "AI service not available", http.StatusServiceUnavailable)
		return
	}

	var req struct {
		Topic  string `json:"topic"`
		Type   string `json:"type"`             // "full_post" or "excerpt"
		Length string `json:"length,omitempty"` // "short", "medium", "long"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Topic == "" {
		http.Error(w, "Topic is required", http.StatusBadRequest)
		return
	}

	if req.Type == "" {
		req.Type = "full_post"
	}

	if req.Length == "" {
		req.Length = "medium"
	}

	// Prepare AI prompt based on type
	var systemPrompt, userPrompt string

	switch req.Type {
	case "excerpt":
		systemPrompt = `You are an expert automotive content writer. Create a compelling excerpt for a blog post about cars.
Write in Indonesian language. Keep it under 200 words. Make it engaging and SEO-friendly.
Focus on the key benefits or interesting facts about the topic.`
		userPrompt = fmt.Sprintf("Tulis excerpt untuk artikel blog tentang: %s", req.Topic)

	case "full_post":
		var lengthInstruction string
		switch req.Length {
		case "short":
			lengthInstruction = "500-800 kata"
		case "medium":
			lengthInstruction = "800-1200 kata"
		case "long":
			lengthInstruction = "1200-1500 kata"
		default:
			lengthInstruction = "800-1200 kata"
		}

		systemPrompt = fmt.Sprintf(`You are an expert automotive content writer. Create a complete blog post about cars.
Write in Indonesian language. Length: %s. Include:
- Engaging introduction
- Main content with detailed information
- Practical tips or insights
- SEO-friendly keywords
- Conclusion
Use proper headings (H2, H3) and structure the content well.`, lengthInstruction)
		userPrompt = fmt.Sprintf("Tulis artikel blog lengkap tentang: %s", req.Topic)

	default:
		http.Error(w, "Invalid type. Use 'excerpt' or 'full_post'", http.StatusBadRequest)
		return
	}

	messages := []llm.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	response, err := h.llmProvider.Chat(r.Context(), messages, nil)
	if err != nil {
		slog.Error("AI generation failed", "error", err)
		http.Error(w, "AI processing failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result := map[string]interface{}{
		"content": strings.TrimSpace(response.Content),
		"type":    req.Type,
		"topic":   req.Topic,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
