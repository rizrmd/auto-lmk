# Development Guide - Auto LMK

> **Generated:** 2025-11-15
> **Project:** Auto LMK - Car Sales Platform with AI-Powered WhatsApp Bot
> **Version:** 1.0 (Production-Ready)

---

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Environment Setup](#environment-setup)
- [Development Workflow](#development-workflow)
- [Database Management](#database-management)
- [Build & Run](#build--run)
- [Testing](#testing)
- [Common Development Tasks](#common-development-tasks)
- [Troubleshooting](#troubleshooting)

---

## ✅ Prerequisites

Before starting development, ensure you have the following installed:

| Requirement | Version | Installation Guide |
|-------------|---------|-------------------|
| **Go** | 1.21+ | [go.dev/doc/install](https://go.dev/doc/install) |
| **PostgreSQL** | 15+ | [postgresql.org/download](https://www.postgresql.org/download/) or use Docker |
| **Node.js & NPM** | 18+ (for Tailwind CSS build) | [nodejs.org](https://nodejs.org/) |
| **golang-migrate** | latest | `brew install golang-migrate` or [github](https://github.com/golang-migrate/migrate) |
| **Air** (optional) | latest | `go install github.com/cosmtrek/air@latest` |
| **Docker** (optional) | latest | [docker.com](https://www.docker.com/) |
| **Git** | latest | [git-scm.com](https://git-scm.com/) |

### Optional Requirements
- **WhatsApp Account:** For WhatsApp bot integration (can use personal number)
- **LLM API Key:** Z.AI, OpenAI, or Anthropic API key for chatbot functionality

---

## 🚀 Quick Start

Get the project running in 5 minutes:

```bash
# 1. Clone repository
git clone https://github.com/riz/auto-lmk.git
cd auto-lmk

# 2. Install Go dependencies
go mod download

# 3. Install Node.js dependencies (for Tailwind CSS)
npm install

# 4. Copy environment template
cp .env.example .env

# 5. Configure environment variables (see below)
vim .env  # or use your preferred editor

# 6. Start PostgreSQL (using Docker)
docker-compose up -d postgres

# 7. Wait for PostgreSQL to be ready
sleep 5

# 8. Run database migrations
make migrate-up

# 9. (Optional) Seed sample data
psql -U autolmk -d autolmk -h localhost < scripts/seed.sql

# 10. Build Tailwind CSS
npm run css:build

# 11. Start development server with hot reload
make dev

# ✅ Server running at http://localhost:8080
```

---

## ⚙️ Environment Setup

### 1. Configure `.env` File

Edit `.env` with your specific values:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=autolmk
DB_PASSWORD=your_secure_password_here  # ⚠️ CHANGE THIS
DB_NAME=autolmk
DB_SSLMODE=disable  # Use 'require' in production

# Server Configuration
PORT=8080
ENV=development  # Options: development, production

# LLM Configuration
LLM_PROVIDER=zai  # Options: zai, openai, anthropic
LLM_API_KEY=f0fd028dc2d4462e82216f0b86d1bfe3.MfQho2gEDOtxm4dy  # ✅ Z.AI API key
LLM_MODEL=glm-4.6  # Z.AI model
ZAI_ENDPOINT=https://api.z.ai/api/coding/paas/v4/chat/completions

# WhatsApp Configuration
WHATSAPP_SESSION_PATH=./whatsapp_sessions

# Security
JWT_SECRET=change-this-to-a-random-secure-string-in-production  # ⚠️ CHANGE THIS
```

### 2. LLM Provider Configuration

**Current Provider: Z.AI (glm-4.6)**

**Configuration:**
```bash
LLM_PROVIDER=zai
LLM_MODEL=glm-4.6
LLM_API_KEY=f0fd028dc2d4462e82216f0b86d1bfe3.MfQho2gEDOtxm4dy
ZAI_ENDPOINT=https://api.z.ai/api/coding/paas/v4/chat/completions
```

**Alternative Providers (Not yet implemented):**

- **OpenAI (GPT-4):**
  ```bash
  LLM_PROVIDER=openai
  LLM_MODEL=gpt-4o-mini
  LLM_API_KEY=sk-your-openai-key-here
  ```

- **Anthropic (Claude):**
  ```bash
  LLM_PROVIDER=anthropic
  LLM_MODEL=claude-3-5-haiku-20241022
  LLM_API_KEY=sk-ant-your-anthropic-key-here
  ```

### 3. Database Setup

**Option A: Docker (Recommended for Development)**
```bash
# Start PostgreSQL container
docker-compose up -d postgres

# Check logs
docker-compose logs -f postgres

# Access PostgreSQL CLI
docker exec -it auto-lmk-postgres psql -U autolmk -d autolmk
```

**Option B: Local PostgreSQL Installation**
```bash
# Create database and user
psql -U postgres
CREATE DATABASE autolmk;
CREATE USER autolmk WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE autolmk TO autolmk;
\q

# Update .env with your credentials
DB_HOST=localhost
DB_USER=autolmk
DB_PASSWORD=your_password
```

---

## 💻 Development Workflow

### Hot Reload Development with Air

**Recommended for active development:**

```bash
# Start development server with auto-reload
make dev

# Or directly:
air
```

**What Air does:**
- Watches for file changes in `cmd/`, `internal/`, `pkg/`
- Automatically rebuilds and restarts the server
- Displays logs in terminal
- Configuration in `.air.toml`

**Air Configuration (`.air.toml`):**
```toml
[build]
  cmd = "go build -o ./tmp/main ./cmd/api"
  bin = "tmp/main"
  exclude_dir = ["node_modules", "whatsapp_sessions", "bin", "tmp"]
  include_ext = ["go"]
  delay = 1000  # ms
```

### Manual Development (Without Hot Reload)

```bash
# Run directly with go run
go run cmd/api/main.go

# Or use Makefile
make run
```

### Frontend Development (Tailwind CSS)

**Build CSS:**
```bash
# Build once
npm run css:build

# Watch for changes (auto-rebuild)
npm run css:watch
```

**Tailwind Configuration:**
- Config file: `tailwind.config.js`
- Build tool: Vite (`vite.config.js`)
- Output: `static/css/output.css`

---

## 🗄️ Database Management

### Migrations

All migrations are in the `migrations/` directory.

**Create a New Migration:**
```bash
# Syntax: make migrate-create name=migration_name
make migrate-create name=add_user_roles

# This creates:
# migrations/000009_add_user_roles.up.sql
# migrations/000009_add_user_roles.down.sql
```

**Run Migrations (Apply Schema Changes):**
```bash
# Apply all pending migrations
make migrate-up

# Example output:
# 000001_create_tenants_table.up.sql  ✓
# 000002_create_users_table.up.sql    ✓
# ...
```

**Rollback Migrations:**
```bash
# Rollback last migration only
make migrate-down

# Rollback specific number of migrations
migrate -path migrations -database "postgresql://..." down 3
```

**Check Migration Status:**
```bash
# View current migration version
migrate -path migrations -database "postgresql://autolmk:autolmk_dev_password@localhost:5432/autolmk?sslmode=disable" version
```

### Seed Data

**Load Sample Data (for testing):**
```bash
psql -U autolmk -d autolmk -h localhost < scripts/seed.sql
```

**Sample data includes:**
- 1 tenant (showroom-jaya.localhost)
- 1 admin user
- 2 sales team members
- 10 sample cars
- Car photos and specs

### Manual Database Operations

**Access PostgreSQL CLI:**
```bash
# Via Docker
docker exec -it auto-lmk-postgres psql -U autolmk -d autolmk

# Via local installation
psql -U autolmk -d autolmk -h localhost
```

**Useful SQL Queries:**
```sql
-- List all tenants
SELECT id, domain, name FROM tenants;

-- List all cars for a tenant
SELECT brand, model, year, price FROM cars WHERE tenant_id = 1;

-- Check WhatsApp conversations
SELECT id, sender_phone, created_at FROM conversations WHERE tenant_id = 1;

-- View database size
SELECT pg_size_pretty(pg_database_size('autolmk'));
```

---

## 🔨 Build & Run

### Development Build

```bash
# Run with go run (no binary)
make run
# or
go run cmd/api/main.go
```

### Production Build

```bash
# Build optimized binary
make build

# Output: bin/api (Linux binary, ~11MB)

# Run the binary
./bin/api

# Or with environment variables
PORT=8080 ENV=production ./bin/api
```

### Cross-Platform Build

```bash
# Build for Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api-linux cmd/api/main.go

# Build for macOS
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/api-macos cmd/api/main.go

# Build for Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/api.exe cmd/api/main.go
```

### Docker Build

```bash
# Build and run with Docker Compose
docker-compose up --build

# Build only the app container
docker build -t auto-lmk:latest .

# Run the container
docker run -p 8080:8080 --env-file .env auto-lmk:latest
```

---

## 🧪 Testing

### Run All Tests

```bash
# Run all tests
make test

# Or directly with go test
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Test Individual Packages

```bash
# Test specific package
go test -v ./internal/service

# Test with verbose output
go test -v ./internal/repository -count=1

# Run specific test function
go test -v -run TestCarRepository_FindByID ./internal/repository
```

### Manual API Testing

**Using `curl`:**

```bash
# Health check
curl http://localhost:8080/health

# Create tenant (root admin)
curl -X POST http://localhost:8080/api/root/tenants \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "showroom-jaya.localhost",
    "name": "Showroom Jaya",
    "whatsapp_number": "6281234567890"
  }'

# List cars (tenant-scoped)
curl http://showroom-jaya.localhost:8080/api/cars

# Create car
curl -X POST http://showroom-jaya.localhost:8080/api/cars \
  -H "Content-Type: application/json" \
  -d '{
    "brand": "Toyota",
    "model": "Avanza",
    "year": 2023,
    "price": 250000000,
    "transmission": "Automatic",
    "fuel_type": "Gasoline",
    "mileage": 5000
  }'
```

**Testing with Local Domain:**

Add to `/etc/hosts`:
```
127.0.0.1  showroom-jaya.localhost
127.0.0.1  another-tenant.localhost
```

Then access:
- http://showroom-jaya.localhost:8080
- http://another-tenant.localhost:8080

---

## 🛠️ Common Development Tasks

### 1. Add a New API Endpoint

**Steps:**
1. Add route in `cmd/api/main.go`
2. Create handler function in `internal/handler/`
3. Create service function in `internal/service/`
4. Create repository function in `internal/repository/`

**Example (Add "Get Car by ID" endpoint):**

```go
// 1. cmd/api/main.go
r.Get("/api/cars/{id}", carHandler.GetByID)

// 2. internal/handler/car_handler.go
func (h *CarHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    carID, _ := strconv.Atoi(chi.URLParam(r, "id"))
    tenantID := model.GetTenantID(r.Context())

    car, err := h.service.GetByID(r.Context(), tenantID, carID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(car)
}

// 3. internal/service/car_service.go
func (s *CarService) GetByID(ctx context.Context, tenantID, carID int) (*model.Car, error) {
    return s.repo.FindByID(ctx, tenantID, carID)
}

// 4. internal/repository/car_repository.go
func (r *CarRepository) FindByID(ctx context.Context, tenantID, carID int) (*model.Car, error) {
    var car model.Car
    err := r.db.QueryRowContext(ctx,
        "SELECT * FROM cars WHERE tenant_id = $1 AND id = $2",
        tenantID, carID,
    ).Scan(&car.ID, &car.Brand, &car.Model, ...)
    return &car, err
}
```

### 2. Add a New Database Table

```bash
# Create migration
make migrate-create name=create_bookings_table

# Edit migrations/000009_create_bookings_table.up.sql
CREATE TABLE IF NOT EXISTS bookings (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER NOT NULL REFERENCES tenants(id),
    customer_name VARCHAR(255) NOT NULL,
    car_id INTEGER NOT NULL REFERENCES cars(id),
    booking_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bookings_tenant_id ON bookings(tenant_id);

# Edit migrations/000009_create_bookings_table.down.sql
DROP TABLE IF EXISTS bookings;

# Run migration
make migrate-up
```

### 3. Add a New Frontend Page

**Steps:**
1. Create template in `templates/pages/`
2. Add route in `cmd/api/main.go`
3. Add handler in `internal/handler/page_handler.go`

**Example (Add "About" page):**

```html
<!-- templates/pages/about.html -->
{{template "layouts/base.html" .}}
{{define "content"}}
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold">About Auto LMK</h1>
    <p class="mt-4">We revolutionize car sales...</p>
  </div>
{{end}}
```

```go
// cmd/api/main.go
r.Get("/about", pageHandler.About)

// internal/handler/page_handler.go
func (h *PageHandler) About(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles(
        "templates/layouts/base.html",
        "templates/pages/about.html",
    ))
    tmpl.Execute(w, nil)
}
```

### 4. Add New LLM Function for WhatsApp Bot

**Edit `internal/llm/bot.go`:**

```go
// Add function definition
{
    Name: "scheduletestDrive",
    Description: "Schedule a test drive appointment for a car",
    Parameters: map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "customerName": map[string]string{"type": "string"},
            "carID":        map[string]string{"type": "integer"},
            "date":         map[string]string{"type": "string"},
        },
        "required": []string{"customerName", "carID", "date"},
    },
}

// Add function execution
case "scheduleTestDrive":
    // Parse parameters
    // Call service layer
    // Return formatted response
```

### 5. Configure Multi-Tenant Domain

**Development (localhost):**

Add to `/etc/hosts`:
```
127.0.0.1  tenant1.localhost
127.0.0.1  tenant2.localhost
```

**Production (with real domain):**

1. Add DNS A record: `showroom-jaya.yourdomain.com → YOUR_SERVER_IP`
2. Create tenant in database:
   ```sql
   INSERT INTO tenants (domain, name, whatsapp_number)
   VALUES ('showroom-jaya.yourdomain.com', 'Showroom Jaya', '628123456789');
   ```
3. Configure Nginx with wildcard SSL (Let's Encrypt)

---

## 🐛 Troubleshooting

### Common Issues

**1. Database Connection Failed**
```
Error: pq: password authentication failed for user "autolmk"
```
**Solution:**
- Check `.env` file has correct DB credentials
- Verify PostgreSQL is running: `docker-compose ps`
- Check PostgreSQL logs: `docker-compose logs postgres`

**2. Port Already in Use**
```
Error: bind: address already in use
```
**Solution:**
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or change port in .env
PORT=8081
```

**3. Migration Failed**
```
Error: Dirty database version
```
**Solution:**
```bash
# Force migration version
migrate -path migrations -database "postgresql://..." force 8

# Then re-run migrations
make migrate-up
```

**4. WhatsApp QR Code Not Showing**
```
WhatsApp client not paired
```
**Solution:**
- Access WhatsApp admin panel: http://localhost:8080/admin/whatsapp
- Scan QR code with WhatsApp app (Linked Devices)
- Check `whatsapp_sessions/` directory is writable

**5. LLM API Errors**
```
Error: 401 Unauthorized - Invalid API key
```
**Solution:**
- Verify `LLM_API_KEY` in `.env`
- Check API key is valid at https://z.ai/
- Ensure `ZAI_ENDPOINT` is correct

**6. Tailwind CSS Not Applied**
```
Styles not loading on frontend
```
**Solution:**
```bash
# Rebuild CSS
npm run css:build

# Check output file exists
ls -lh static/css/output.css

# Watch for changes during development
npm run css:watch
```

---

## 📦 Makefile Commands Reference

| Command | Description | Example |
|---------|-------------|---------|
| `make help` | Show all available commands | `make help` |
| `make dev` | Run with hot reload (Air) | `make dev` |
| `make build` | Build production binary | `make build` |
| `make run` | Run without hot reload | `make run` |
| `make test` | Run all tests | `make test` |
| `make migrate-up` | Apply database migrations | `make migrate-up` |
| `make migrate-down` | Rollback last migration | `make migrate-down` |
| `make migrate-create` | Create new migration | `make migrate-create name=add_table` |
| `make docker-up` | Start Docker containers | `make docker-up` |
| `make docker-down` | Stop Docker containers | `make docker-down` |
| `make tidy` | Tidy Go modules | `make tidy` |
| `make clean` | Clean build artifacts | `make clean` |

---

## 🔄 Development Lifecycle

### Daily Development Flow

```bash
# Morning: Start development
make docker-up          # Start PostgreSQL
npm run css:watch       # Start CSS watcher in terminal 1
make dev                # Start Go server in terminal 2

# During development:
# - Edit code (Air auto-reloads)
# - Edit templates (refresh browser)
# - Edit Tailwind classes (CSS auto-rebuilds)

# Evening: Stop services
Ctrl+C                  # Stop Air and CSS watcher
make docker-down        # Stop Docker containers
```

### Feature Development Flow

```bash
# 1. Create feature branch
git checkout -b feature/add-booking-system

# 2. Create database migration
make migrate-create name=create_bookings_table

# 3. Edit migration files
vim migrations/000009_create_bookings_table.up.sql
vim migrations/000009_create_bookings_table.down.sql

# 4. Run migration
make migrate-up

# 5. Implement feature (handler → service → repository)
# ...

# 6. Test manually
curl http://localhost:8080/api/bookings

# 7. Write unit tests
# ...

# 8. Run tests
make test

# 9. Build and verify
make build
./bin/api

# 10. Commit and push
git add .
git commit -m "Add booking system feature"
git push origin feature/add-booking-system
```

---

## 🚀 Next Steps

Once development environment is set up:

1. **Explore the codebase:** Start with `cmd/api/main.go`
2. **Review architecture:** Read `docs/source-tree-analysis.md`
3. **Test API endpoints:** Use curl or Postman
4. **Pair WhatsApp bot:** Access `/admin/whatsapp`
5. **Add a new feature:** Follow the patterns in existing handlers

**Happy coding! 🎉**
