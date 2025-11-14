# Auto LMK - Multi-Tenant Car Sales SaaS Platform

Platform penjualan mobil multi-tenant dengan fitur WhatsApp bot menggunakan Whatsmeow dan LLM integration.

## 🚀 Tech Stack

- **Backend:** Go (Golang)
- **Frontend:** HTMX
- **Database:** PostgreSQL
- **WhatsApp Integration:** Whatsmeow
- **LLM:** TBD (OpenAI/Anthropic/Local)

## 📁 Project Structure

```
auto-lmk/
├── cmd/
│   └── api/              # Main application entry point
├── internal/
│   ├── handler/          # HTTP handlers
│   ├── middleware/       # Middleware (tenant isolation, auth, etc)
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic
│   └── model/            # Domain models
├── pkg/
│   ├── config/           # Configuration management
│   ├── database/         # Database connection & utilities
│   ├── logger/           # Logging utilities
│   └── security/         # Security utilities
├── migrations/           # Database migrations
├── docs/                 # Documentation
└── docker-compose.yml    # Local development environment
```

## 🏗️ Setup Instructions

### Prerequisites

- Go 1.21 or higher
- Docker & Docker Compose
- Make (optional)

### Local Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/riz/auto-lmk.git
   cd auto-lmk
   ```

2. **Start PostgreSQL with Docker Compose**
   ```bash
   docker-compose up -d
   ```

3. **Create .env file**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Install Go dependencies**
   ```bash
   go mod download
   ```

5. **Run database migrations**
   ```bash
   # Migration tool TBD
   ```

6. **Run the application**
   ```bash
   go run cmd/api/main.go
   ```

## 🔧 Development

### Database

- **Host:** localhost
- **Port:** 5432
- **Database:** autolmk
- **User:** autolmk
- **Password:** autolmk_dev_password

### Environment Variables

Create a `.env` file in the root directory:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=autolmk
DB_PASSWORD=autolmk_dev_password
DB_NAME=autolmk

# Server
PORT=8080
ENV=development

# LLM (TBD)
LLM_PROVIDER=
LLM_API_KEY=

# WhatsApp (TBD)
WHATSAPP_SESSION_PATH=./whatsapp_sessions
```

## 📊 Roadmap

Mengikuti Implementation Roadmap 10-minggu yang terletak di `docs/implementation-roadmap-2025-11-14.md`

### Current Phase: Week 1 - Setup & Foundation

- [x] Project initialization
- [x] Git repository setup
- [x] Docker Compose for PostgreSQL
- [x] Basic project structure
- [ ] Database selection & setup
- [ ] Core dependencies selection
- [ ] LLM provider research
- [ ] Whatsmeow investigation

## 🎯 Core Features (MVP)

### Multi-Tenancy
- Row-level tenant isolation dengan `tenant_id`
- Domain-based tenant identification
- Context propagation untuk security

### WhatsApp Bot
- Multi-number support (1 bot instance, multiple tenant numbers)
- LLM-powered natural language understanding
- Hybrid command parser (structured + natural language)
- Sales vs Customer differentiation

### Web Interface
- Public website per tenant (SEO-optimized)
- Tenant admin dashboard
- Root admin interface
- HTMX-based responsive UI

## 📝 License

TBD

## 👤 Author

BMad

---

**Status:** Week 1 - Foundation Setup 🏗️
