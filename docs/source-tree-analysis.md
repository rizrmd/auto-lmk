# Source Tree Analysis - Auto LMK

> **Generated:** 2025-11-15
> **Project:** Auto LMK - Car Sales Platform with AI-Powered WhatsApp Bot
> **Type:** Backend (Go + PostgreSQL) with Server-Rendered Frontend (HTMX)

---

## 📂 Complete Directory Structure (Annotated)

```
auto-lmk/
├── cmd/                              # Application entry points
│   └── api/
│       └── main.go                   # 🚀 PRIMARY ENTRY POINT - HTTP server bootstrap
│
├── internal/                         # Private application code (not importable by external packages)
│   ├── handler/                      # 🔌 HTTP request handlers (6 handlers)
│   │   ├── car_handler.go            # Car CRUD operations
│   │   ├── tenant_handler.go         # Root admin: tenant management
│   │   ├── sales_handler.go          # Sales team management
│   │   ├── conversation_handler.go   # View WhatsApp conversations
│   │   ├── whatsapp_handler.go       # WhatsApp bot pairing & management
│   │   └── page_handler.go           # Frontend page rendering
│   │
│   ├── service/                      # 💼 Business logic layer (2 services)
│   │   ├── car_service.go            # Car validation & bot search logic
│   │   └── whatsapp_service.go       # WhatsApp bot orchestration
│   │
│   ├── repository/                   # 💾 Data access layer (4 repositories)
│   │   ├── tenant_repository.go      # Tenant CRUD (root admin only)
│   │   ├── car_repository.go         # Car CRUD (tenant-scoped)
│   │   ├── sales_repository.go       # Sales team CRUD (tenant-scoped)
│   │   └── conversation_repository.go # Conversation CRUD (tenant-scoped)
│   │
│   ├── model/                        # 📦 Domain models & DTOs (5 models)
│   │   ├── tenant.go                 # Tenant entity
│   │   ├── car.go                    # Car entity & specs
│   │   ├── sales.go                  # Sales team entity
│   │   ├── conversation.go           # Conversation & message entities
│   │   └── context.go                # Tenant context helpers (multi-tenant)
│   │
│   ├── middleware/                   # 🛡️ HTTP middleware (CRITICAL)
│   │   └── tenant.go                 # Multi-tenant isolation via domain extraction
│   │
│   ├── whatsapp/                     # 📱 WhatsApp integration (1 file)
│   │   └── client.go                 # Multi-tenant WhatsApp manager (Whatsmeow)
│   │
│   └── llm/                          # 🤖 LLM integration (3 files)
│       ├── provider.go               # LLM provider abstraction (Z.AI/OpenAI/Anthropic)
│       ├── adapter.go                # LLM adapter pattern
│       └── bot.go                    # Conversation bot with function calling
│
├── pkg/                              # 📚 Public packages (reusable, importable)
│   ├── config/                       # ⚙️ Environment configuration
│   │   └── config.go                 # .env loader (godotenv)
│   │
│   ├── database/                     # 🗄️ Database connection
│   │   └── database.go               # PostgreSQL connection pool (lib/pq)
│   │
│   ├── logger/                       # 📝 Structured logging
│   │   └── logger.go                 # Logger setup (slog)
│   │
│   └── security/                     # 🔐 Security utilities
│       ├── password.go               # Bcrypt password hashing
│       └── jwt.go                    # JWT token generation & validation
│
├── migrations/                       # 🗃️ Database migrations (16 files for 8 tables)
│   ├── 000001_create_tenants_table.{up,down}.sql
│   ├── 000002_create_users_table.{up,down}.sql
│   ├── 000003_create_sales_table.{up,down}.sql
│   ├── 000004_create_cars_table.{up,down}.sql
│   ├── 000005_create_car_specs_table.{up,down}.sql
│   ├── 000006_create_car_photos_table.{up,down}.sql
│   ├── 000007_create_conversations_table.{up,down}.sql
│   └── 000008_create_messages_table.{up,down}.sql
│
├── templates/                        # 🎨 Server-rendered HTML templates (17 files)
│   ├── layouts/                      # Layout shells
│   │   └── base.html                 # Public site layout (HTMX)
│   │
│   ├── pages/                        # Public pages (4 pages)
│   │   ├── home.html                 # Homepage with hero
│   │   ├── cars.html                 # Car catalog with filters (HTMX pagination)
│   │   ├── car-detail.html           # Car details with WhatsApp CTA
│   │   └── contact.html              # Contact page
│   │
│   ├── admin/                        # Admin dashboard (4 files)
│   │   ├── layout.html               # Admin shell with sidebar
│   │   ├── dashboard.html            # Admin overview & stats
│   │   ├── cars.html                 # Car management table
│   │   └── whatsapp.html             # WhatsApp pairing UI
│   │
│   └── components/                   # Reusable UI components (8 components)
│       ├── nav.html                  # Navigation bar
│       ├── hero.html                 # Hero section
│       ├── footer.html               # Footer
│       ├── card.html                 # Car card
│       ├── button.html               # Button styles
│       ├── input.html                # Form input
│       ├── gallery.html              # Image gallery
│       └── pagination.html           # Pagination component
│
├── static/                           # 🌐 Static assets (CSS, JS, images)
│   ├── css/                          # Compiled CSS (Tailwind via Vite)
│   ├── js/                           # JavaScript files
│   └── images/                       # Static images
│
├── scripts/                          # 🛠️ Helper scripts
│   └── seed.sql                      # Sample data for testing
│
├── docs/                             # 📖 Documentation
│   ├── bmm-workflow-status.yaml      # BMM workflow tracking
│   ├── project-scan-report.json      # Scan state file
│   └── sprint-artifacts/             # Sprint planning artifacts
│
├── whatsapp_sessions/                # 📂 WhatsApp session data (multi-tenant)
│   └── [tenant-specific session files]
│
├── bin/                              # 🔨 Compiled binaries (gitignored)
│   └── api                           # Built Go binary
│
├── node_modules/                     # 📦 NPM dependencies (gitignored)
│   └── [Vite, Tailwind CSS]
│
├── Configuration Files (Root Level)
├── go.mod                            # Go module definition
├── go.sum                            # Go dependency checksums
├── package.json                      # NPM dependencies (Tailwind, Vite)
├── package-lock.json                 # NPM lockfile
├── tailwind.config.js                # Tailwind CSS configuration
├── vite.config.js                    # Vite build configuration
├── docker-compose.yml                # Docker setup (PostgreSQL)
├── Makefile                          # Development commands
├── .air.toml                         # Hot reload configuration
├── .env.example                      # Environment variable template
└── README.md                         # 📄 Main documentation (592 lines)
```

