package model

import "time"

// SearchAnalytics represents search keyword analytics
type SearchAnalytics struct {
	Keyword        string    `json:"keyword"`
	SearchCount    int       `json:"search_count"`
	LastSearchedAt time.Time `json:"last_searched_at"`
}

// CarViewAnalytics represents car view analytics
type CarViewAnalytics struct {
	CarID        int       `json:"car_id"`
	Brand        string    `json:"brand"`
	Model        string    `json:"model"`
	Year         int       `json:"year"`
	ViewCount    int       `json:"view_count"`
	LastViewedAt time.Time `json:"last_viewed_at"`
}

// SearchTrend represents search trends over time
type SearchTrend struct {
	Date        time.Time `json:"date"`
	SearchCount int       `json:"search_count"`
}

// WhatsAppSettings represents WhatsApp configuration for a tenant
type WhatsAppSettings struct {
	TenantID       int       `json:"tenant_id"`
	BotNumber      *string   `json:"bot_number,omitempty"`
	FallbackNumber string    `json:"fallback_number"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
