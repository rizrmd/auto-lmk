package model

import "time"

type Sales struct {
	ID           int       `json:"id"`
	TenantID     int       `json:"tenant_id"`
	PhoneNumber  string    `json:"phone_number"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	RegisteredAt time.Time `json:"registered_at"`
}

type CreateSalesRequest struct {
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
}
