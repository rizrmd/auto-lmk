package handler

import (
	"testing"

	"github.com/riz/auto-lmk/internal/model"
)

func TestCreateSalesRequest_Validate_ValidPhone(t *testing.T) {
	req := model.CreateSalesRequest{
		PhoneNumber: "08123456789",
		Name:        "John Doe",
	}

	err := req.Validate()
	if err != nil {
		t.Errorf("Expected no error for valid phone, got %v", err)
	}
}

func TestCreateSalesRequest_Validate_ValidE164Phone(t *testing.T) {
	req := model.CreateSalesRequest{
		PhoneNumber: "+628123456789",
		Name:        "John Doe",
	}

	err := req.Validate()
	if err != nil {
		t.Errorf("Expected no error for valid E164 phone, got %v", err)
	}
}

func TestCreateSalesRequest_Validate_InvalidPhone(t *testing.T) {
	req := model.CreateSalesRequest{
		PhoneNumber: "invalid",
		Name:        "John Doe",
	}

	err := req.Validate()
	if err == nil {
		t.Error("Expected error for invalid phone")
	}
	if err.Error() != "Format nomor telepon tidak valid" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestCreateSalesRequest_Validate_EmptyName(t *testing.T) {
	req := model.CreateSalesRequest{
		PhoneNumber: "08123456789",
		Name:        "",
	}

	err := req.Validate()
	if err == nil {
		t.Error("Expected error for empty name")
	}
	if err.Error() != "Nama tidak boleh kosong" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}