---

## 🔑 Critical Directories Explained

### 1. **`cmd/api/`** - Application Entry Point
**Purpose:** Bootstrap the HTTP server
**Key File:** `main.go`
**What it does:**
- Loads environment configuration
- Initializes database connection pool
- Sets up Chi router with middleware
- Registers all HTTP routes
- Starts WhatsApp clients for tenants
- Listens on configured port (default 8080)

**Entry Point Flow:**
```
main.go → Load Config → Connect DB → Setup Router → Register Routes → Start Server
```

---

### 2. **`internal/handler/`** - HTTP Request Handlers
**Purpose:** Handle incoming HTTP requests, validate input, call services
**Pattern:** Handler → Service → Repository → Database

**Handler Breakdown:**

| Handler | Routes | Scope | Purpose |
|---------|--------|-------|---------|
| `tenant_handler.go` | `/api/root/tenants` | Root Admin | Create/manage tenants |
| `car_handler.go` | `/api/cars` | Tenant | CRUD for car inventory |
| `sales_handler.go` | `/api/sales` | Tenant | Manage sales team |
| `conversation_handler.go` | `/api/conversations` | Tenant | View WhatsApp conversations |
| `whatsapp_handler.go` | `/admin/whatsapp` | Tenant | Pair WhatsApp, check status |
| `page_handler.go` | `/`, `/cars`, `/contact` | Public | Render public pages |

**Multi-Tenant Isolation:**
All tenant-scoped handlers rely on `middleware/tenant.go` to extract `tenant_id` from the request domain and inject it into the request context.

---

### 3. **`internal/middleware/tenant.go`** - Multi-Tenant Isolation (CRITICAL)
**Purpose:** Enforce row-level security at middleware level
**How it works:**
1. Extract domain from `Host` header (e.g., `showroom-jaya.localhost`)
2. Query database: `SELECT id FROM tenants WHERE domain = ?`
3. Add `tenant_id` to request context
4. All downstream repositories automatically filter by `tenant_id`

**Result:** Complete data isolation between tenants - no cross-tenant data leaks possible.

---

### 4. **`internal/service/`** - Business Logic Layer
**Purpose:** Encapsulate business rules, validation, and cross-cutting concerns

**Services:**
- **`car_service.go`**
  - Validates car data before saving
  - Implements car search logic for WhatsApp bot
  - Filters cars by brand, price range, transmission, etc.

- **`whatsapp_service.go`**
  - Orchestrates WhatsApp bot lifecycle (pairing, messaging, event handling)
  - Routes messages to LLM bot
  - Manages multi-tenant WhatsApp clients

