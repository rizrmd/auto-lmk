package model

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type Sales struct {
	ID           int       `json:"id"`
	TenantID     int       `json:"tenant_id"`
	PhoneNumber  string    `json:"phone_number"`
	Name         string    `json:"name"`
	Role         string    `json:"role,omitempty"`
	Status       string    `json:"status"`
	RegisteredAt time.Time `json:"registered_at"`
}

type CreateSalesRequest struct {
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Role        string `json:"role,omitempty"`
}

func (r *CreateSalesRequest) Validate() error {
	// Phone validation - E.164 (international) or Indonesia local format
	// E.164: +[country code][subscriber number] - min 8 digits, max 15 digits total
	// Indonesia local: 08[subscriber] - 10-13 digits total
	phoneRegex := regexp.MustCompile(`^\+[1-9]\d{7,14}$|^08\d{8,11}$`)
	if !phoneRegex.MatchString(r.PhoneNumber) {
		return errors.New("Format nomor telepon tidak valid")
	}

	// Name validation
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("Nama tidak boleh kosong")
	}

	return nil
}
