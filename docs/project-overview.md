# Project Overview - Auto LMK

> **Generated:** 2025-11-15
> **Status:** Production-Ready (v1.0)

---

## 🎯 Project Summary

**Auto LMK** adalah platform SaaS multi-tenant untuk showroom mobil dengan WhatsApp bot AI-powered yang memungkinkan customer mencari mobil secara natural melalui chat dalam Bahasa Indonesia.

### Quick Facts

| Attribute | Value |
|-----------|-------|
| **Name** | Auto LMK |
| **Type** | Multi-Tenant SaaS Platform |
| **Domain** | Car Sales / Automotive |
| **Status** | Production-Ready (v1.0) |
| **Architecture** | Backend Monolith + Server-Rendered Frontend |
| **Primary Language** | Go 1.25.3 |
| **Database** | PostgreSQL 15+ |
| **Lines of Code** | ~2,500+ Go code |
| **Repository** | github.com/riz/auto-lmk |

---

## 🌟 What Auto LMK Does

### For Customers
- 💬 **Natural Language Car Search via WhatsApp**
  - "Ada Toyota budget 200 juta?" → Bot finds matching cars
  - 24/7 availability, instant responses
  - Rich media (car photos, specs)
  - Smart recommendations from LLM

### For Dealerships (Tenants)
- 🎯 **Lead Management:** Automatic lead capture from WhatsApp
- 👥 **Sales Team Management:** Track sales team conversations
- 📈 **Inventory Management:** CRUD for car inventory
- 📊 **Analytics Dashboard:** Monitor leads and conversations
- 🔧 **Easy Setup:** Domain-based tenant isolation, simple WhatsApp pairing

### For Platform Admins
- 🏗️ **Multi-Tenant Control:** Manage multiple dealerships
- 🔐 **Security:** Row-level tenant isolation
- 📊 **Scalable:** PostgreSQL, horizontal scaling ready
- 🐳 **Docker Support:** Easy deployment

---

## 🏗️ Technology Stack Summary

| Category | Technologies |
|----------|-------------|
| **Backend** | Go 1.25.3, Chi Router v5 |
| **Database** | PostgreSQL 15-alpine |
| **Frontend** | HTMX, Tailwind CSS v4, Alpine.js |
| **AI/LLM** | Z.AI (glm-4.6) - OpenAI-compatible |
| **WhatsApp** | Whatsmeow (WhatsApp Web API) |
| **Build Tools** | Vite 7.2.2, golang-migrate, Air |
| **Security** | Bcrypt, JWT, golang.org/x/crypto |
| **Deployment** | Docker, Docker Compose, Nginx, Systemd |

---

## 🏛️ Architecture Type

**Pattern:** Clean Architecture / Layered Architecture

```
Handler → Service → Repository → Database
```

**Key Characteristics:**
- ✅ Multi-tenant with domain-based isolation
- ✅ Monolith (single cohesive codebase)
- ✅ Server-rendered templates (HTMX)
- ✅ RESTful API + HTML endpoints
- ✅ Stateless design (horizontally scalable)

---

## 📊 Repository Structure

### Project Organization

| Component | Location | Purpose |
|-----------|----------|---------|
| **Entry Point** | `cmd/api/main.go` | HTTP server bootstrap |
| **Handlers** | `internal/handler/` | HTTP request handling (6 handlers) |
| **Services** | `internal/service/` | Business logic (2 services) |
| **Repositories** | `internal/repository/` | Data access (4 repositories) |
| **Models** | `internal/model/` | Domain entities (5 models) |
| **Middleware** | `internal/middleware/` | Tenant extraction (CRITICAL) |
| **WhatsApp** | `internal/whatsapp/` | WhatsApp client manager |
| **LLM** | `internal/llm/` | LLM provider abstraction (3 files) |
| **Packages** | `pkg/` | Reusable utilities (config, db, logger, security) |
| **Migrations** | `migrations/` | Database schema (16 files for 8 tables) |
| **Templates** | `templates/` | Server-rendered UI (17 HTML files) |
| **Static Assets** | `static/` | CSS, JS, images |

**File Counts:**
- Go files: 27
- Templates: 17
- Migrations: 16 (8 up, 8 down)
- Documentation: 6+ files

---

## 🗄️ Database Schema Overview

### Tables (8 Total)

| Table | Purpose | Records (Est.) | Tenant-Scoped |
|-------|---------|----------------|---------------|
| `tenants` | Dealerships | 10-100 | ❌ (Root) |
| `users` | Admin users | 100-1000 | ✅ |
| `sales` | Sales team | 10-100/tenant | ✅ |
| `cars` | Vehicle inventory | 50-500/tenant | ✅ |
| `car_specs` | Car specifications (EAV) | 200-2000 | ✅ |
| `car_photos` | Car images | 250-2500 | ✅ |
| `conversations` | WhatsApp chats | 1000-10000/tenant | ✅ |
| `messages` | Chat messages | 10000-100000/tenant | ✅ |

**Total Database Size (Estimated):** 1-10 GB for 10 tenants

---

## 🔌 API Overview

### Endpoint Categories

**Public Pages (No Auth):**
- `/` - Homepage
- `/cars` - Car catalog with filters
- `/cars/:id` - Car details
- `/contact` - Contact page

