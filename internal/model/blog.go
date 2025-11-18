package model

import "time"

// BlogPost represents a blog post
type BlogPost struct {
	ID          int        `json:"id"`
	TenantID    int        `json:"tenant_id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Content     string     `json:"content"`
	Excerpt     string     `json:"excerpt,omitempty"`
	Status      string     `json:"status"` // 'draft' or 'published'
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedBy   *int       `json:"created_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// BlogPostRequest represents a request to create/update a blog post
type BlogPostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Excerpt string `json:"excerpt,omitempty"`
	Status  string `json:"status"`
}
