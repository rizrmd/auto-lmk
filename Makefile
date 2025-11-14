.PHONY: help dev build run migrate-up migrate-down migrate-create docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev: ## Run with hot reload (air)
	air

build: ## Build the application
	go build -o bin/api cmd/api/main.go

run: ## Run the application
	go run cmd/api/main.go

migrate-up: ## Run database migrations
	migrate -path migrations -database "postgresql://autolmk:autolmk_dev_password@localhost:5432/autolmk?sslmode=disable" up

migrate-down: ## Rollback last migration
	migrate -path migrations -database "postgresql://autolmk:autolmk_dev_password@localhost:5432/autolmk?sslmode=disable" down 1

migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	migrate create -ext sql -dir migrations -seq $(name)

docker-up: ## Start Docker containers
	docker-compose up -d

docker-down: ## Stop Docker containers
	docker-compose down

test: ## Run tests
	go test -v ./...

tidy: ## Tidy go modules
	go mod tidy

clean: ## Clean build artifacts
	rm -rf bin/ tmp/
