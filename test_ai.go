package main

import (
	"context"
	"fmt"
	"log"

	"github.com/riz/auto-lmk/internal/handler"
	"github.com/riz/auto-lmk/internal/repository"
	"github.com/riz/auto-lmk/pkg/config"
	"github.com/riz/auto-lmk/pkg/database"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create repository
	carRepo := repository.NewCarRepository(db.DB)

	// Create handler
	handler := handler.NewCarHandler(carRepo)

	// Create context with tenant_id
	ctx := context.Background()
	// Note: In a real request, tenant_id would be extracted from middleware
	// For testing, we'll simulate tenant_id = 1

	fmt.Println("Testing enhanced AI generation...")

	// Test cases
	testCases := []string{
		"hyundai creta 2022 hitam 15000 km",
		"Tesla Model Y 2023 merah",
		"mercedes benz c300",
		"toyota avanza",
		"Honda CR-V",
	}

	for i, testCase := range testCases {
		fmt.Printf("\n=== Test %d: %s ===\n", i+1, testCase)

		// This will fail without proper tenant context, but let's see what happens
		result := handler.GenerateCarDataFromDatabase(ctx, testCase)

		fmt.Printf("Brand: %v\n", result["brand"])
		fmt.Printf("Model: %v\n", result["model"])
		fmt.Printf("Year: %v\n", result["year"])
		fmt.Printf("Price: %v\n", result["price"])
	}
}