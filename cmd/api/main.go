package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Auto LMK - Multi-Tenant Car Sales Platform")
	fmt.Println("===========================================")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server will run on port %s", port)
	log.Println("Setup complete! Ready for development...")

	// TODO: Initialize server, database, routes, etc.
}
