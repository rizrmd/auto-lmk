package model

import "time"

// BrandingSettings represents tenant-specific branding configuration
type BrandingSettings struct {
	TenantID       int       `json:"tenant_id" db:"tenant_id"`
	LogoPath       *string   `json:"logo_path,omitempty" db:"logo_path"`
	FaviconPath    *string   `json:"favicon_path,omitempty" db:"favicon_path"`
	CustomTitle    *string   `json:"custom_title,omitempty" db:"custom_title"`
	CustomSubtitle *string   `json:"custom_subtitle,omitempty" db:"custom_subtitle"`
	PromoText      *string   `json:"promo_text,omitempty" db:"promo_text"`
	HeaderStyle    string    `json:"header_style" db:"header_style"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// BrandingUpdateRequest represents the request body for updating branding settings
type BrandingUpdateRequest struct {
	CustomTitle    *string `json:"custom_title"`
	CustomSubtitle *string `json:"custom_subtitle"`
	PromoText      *string `json:"promo_text"`
	HeaderStyle    string  `json:"header_style"`
}
