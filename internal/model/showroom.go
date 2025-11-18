package model

import "time"

// ShowroomSettings represents the showroom information for a tenant
type ShowroomSettings struct {
	TenantID      int        `json:"tenant_id" db:"tenant_id"`
	Address       *string    `json:"address,omitempty" db:"address"`
	Phone         *string    `json:"phone,omitempty" db:"phone"`
	Email         *string    `json:"email,omitempty" db:"email"`
	BusinessHours *string    `json:"business_hours,omitempty" db:"business_hours"`
	Latitude      *float64   `json:"latitude,omitempty" db:"latitude"`
	Longitude     *float64   `json:"longitude,omitempty" db:"longitude"`
	MapEmbed      *string    `json:"map_embed,omitempty" db:"map_embed"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// ShowroomUpdateRequest represents the request to update showroom settings
type ShowroomUpdateRequest struct {
	Address       *string  `json:"address,omitempty"`
	Phone         *string  `json:"phone,omitempty"`
	Email         *string  `json:"email,omitempty"`
	BusinessHours *string  `json:"business_hours,omitempty"`
	Latitude      *float64 `json:"latitude,omitempty"`
	Longitude     *float64 `json:"longitude,omitempty"`
	MapEmbed      *string  `json:"map_embed,omitempty"`
}
