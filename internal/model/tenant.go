package model

import "time"

type Tenant struct {
	ID             int       `json:"id"`
	Domain         string    `json:"domain"`
	Name           string    `json:"name"`
	WhatsAppNumber *string   `json:"whatsapp_number,omitempty"`
	PairingStatus  string    `json:"pairing_status"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateTenantRequest struct {
	Domain         string  `json:"domain"`
	Name           string  `json:"name"`
	WhatsAppNumber *string `json:"whatsapp_number,omitempty"`
}