---

### 5. **`internal/repository/`** - Data Access Layer
**Purpose:** Interact with PostgreSQL database, execute SQL queries
**Pattern:** All queries are **tenant-scoped** (filter by `tenant_id`)

**Repositories:**
- `tenant_repository.go` - Tenant CRUD (root admin only, no tenant filter)
- `car_repository.go` - Car CRUD (tenant-scoped)
- `sales_repository.go` - Sales team CRUD (tenant-scoped)
- `conversation_repository.go` - Conversation CRUD (tenant-scoped)

**Example Tenant-Scoped Query:**
```sql
SELECT * FROM cars WHERE tenant_id = $1 AND brand = $2
```

---

### 6. **`internal/model/`** - Domain Models
**Purpose:** Define data structures (entities, DTOs, request/response types)

**Key Models:**
- `tenant.go` - Tenant entity (domain, name, whatsapp_number)
- `car.go` - Car entity (brand, model, year, price, specs, photos)
- `sales.go` - Sales team member (phone_number, name, is_active)
- `conversation.go` - Conversation & Message entities
- `context.go` - Tenant context helpers (`GetTenantID`, `SetTenantID`)

---

### 7. **`internal/whatsapp/client.go`** - WhatsApp Integration
**Purpose:** Multi-tenant WhatsApp client manager using Whatsmeow
**Features:**
- Maintains separate WhatsApp sessions per tenant
- QR code generation for pairing
- Message sending/receiving
- Event handling (incoming messages, connection status)
- Session persistence in `whatsapp_sessions/` directory

**Session Storage:**
```
whatsapp_sessions/
└── tenant-{tenant_id}/
    ├── device.store
    └── session.json
```

---

### 8. **`internal/llm/`** - LLM Integration
**Purpose:** Abstract LLM provider (Z.AI, OpenAI, Anthropic) and implement conversation bot

**Architecture:**
- **`provider.go`** - Provider interface and factory
- **`adapter.go`** - Adapter pattern for different LLM APIs
- **`bot.go`** - Conversation bot with function calling

**Function Calling:**
The bot can call these functions:
- `searchCars(brand, maxPrice, transmission, ...)` - Find cars matching criteria
- `getCarDetails(carID)` - Get full car details
- `createLead(name, phone, carID)` - Capture lead information

**Current LLM:** Z.AI (glm-4.6) with OpenAI-compatible API

---

### 9. **`pkg/`** - Shared Packages
**Purpose:** Reusable packages that can be imported by external projects

| Package | Purpose | Key Functions |
|---------|---------|---------------|
| `config/` | Environment config | `Load()` - Load .env file |
| `database/` | DB connection | `Connect()` - PostgreSQL pool |
| `logger/` | Structured logging | `New()` - Create logger |
| `security/password` | Password hashing | `Hash()`, `Verify()` - Bcrypt |
| `security/jwt` | JWT tokens | `Generate()`, `Validate()` - JWT |

---

### 10. **`migrations/`** - Database Schema
**Purpose:** Version-controlled database schema changes
**Tool:** golang-migrate

**Database Schema (8 tables):**

| Table | Purpose | Key Columns | Tenant-Scoped? |
|-------|---------|-------------|----------------|
| `tenants` | Dealerships | id, domain, name, whatsapp_number | ❌ (root) |
| `users` | Admin users | id, email, password_hash, tenant_id | ✅ |
| `sales` | Sales team | id, phone_number, name, tenant_id | ✅ |
| `cars` | Vehicle inventory | id, brand, model, year, price, tenant_id | ✅ |
| `car_specs` | Additional specs (EAV) | id, car_id, key, value | ✅ (via car) |
| `car_photos` | Car images | id, car_id, file_path | ✅ (via car) |
| `conversations` | WhatsApp chats | id, sender_phone, is_sales, tenant_id | ✅ |
| `messages` | Chat messages | id, conversation_id, message_text, direction | ✅ (via conversation) |

**Migration Commands:**
```bash
make migrate-up      # Run all pending migrations
make migrate-down    # Rollback last migration
make migrate-create  # Create new migration file
```

---

### 11. **`templates/`** - Server-Rendered UI
**Purpose:** HTMX-powered server-rendered templates (no heavy JS framework)

**Template Organization:**

| Category | Files | Purpose |
|----------|-------|---------|
| **Layouts** | base.html, admin/layout.html | Page shells |
| **Public Pages** | home.html, cars.html, car-detail.html, contact.html | Customer-facing pages |
| **Admin Pages** | dashboard.html, cars.html, whatsapp.html | Admin management |
| **Components** | nav, hero, footer, card, button, input, gallery, pagination | Reusable UI elements |