**Root Admin API:**
- `/api/root/tenants` - Tenant CRUD (super admin only)

**Tenant-Scoped API:**
- `/api/cars` - Car management
- `/api/sales` - Sales team management
- `/api/conversations` - WhatsApp conversation viewing
- `/admin/whatsapp` - WhatsApp pairing & management

**Total Endpoints:** ~15-20

---

## 🚀 Deployment Options

### Development
```bash
make dev                # Hot reload with Air
docker-compose up -d    # PostgreSQL
```

### Production (Docker)
```bash
docker-compose -f docker-compose.prod.yml up --build
```

### Production (Systemd)
```bash
make build
systemctl start auto-lmk
```

**Reverse Proxy:** Nginx with SSL/TLS

---

## 🔐 Security Features

| Feature | Status | Implementation |
|---------|--------|----------------|
| Multi-Tenant Isolation | ✅ | Domain-based middleware |
| Password Hashing | ✅ | Bcrypt (cost factor 12) |
| JWT Authentication | ⏳ | Ready, not enforced |
| SQL Injection Prevention | ✅ | Parameterized queries |
| CORS | ✅ | go-chi/cors |
| HTTPS | ✅ | Nginx SSL termination |

---

## 📈 Development Status

### ✅ Completed (v1.0)

- [x] Multi-tenant architecture with domain-based isolation
- [x] Complete CRUD API for cars, sales, conversations
- [x] WhatsApp bot integration (Whatsmeow)
- [x] QR code pairing for WhatsApp
- [x] LLM integration (Z.AI glm-4.6)
- [x] Function calling (searchCars, getCarDetails, createLead)
- [x] HTMX public website templates
- [x] Admin dashboard templates
- [x] Docker deployment setup
- [x] Database migrations (8 tables)
- [x] Comprehensive README

### 🚧 In Progress (v1.1)

- [ ] Complete admin car edit/create forms
- [ ] File upload for car photos
- [ ] User authentication (JWT implementation)
- [ ] Root admin dashboard
- [ ] Comprehensive testing

### 🔮 Planned (v2.0)

- [ ] Mobile app (React Native / Flutter)
- [ ] Advanced analytics dashboard
- [ ] Email/SMS notifications
- [ ] Multi-language support (English, Malay)
- [ ] Payment gateway integration
- [ ] Test drive scheduling system
- [ ] CRM integration (Salesforce, HubSpot)

---

## 👥 Target Users

### Primary Users

1. **Car Dealership Owners**
   - Manage inventory
   - View leads and analytics
   - Configure WhatsApp bot

2. **Sales Team Members**
   - Use WhatsApp for internal communication
   - Access lead information
   - Check inventory availability

3. **Customers**
   - Search for cars via WhatsApp
   - View car details and photos
   - Request information, schedule visits

4. **Platform Administrators**
   - Manage multiple tenants
   - Monitor system health
   - Configure global settings

---

## 🎯 Business Value

### Key Differentiators

1. **AI-Powered Natural Language Search**
   - No app download needed
   - Works on any device (WhatsApp)
   - Bahasa Indonesia support

2. **Multi-Tenant SaaS**
   - One platform, multiple dealerships
   - Complete data isolation
   - Cost-effective scaling

3. **Instant Lead Capture**
   - Automatic lead extraction from conversations
   - No manual data entry
   - Integration-ready

4. **24/7 Availability**
   - Bot responds instantly anytime
   - Reduces sales team workload
   - Improves customer experience

---

## 📊 Technical Metrics

| Metric | Value |
|--------|-------|
| **API Response Time** | < 100ms (target) |
| **Page Load Time** | < 2s (target) |
| **WhatsApp Reply Time** | ~2-3s (actual) |
| **Binary Size** | ~11MB |
| **Docker Image** | ~25MB (alpine) |
| **Build Time** | ~3 seconds |
| **Concurrent Users** | 1000+ (estimated) |

---

## 🔗 Related Documentation

| Document | Purpose |
|----------|---------|
| [README.md](../README.md) | Main project documentation |
| [architecture.md](./architecture.md) | System architecture & design |
| [source-tree-analysis.md](./source-tree-analysis.md) | Complete directory structure |
| [development-guide.md](./development-guide.md) | Development setup & workflow |

---

## 🚦 Getting Started (Quick Links)

### For Developers
1. **Setup:** See [development-guide.md](./development-guide.md#quick-start)
2. **Architecture:** See [architecture.md](./architecture.md)
3. **Code Structure:** See [source-tree-analysis.md](./source-tree-analysis.md)

### For Product Managers
1. **Features:** See [README.md](../README.md#key-features)
2. **Roadmap:** See [README.md](../README.md#roadmap)

### For DevOps
1. **Deployment:** See [development-guide.md](./development-guide.md#deployment-architecture)
2. **Docker:** See [README.md](../README.md#deployment)

---

## 📞 Support & Contact

- **Repository:** https://github.com/riz/auto-lmk
- **Issues:** https://github.com/riz/auto-lmk/issues
- **Email:** support@auto-lmk.com

---

**Auto LMK - Revolutionizing car sales, one chat at a time.** 🚗💬🤖
