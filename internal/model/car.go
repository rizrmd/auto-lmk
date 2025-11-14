package model

import "time"

type Car struct {
	ID           int       `json:"id"`
	TenantID     int       `json:"tenant_id"`
	Brand        string    `json:"brand"`
	Model        string    `json:"model"`
	Year         int       `json:"year"`
	Price        int64     `json:"price"`
	Mileage      *int      `json:"mileage,omitempty"`
	Transmission *string   `json:"transmission,omitempty"`
	FuelType     *string   `json:"fuel_type,omitempty"`
	EngineCC     *int      `json:"engine_cc,omitempty"`
	Seats        *int      `json:"seats,omitempty"`
	Color        *string   `json:"color,omitempty"`
	Description  *string   `json:"description,omitempty"`
	Status       string    `json:"status"`
	IsFeatured   bool      `json:"is_featured"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CarPhoto struct {
	ID           int       `json:"id"`
	CarID        int       `json:"car_id"`
	FilePath     string    `json:"file_path"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
}

type CarSpec struct {
	ID        int       `json:"id"`
	CarID     int       `json:"car_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}