**Technology:**
- **HTMX:** Dynamic updates without full page reloads
- **Tailwind CSS v4:** Utility-first styling (compiled by Vite)
- **Alpine.js:** Minimal JavaScript for interactivity (mentioned in README)

---

## 🔄 Request Flow

### Public Page Request (Customer)
```
Customer Browser
    ↓
Nginx (SSL termination)
    ↓
Go App :8080
    ↓
middleware/tenant.go → Extract tenant_id from domain
    ↓
handler/page_handler.go → Render template
    ↓
templates/pages/cars.html → HTMX partial updates
    ↓
Response (HTML)
```

### API Request (Admin/Sales)
```
Admin Dashboard
    ↓
POST /api/cars (create car)
    ↓
middleware/tenant.go → Extract tenant_id
    ↓
handler/car_handler.go → Validate input
    ↓
service/car_service.go → Business logic
    ↓
repository/car_repository.go → INSERT INTO cars (tenant_id, ...)
    ↓
PostgreSQL
    ↓
Response (JSON)
```

### WhatsApp Message Flow
```
Customer WhatsApp
    ↓
Whatsmeow Client (whatsapp/client.go)
    ↓
whatsapp_service.go → Route to LLM bot
    ↓
llm/bot.go → Process message with Z.AI (glm-4.6)
    ↓
Function Calling → searchCars() / getCarDetails()
    ↓
service/car_service.go → Search database
    ↓
repository/car_repository.go → SELECT * FROM cars WHERE tenant_id = ?
    ↓
Response → LLM formats reply
    ↓
whatsapp_service.go → Send WhatsApp message
    ↓
Customer WhatsApp (reply received)
```

---

## 🧩 Integration Points

Since this is a **monolith** (single cohesive codebase), all integration happens **within the same process**:

| Integration | How It Works |
|-------------|--------------|
| **HTTP → Service → Repository** | Direct function calls (in-process) |
| **WhatsApp → LLM → Database** | In-process function calls via whatsapp_service |
| **Frontend → Backend** | Same server renders templates and handles API calls |

**External Integrations:**
- **PostgreSQL:** Via `lib/pq` driver (TCP connection)
- **WhatsApp Web:** Via Whatsmeow library (WebSocket)
- **Z.AI LLM:** Via HTTPS API (OpenAI-compatible endpoint)

---

## 📦 Deployment Structure

### Development
```
make dev
    ↓
Air (hot reload)
    ↓
Go App :8080 (watches file changes)
    ↓
Docker Compose → PostgreSQL :5432
```

### Production (Docker)
```
docker-compose -f docker-compose.prod.yml up
    ↓
Nginx :443 (SSL)
    ↓
Go App (Docker container) :8080
    ↓
PostgreSQL (Docker container) :5432
```

### Production (Systemd)
```
Nginx :443 (SSL reverse proxy)
    ↓
Go Binary (systemd service) :8080
    ↓
PostgreSQL (system service) :5432
```

---

## 🔐 Security Architecture

### Multi-Tenant Isolation
**Enforcement Points:**
1. **Middleware Level:** `middleware/tenant.go` extracts tenant_id from domain
2. **Context Level:** `model/context.go` helpers inject tenant_id into queries
3. **Repository Level:** All queries filter by `WHERE tenant_id = $1`

**Result:** Row-level security enforced at application level.

### Authentication & Authorization
- **Password Hashing:** Bcrypt (cost factor 12) via `pkg/security/password.go`
- **JWT Tokens:** Stateless auth (ready but not fully implemented)
- **HTTPS:** SSL termination at Nginx (production)

### Input Validation
- **SQL Injection:** Prevented by parameterized queries (`$1`, `$2`, etc.)
- **CORS:** Configured via `go-chi/cors` middleware
- **Request Validation:** Handler-level validation before service layer

---

## 🎯 Key Takeaways for AI-Assisted Development

1. **Entry Point:** `cmd/api/main.go` is where the application starts
2. **Add New Feature:** Create handler → service → repository (follow existing pattern)
3. **Multi-Tenant:** ALWAYS filter by `tenant_id` in repositories
4. **Migrations:** Use `make migrate-create` for database changes
5. **Templates:** HTMX templates in `templates/` for frontend
6. **WhatsApp Bot:** Logic in `internal/llm/bot.go` and `internal/whatsapp/client.go`
7. **LLM Integration:** Z.AI (glm-4.6) via OpenAI-compatible API
8. **Testing:** `make test` to run Go tests
9. **Hot Reload:** `make dev` with Air for development

---

**📌 For Brownfield PRD:** This source tree provides complete context for planning new features, understanding data flow, and identifying reusable components.
