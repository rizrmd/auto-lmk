package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/riz/auto-lmk/internal/model"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

// LogSearchEvent logs a search event for analytics
func (r *AnalyticsRepository) LogSearchEvent(ctx context.Context, keyword string, resultCount int) error {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO search_analytics (tenant_id, keyword, search_count, last_searched_at)
		VALUES ($1, $2, 1, CURRENT_TIMESTAMP)
		ON CONFLICT (tenant_id, keyword)
		DO UPDATE SET
			search_count = search_analytics.search_count + 1,
			last_searched_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err = r.db.ExecContext(ctx, query, tenantID, keyword)
	return err
}

// LogCarView logs a car view event for analytics
func (r *AnalyticsRepository) LogCarView(ctx context.Context, carID int) error {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO car_views (tenant_id, car_id, view_count, last_viewed_at)
		VALUES ($1, $2, 1, CURRENT_TIMESTAMP)
		ON CONFLICT (tenant_id, car_id)
		DO UPDATE SET
			view_count = car_views.view_count + 1,
			last_viewed_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err = r.db.ExecContext(ctx, query, tenantID, carID)
	return err
}

// GetTopSearchKeywords returns top search keywords for a tenant
func (r *AnalyticsRepository) GetTopSearchKeywords(ctx context.Context, limit int) ([]*model.SearchAnalytics, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT keyword, search_count, last_searched_at
		FROM search_analytics
		WHERE tenant_id = $1
		ORDER BY search_count DESC, last_searched_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.SearchAnalytics
	for rows.Next() {
		analytics := &model.SearchAnalytics{}
		err := rows.Scan(&analytics.Keyword, &analytics.SearchCount, &analytics.LastSearchedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, analytics)
	}

	return results, nil
}

// GetTopViewedCars returns top viewed cars for a tenant
func (r *AnalyticsRepository) GetTopViewedCars(ctx context.Context, limit int) ([]*model.CarViewAnalytics, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT c.id, c.brand, c.model, c.year, cv.view_count, cv.last_viewed_at
		FROM car_views cv
		JOIN cars c ON cv.car_id = c.id
		WHERE cv.tenant_id = $1 AND c.tenant_id = $1
		ORDER BY cv.view_count DESC, cv.last_viewed_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.CarViewAnalytics
	for rows.Next() {
		analytics := &model.CarViewAnalytics{}
		err := rows.Scan(&analytics.CarID, &analytics.Brand, &analytics.Model, &analytics.Year, &analytics.ViewCount, &analytics.LastViewedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, analytics)
	}

	return results, nil
}

// GetSearchTrends returns search trends for date range
func (r *AnalyticsRepository) GetSearchTrends(ctx context.Context, startDate, endDate time.Time) ([]*model.SearchTrend, error) {
	tenantID, err := model.GetTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT DATE(last_searched_at) as date, COUNT(*) as search_count
		FROM search_analytics
		WHERE tenant_id = $1 AND last_searched_at BETWEEN $2 AND $3
		GROUP BY DATE(last_searched_at)
		ORDER BY date
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.SearchTrend
	for rows.Next() {
		trend := &model.SearchTrend{}
		err := rows.Scan(&trend.Date, &trend.SearchCount)
		if err != nil {
			return nil, err
		}
		results = append(results, trend)
	}

	return results, nil
}
