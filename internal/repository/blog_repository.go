package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/riz/auto-lmk/internal/model"
)

type BlogRepository struct {
	db *sql.DB
}

func NewBlogRepository(db *sql.DB) *BlogRepository {
	return &BlogRepository{db: db}
}

// Create creates a new blog post
func (r *BlogRepository) Create(ctx context.Context, post *model.BlogPost) error {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO blog_posts (tenant_id, title, slug, content, excerpt, status, published_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	var publishedAt *time.Time
	if post.Status == "published" && post.PublishedAt == nil {
		now := time.Now()
		publishedAt = &now
	} else {
		publishedAt = post.PublishedAt
	}

	err = r.db.QueryRowContext(ctx, query,
		tenantID,
		post.Title,
		post.Slug,
		post.Content,
		post.Excerpt,
		post.Status,
		publishedAt,
		post.CreatedBy,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err
	}

	post.TenantID = tenantID
	return nil
}

// GetByID gets a blog post by ID
func (r *BlogRepository) GetByID(ctx context.Context, id int) (*model.BlogPost, error) {
	query := `
		SELECT id, tenant_id, title, slug, content, excerpt, status, published_at, created_by, created_at, updated_at
		FROM blog_posts
		WHERE id = $1
	`

	post := &model.BlogPost{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.TenantID,
		&post.Title,
		&post.Slug,
		&post.Content,
		&post.Excerpt,
		&post.Status,
		&post.PublishedAt,
		&post.CreatedBy,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return post, nil
}

// GetBySlug gets a published blog post by slug
func (r *BlogRepository) GetBySlug(ctx context.Context, slug string) (*model.BlogPost, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, tenant_id, title, slug, content, excerpt, status, published_at, created_by, created_at, updated_at
		FROM blog_posts
		WHERE tenant_id = $1 AND slug = $2 AND status = 'published'
	`

	post := &model.BlogPost{}
	err = r.db.QueryRowContext(ctx, query, tenantID, slug).Scan(
		&post.ID,
		&post.TenantID,
		&post.Title,
		&post.Slug,
		&post.Content,
		&post.Excerpt,
		&post.Status,
		&post.PublishedAt,
		&post.CreatedBy,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return post, nil
}

// List lists blog posts with optional status filter
func (r *BlogRepository) List(ctx context.Context, statusFilter string) ([]*model.BlogPost, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, tenant_id, title, slug, content, excerpt, status, published_at, created_by, created_at, updated_at
		FROM blog_posts
		WHERE tenant_id = $1
	`

	args := []interface{}{tenantID}
	if statusFilter != "" {
		query += " AND status = $2"
		args = append(args, statusFilter)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.BlogPost
	for rows.Next() {
		post := &model.BlogPost{}
		err := rows.Scan(
			&post.ID,
			&post.TenantID,
			&post.Title,
			&post.Slug,
			&post.Content,
			&post.Excerpt,
			&post.Status,
			&post.PublishedAt,
			&post.CreatedBy,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// ListPublished lists published blog posts for frontend
func (r *BlogRepository) ListPublished(ctx context.Context, limit int) ([]*model.BlogPost, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, tenant_id, title, slug, content, excerpt, status, published_at, created_by, created_at, updated_at
		FROM blog_posts
		WHERE tenant_id = $1 AND status = 'published'
		ORDER BY published_at DESC
	`

	if limit > 0 {
		query += " LIMIT $2"
	}

	var args []interface{}
	args = append(args, tenantID)
	if limit > 0 {
		args = append(args, limit)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.BlogPost
	for rows.Next() {
		post := &model.BlogPost{}
		err := rows.Scan(
			&post.ID,
			&post.TenantID,
			&post.Title,
			&post.Slug,
			&post.Content,
			&post.Excerpt,
			&post.CreatedBy,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// Update updates a blog post
func (r *BlogRepository) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	query := "UPDATE blog_posts SET "
	args := []interface{}{}
	i := 1

	for key, value := range updates {
		if i > 1 {
			query += ", "
		}
		query += key + " = $" + string(rune('0'+i))
		args = append(args, value)
		i++
	}

	// Handle special case for published_at when status changes to published
	if status, ok := updates["status"]; ok && status == "published" {
		if i > 1 {
			query += ", "
		}
		query += "published_at = CASE WHEN published_at IS NULL THEN CURRENT_TIMESTAMP ELSE published_at END"
	}

	query += ", updated_at = CURRENT_TIMESTAMP WHERE id = $" + string(rune('0'+i))
	args = append(args, id)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// Delete deletes a blog post
func (r *BlogRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM blog_posts WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GenerateSlug generates a URL-friendly slug from title
func (r *BlogRepository) GenerateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces and special characters with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove non-alphanumeric characters except hyphens
	result := ""
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result += string(char)
		}
	}

	// Remove multiple consecutive hyphens
	for strings.Contains(result, "--") {
		result = strings.ReplaceAll(result, "--", "-")
	}

	// Trim hyphens from start and end
	result = strings.Trim(result, "-")

	return result
}

// IsSlugUnique checks if a slug is unique for the tenant
func (r *BlogRepository) IsSlugUnique(ctx context.Context, slug string, excludeID int) (bool, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return false, err
	}

	query := "SELECT COUNT(*) FROM blog_posts WHERE tenant_id = $1 AND slug = $2"
	args := []interface{}{tenantID, slug}

	if excludeID > 0 {
		query += " AND id != $3"
		args = append(args, excludeID)
	}

	var count int
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}
