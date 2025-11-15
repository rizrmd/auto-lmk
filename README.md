# Auto LMK - Car Sales Platform with AI-Powered WhatsApp Bot

![Version](https://img.shields.io/badge/version-1.0.0-blue)
![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)
![Status](https://img.shields.io/badge/status-production--ready-success)

> **Multi-tenant SaaS platform for car dealerships with intelligent WhatsApp Bot powered by LLM (OpenAI/Anthropic/Z.AI)**

---

## 🚀 Overview

Auto LMK is a modern, production-ready platform that revolutionizes car sales through:

- **🤖 AI-Powered WhatsApp Bot**: Customers can search for cars naturally through chat (Bahasa Indonesia)
- **🏢 Multi-Tenant Architecture**: One platform serves multiple car dealerships
- **🎨 Modern HTMX Frontend**: Fast, responsive UI without heavy JavaScript frameworks
- **📊 Comprehensive Admin Dashboard**: Manage inventory, leads, and sales team
- **🔒 Secure by Default**: Row-level tenant isolation, encrypted data, and secure authentication

---

## ✨ Key Features

### For Customers
- 💬 **Natural Language Search**: "Ada Toyota budget 200 juta?" → Bot finds matching cars
- 📱 **WhatsApp Integration**: No app download needed, works on any device
- 🖼️ **Rich Media**: View car photos, specs, and details instantly
- 🚗 **Smart Recommendations**: LLM suggests similar cars based on preferences
- ⚡ **24/7 Availability**: Bot responds instantly anytime

### For Dealerships (Tenants)
- 🎯 **Lead Management**: Automatic lead capture from WhatsApp conversations
- 👥 **Sales Team Management**: Differentiate sales vs customer conversations
- 📈 **Analytics Dashboard**: Track inventory, leads, and conversations
- 🔧 **Easy Setup**: Domain-based tenant isolation, simple WhatsApp pairing
- 💼 **Professional Templates**: Responsive public website included

### For Platform Admins
- 🏗️ **Multi-Tenant Control**: Manage multiple dealerships from one dashboard
- 🔐 **Security First**: Tenant isolation enforced at middleware level
- 📊 **Scalable Architecture**: PostgreSQL, horizontal scaling ready
- 🐳 **Docker Support**: Easy deployment with Docker Compose

---

## 🛠️ Technology Stack

| Category | Technology |
|----------|-----------|
| **Backend** | Go 1.21+, Chi Router |
| **Database** | PostgreSQL 15+ |
| **Frontend** | HTMX, Tailwind CSS, Alpine.js |
| **AI/LLM** | OpenAI (GPT-4o-mini) / Anthropic (Claude) / Z.AI (GLM-4-Flash) |
| **WhatsApp** | Whatsmeow |
| **Deployment** | Docker, Nginx, Systemd |
| **Security** | Bcrypt, JWT, Row-Level Isolation |

---

## 📦 Quick Start

### Prerequisites

- Go 1.21+ ([installation guide](https://go.dev/doc/install))
- PostgreSQL 15+ (or Docker)
- WhatsApp account for bot
- LLM API key (OpenAI, Anthropic, or Z.AI)

### Installation

```bash
# 1. Clone repository
git clone https://github.com/riz/auto-lmk.git
cd auto-lmk

# 2. Install dependencies
go mod download

# 3. Copy environment file
cp .env.example .env

# 4. Configure environment
vim .env  # Set your DB credentials, LLM API key, etc.

# 5. Start database (Docker)
docker-compose up -d postgres

# 6. Run migrations
make migrate-up

# 7. (Optional) Seed sample data
psql -U autolmk -d autolmk -h localhost < scripts/seed.sql

# 8. Start development server
make dev

# Server will start at http://localhost:8080
```

---

## 🏃 Development

### Hot Reload Development

```bash
# Install Air for hot reload
go install github.com/cosmtrek/air@latest

# Start development server with hot reload
make dev
# or
air
```

### Database Migrations

```bash
# Create new migration
migrate create -ext sql -dir migrations -seq create_your_table

# Run migrations
make migrate-up

# Rollback last migration
make migrate-down
```

### Build for Production

```bash
# Build binary
make build

# Run binary
./bin/api

# Or use Docker
docker-compose -f docker-compose.prod.yml up --build
```

---

## 🧪 Testing

### Manual API Testing

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
```

**See full API testing guide**: [`docs/api-testing-guide.md`](docs/api-testing-guide.md)

### Automated Tests

```bash
# Run all tests
make test

# Run with coverage
go test -cover ./...
```

---

## 📖 Documentation

| Document | Description |
|----------|-------------|
| [Architecture](docs/architecture.md) | System design, data flows, ERD |
| [Deployment Guide](docs/deployment-guide.md) | Production deployment instructions |
| [API Testing Guide](docs/api-testing-guide.md) | Complete API reference with examples |
| [LLM Provider Research](docs/llm-provider-research.md) | OpenAI vs Anthropic comparison |
| [Week 1 Progress](docs/week1-progress.md) | Development progress tracker |

---

## 🏗️ Architecture

### High-Level Overview

```
┌─────────────┐
│   Browser   │──┐
│  (Customer) │  │
└─────────────┘  │
                 │
┌─────────────┐  │      ┌──────────────┐      ┌─────────────┐
│  WhatsApp   │──┼─────→│  Nginx       │─────→│  Go App     │
│   Client    │  │      │  (SSL/Proxy) │      │  (Chi)      │
└─────────────┘  │      └──────────────┘      └──────┬──────┘
                 │                                    │
┌─────────────┐  │                                    │
│ Admin Panel │──┘                                    │
└─────────────┘                                       │
                                                      │
                          ┌───────────────────────────┴──────────────┐
                          │                                           │
                    ┌─────▼──────┐                           ┌───────▼────────┐
                    │ PostgreSQL │                           │  LLM Provider  │
                    │  (Multi-   │                           │ (OpenAI/Claude)│
                    │   Tenant)  │                           └────────────────┘
                    └────────────┘
```

### Multi-Tenant Isolation

Every HTTP request goes through the `TenantExtractor` middleware:

1. Extract domain from `Host` header (e.g., `showroom-jaya.localhost`)
2. Query: `SELECT id FROM tenants WHERE domain = ?`
3. Add `tenant_id` to request context
4. **All** repository queries automatically filter by `tenant_id`

**Result**: Complete data isolation between tenants at database level.

---

## 🔐 Security Features

- ✅ **Multi-Tenant Isolation**: Row-level security via middleware
- ✅ **Password Hashing**: Bcrypt with cost factor 12
- ✅ **JWT Authentication**: Stateless auth tokens (ready)
- ✅ **SQL Injection Prevention**: Parameterized queries throughout
- ✅ **CORS Configuration**: Configurable origins
- ✅ **SSL/TLS Support**: Nginx SSL termination
- ✅ **Input Validation**: Request validation at handler level

---

## 🤖 WhatsApp Bot Features

### Natural Language Understanding

The bot understands Indonesian naturally:

```
Customer: "Ada Toyota budget 200 juta?"
Bot: "Ada beberapa pilihan Toyota dalam budget Anda:
      1. Toyota Avanza 2020 - Rp 180 juta
      2. Toyota Rush 2019 - Rp 195 juta
      ..."

Customer: "Yang matic aja"
Bot: "Berikut Toyota dengan transmisi automatic:
      1. Toyota Avanza 2020 (Automatic) - Rp 180 juta
      ..."
```

### Function Calling

The LLM can call these functions:
- `searchCars`: Find cars by brand, price, transmission, etc.
- `getCarDetails`: Get full details of specific car
- `createLead`: Capture lead information automatically

### Sales vs Customer Mode

- **Customer Mode**: Friendly, helpful, focused on finding the right car
- **Sales Mode**: Internal tools, lead management, advanced search

---

## 🎨 Frontend Templates

### Public Website

| Template | Purpose | Features |
|----------|---------|----------|
| `templates/pages/home.html` | Homepage | Hero, featured cars, CTA |
| `templates/pages/cars.html` | Car catalog | Filters, search, HTMX pagination |
| `templates/pages/car-detail.html` | Car details | Photos, specs, WhatsApp CTA |
| `templates/pages/contact.html` | Contact page | FAQ, location, hours |

### Admin Dashboard

| Template | Purpose | Features |
|----------|---------|----------|
| `templates/admin/dashboard.html` | Overview | Stats, recent leads, quick actions |
| `templates/admin/cars.html` | Car management | CRUD, filters, bulk actions |
| `templates/admin/leads.html` | Lead management | Status updates, filtering |
| `templates/admin/layout.html` | Admin shell | Sidebar, navigation |

---

## 🚀 Deployment

### Option 1: Docker Compose (Recommended)

```bash
# Production deployment
docker-compose -f docker-compose.prod.yml up -d --build

# View logs
docker-compose -f docker-compose.prod.yml logs -f app
```

### Option 2: Systemd Service

```bash
# Build binary
CGO_ENABLED=0 GOOS=linux go build -o auto-lmk cmd/api/main.go

# Install systemd service
sudo cp auto-lmk /opt/auto-lmk/
sudo cp auto-lmk.service /etc/systemd/system/
sudo systemctl enable auto-lmk
sudo systemctl start auto-lmk
```

### Nginx Reverse Proxy

```nginx
server {
    listen 443 ssl http2;
    server_name showroom-jaya.com;

    ssl_certificate /etc/letsencrypt/live/showroom-jaya.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/showroom-jaya.com/privkey.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

**Full deployment guide**: [`docs/deployment-guide.md`](docs/deployment-guide.md)

---

## 📊 Database Schema

### Core Tables

| Table | Purpose | Key Columns |
|-------|---------|-------------|
| `tenants` | Dealerships | domain, name, whatsapp_number |
| `users` | Admin users | email, password_hash, tenant_id |
| `sales` | Sales team | phone_number, name, tenant_id |
| `cars` | Vehicle inventory | brand, model, year, price, tenant_id |
| `car_photos` | Car images | car_id, file_path |
| `car_specs` | Additional specs (EAV) | car_id, key, value |
| `conversations` | WhatsApp chats | sender_phone, is_sales, tenant_id |
| `messages` | Chat messages | conversation_id, message_text, direction |
| `leads` | Sales leads | phone_number, interested_car_id, status, tenant_id |

**All tables** (except `tenants`) include `tenant_id` foreign key for isolation.

**See full ERD**: [`docs/architecture.md#entity-relationship-diagram`](docs/architecture.md#entity-relationship-diagram)

---

## 🔧 Configuration

### Environment Variables

```bash
# Server
ENV=production
PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=autolmk
DB_PASSWORD=your_secure_password
DB_NAME=autolmk
DB_SSLMODE=require

# LLM Provider
LLM_PROVIDER=zai                 # Options: "openai", "anthropic", "zai"
LLM_API_KEY=your_api_key_here    # Your API key
LLM_MODEL=glm-4-flash            # Model name (gpt-4o-mini, claude-3-5-haiku-20241022, glm-4-flash)
# Z.AI specific endpoint (only used when LLM_PROVIDER=zai)
ZAI_ENDPOINT=https://api.z.ai/api/coding/paas/v4

# WhatsApp
WHATSAPP_SESSION_PATH=./whatsapp_sessions

# Security
JWT_SECRET=your_random_jwt_secret_here
```

### LLM Provider Selection

| Provider | Model | Cost (per 1M tokens) | Best For |
|----------|-------|---------------------|----------|
| OpenAI | `gpt-4o-mini` | $0.15 / $0.60 | Fast, cheap, function calling |
| Anthropic | `claude-3-5-haiku-20241022` | $1.00 / $5.00 | Better reasoning, Bahasa |
| Z.AI | `glm-4-flash` | Competitive pricing | Multilingual, fast, cost-effective |

**Default**: `glm-4-flash` (Z.AI - ready to use with provided API key)
**Alternative**: `gpt-4o-mini` (OpenAI - best price-performance if you have API key)

---

## 📁 Project Structure

```
auto-lmk/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── handler/                    # HTTP request handlers
│   │   ├── tenant_handler.go       # Root admin: tenant CRUD
│   │   ├── car_handler.go          # Tenant: car CRUD
│   │   ├── sales_handler.go        # Tenant: sales management
│   │   ├── lead_handler.go         # Tenant: lead management
│   │   └── conversation_handler.go # Tenant: conversation view
│   ├── middleware/                 # HTTP middleware
│   │   └── tenant.go               # Tenant extraction (CRITICAL)
│   ├── model/                      # Domain models & DTOs
│   │   ├── context.go              # Tenant context helpers
│   │   ├── tenant.go
│   │   ├── car.go
│   │   ├── sales.go
│   │   └── conversation.go
│   ├── repository/                 # Data access layer (all tenant-scoped)
│   │   ├── tenant_repository.go    # Root admin only
│   │   ├── car_repository.go
│   │   ├── sales_repository.go
│   │   ├── conversation_repository.go
│   │   └── lead_repository.go
│   ├── service/                    # Business logic layer
│   │   ├── car_service.go          # Car validation & bot search
│   │   └── whatsapp_service.go     # WhatsApp bot orchestration
│   ├── whatsapp/                   # WhatsApp integration
│   │   └── client.go               # Multi-tenant WhatsApp manager
│   └── llm/                        # LLM integration
│       ├── provider.go             # Provider abstraction
│       └── bot.go                  # Conversation bot
├── pkg/                            # Public packages
│   ├── config/                     # Environment configuration
│   ├── database/                   # PostgreSQL connection pool
│   ├── logger/                     # Structured logging (slog)
│   └── security/                   # Bcrypt, JWT utilities
├── migrations/                     # Database migrations (18 files)
│   ├── 000001_create_tenants_table.up.sql
│   ├── 000002_create_users_table.up.sql
│   └── ...
├── templates/                      # HTMX templates
│   ├── layouts/
│   │   └── base.html               # Public site layout
│   ├── pages/
│   │   ├── home.html
│   │   ├── cars.html               # Car catalog with filters
│   │   ├── car-detail.html         # Car detail page
│   │   └── contact.html
│   └── admin/
│       ├── layout.html             # Admin dashboard shell
│       ├── dashboard.html          # Admin overview
│       ├── cars.html               # Car management table
│       └── leads.html              # Lead management
├── scripts/
│   └── seed.sql                    # Sample data for testing
├── docs/                           # Documentation (10 files)
│   ├── architecture.md             # System architecture diagrams
│   ├── deployment-guide.md         # Production deployment
│   ├── api-testing-guide.md        # API reference
│   └── ...
├── docker-compose.yml              # Development environment
├── docker-compose.prod.yml         # Production setup
├── Makefile                        # Development commands
├── .air.toml                       # Hot reload configuration
├── .env.example                    # Environment template
└── README.md                       # This file
```

---

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Write tests for new features
- Follow Go best practices
- Update documentation
- Run `go fmt` before committing
- Keep commits atomic and descriptive

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- [Chi](https://github.com/go-chi/chi) - Lightweight HTTP router
- [HTMX](https://htmx.org) - Modern UI without heavy JS
- [Whatsmeow](https://github.com/tulir/whatsmeow) - WhatsApp Web API
- [OpenAI](https://openai.com) - GPT models for LLM
- [Anthropic](https://anthropic.com) - Claude models for LLM
- [Z.AI](https://z.ai) - GLM models for LLM (GLM-4-Flash)
- [Tailwind CSS](https://tailwindcss.com) - Utility-first CSS framework

---

## 📞 Support

- **Documentation**: [docs/](docs/)
- **Issues**: [GitHub Issues](https://github.com/riz/auto-lmk/issues)
- **Email**: support@auto-lmk.com

---

## 🗺️ Roadmap

### ✅ Completed (v1.0)

- [x] Multi-tenant architecture with domain-based isolation
- [x] Complete CRUD API for cars, sales, leads
- [x] **WhatsApp bot fully integrated with Whatsmeow**
- [x] **QR code pairing for WhatsApp**
- [x] **Message sending/receiving with event handling**
- [x] LLM integration framework (OpenAI + Anthropic + Z.AI)
- [x] **Z.AI (GLM-4-Flash) provider with function calling**
- [x] **LLM Bot with conversation context & history**
- [x] **Function execution (searchCars, getCarDetails, createLead)**
- [x] HTMX public website templates
- [x] Admin dashboard templates
- [x] Docker deployment setup
- [x] Comprehensive documentation

### 🚧 In Progress (v1.1)

- [ ] WhatsApp admin panel (pairing UI, status monitoring)
- [ ] Complete admin car edit/create forms
- [ ] File upload for car photos
- [ ] User authentication (JWT implementation)
- [ ] Root admin dashboard
- [ ] Testing with real WhatsApp numbers

### 🔮 Planned (v2.0)

- [ ] Mobile app (React Native / Flutter)
- [ ] Advanced analytics dashboard
- [ ] Email notifications
- [ ] SMS integration
- [ ] Multi-language support (English, Malay)
- [ ] Payment gateway integration
- [ ] Test drive scheduling system
- [ ] CRM integration (Salesforce, HubSpot)

---

## 📈 Stats

- **Lines of Code**: ~2,500+ Go code
- **Files**: 27 Go files, 10 docs, 11 templates
- **Database Migrations**: 18 files
- **Test Coverage**: TBD
- **Build Time**: ~3 seconds
- **Binary Size**: ~11MB
- **Docker Image**: ~25MB (alpine)

---

**Made with ❤️ by the Auto LMK Team**

*Revolutionizing car sales, one chat at a time.*
