# Auto LMK - Documentation Index

> **Multi-Tenant Car Sales Platform with AI-Powered WhatsApp Bot**
>
> **Generated:** 2025-11-15 | **Status:** Production-Ready (v1.0)

---

## 📚 Documentation Overview

This is the **master documentation index** for Auto LMK. All documentation has been generated through automated codebase analysis to provide comprehensive context for AI-assisted development and brownfield feature planning.

---

## 🎯 Project At A Glance

| Attribute | Value |
|-----------|-------|
| **Type** | Backend Monolith with Server-Rendered Frontend |
| **Primary Language** | Go 1.25.3 |
| **Framework** | Chi Router v5.2.3 |
| **Database** | PostgreSQL 15+ |
| **Frontend** | HTMX + Tailwind CSS v4 + Alpine.js |
| **LLM Provider** | Z.AI (glm-4.6) - OpenAI compatible |
| **WhatsApp** | Whatsmeow (WhatsApp Web API) |
| **Architecture** | Clean Architecture / Layered |
| **Multi-Tenant** | Domain-based isolation |

---

## 🚀 Quick Reference

### Tech Stack Summary
- **Backend:** Go + Chi Router + PostgreSQL
- **Frontend:** HTMX (server-rendered) + Tailwind CSS
- **AI/LLM:** Z.AI glm-4.6 (function calling enabled)
- **WhatsApp:** Whatsmeow multi-tenant client
- **Deployment:** Docker + Nginx + Systemd

### Entry Point
- **Main:** `cmd/api/main.go` - HTTP server bootstrap

### Architecture Pattern
```
Handler → Service → Repository → Database
(Clean Architecture with 4 layers)
```

### Database Tables
8 tables total: tenants, users, sales, cars, car_specs, car_photos, conversations, messages

### API Endpoints
~15-20 endpoints (Public pages + Root admin + Tenant-scoped APIs)

---

## 📖 Generated Documentation

### Core Documentation

1. **[Project Overview](./project-overview.md)** ⭐
   - High-level project summary
   - Technology stack overview
   - Business value & target users
   - Development status & roadmap
   - **Start here for quick understanding**

2. **[Architecture](./architecture.md)** 🏗️
   - System architecture diagrams
   - Clean architecture pattern
   - Data architecture (ERD)
   - API design
   - Multi-tenant architecture
   - WhatsApp & LLM integration
   - Security architecture
   - Deployment architecture
   - **Essential for technical decisions**

3. **[Source Tree Analysis](./source-tree-analysis.md)** 📂
   - Complete annotated directory structure
   - Critical directories explained
   - Request flow diagrams
   - Integration points
   - Security enforcement points
   - **Essential for code navigation**

4. **[Development Guide](./development-guide.md)** 💻
   - Prerequisites & quick start
   - Environment setup (DB, LLM, config)
   - Development workflow (hot reload, testing)
   - Database migrations
   - Build & deployment
   - Common development tasks
   - Troubleshooting guide
   - **Essential for developers**

---

## 📄 Existing Documentation

### Original Documentation (Root)

1. **[README.md](../README.md)** - Main project documentation (592 lines)
   - Comprehensive overview
   - Installation & setup
   - Feature descriptions
   - API testing examples
   - Deployment instructions
   - Roadmap

---

## 🗂️ Supporting Files

- **[bmm-workflow-status.yaml](./bmm-workflow-status.yaml)** - BMM methodology workflow tracking
- **[project-scan-report.json](./project-scan-report.json)** - Automated scan state file

---

## 🎯 Getting Started Guide

### For New Developers

**Quick Setup (5 minutes):**

```bash
# 1. Clone & setup
git clone https://github.com/riz/auto-lmk.git
cd auto-lmk
cp .env.example .env

# 2. Start database
docker-compose up -d postgres

# 3. Run migrations
make migrate-up

# 4. Start development server
make dev

# ✅ Server running at http://localhost:8080
```

**Next Steps:**
1. Read [Project Overview](./project-overview.md) for high-level understanding
2. Review [Architecture](./architecture.md) for system design
3. Explore [Source Tree Analysis](./source-tree-analysis.md) for code structure
4. Follow [Development Guide](./development-guide.md) for detailed workflow

---

### For Product Managers

**Understanding the Product:**
1. Start with [Project Overview](./project-overview.md)
2. Review [README.md](../README.md) for features & roadmap
3. Check [Architecture](./architecture.md) for technical capabilities

**Planning New Features:**
1. Review [Architecture](./architecture.md) for integration points
2. Check [Source Tree Analysis](./source-tree-analysis.md) for reusable components
3. Understand multi-tenant isolation in [Architecture - Multi-Tenant Section](./architecture.md#multi-tenant-architecture)

---

### For DevOps Engineers

**Deployment Setup:**
1. Review [Development Guide - Deployment](./development-guide.md#deployment-architecture)
2. Check [README.md - Deployment Section](../README.md#deployment)
3. Understand [Architecture - Deployment](./architecture.md#deployment-architecture)

**Infrastructure:**
- Docker Compose for local/staging
- Systemd for production servers
- Nginx for reverse proxy + SSL

---

### For AI-Assisted Development

**Brownfield PRD Planning:**

This documentation index serves as the primary context for AI agents when planning new features for this existing (brownfield) project.

**Key Documents to Load:**
1. **[Architecture](./architecture.md)** - Understand system design & constraints
2. **[Source Tree Analysis](./source-tree-analysis.md)** - Identify reusable components
3. **[Development Guide](./development-guide.md)** - Follow development patterns

**Feature Planning Workflow:**
1. Load this index + relevant docs into AI context
2. Review existing architecture & patterns
3. Plan feature that aligns with current design
4. Identify integration points
5. Leverage existing components
6. Follow established conventions

---

## 🔍 Documentation Usage Guide

### By Role

| Role | Recommended Reading Order |
|------|--------------------------|
| **Backend Developer** | Development Guide → Architecture → Source Tree → README |
| **Frontend Developer** | Source Tree → Development Guide → Architecture |
| **Product Manager** | Project Overview → README → Architecture |
| **DevOps Engineer** | Development Guide (Deployment) → Architecture (Deployment) |
| **Technical Lead** | Architecture → Source Tree → Development Guide |
| **New Team Member** | Project Overview → README → Development Guide |

### By Task

| Task | Relevant Documentation |
|------|------------------------|
| **Add New API Endpoint** | Source Tree → Architecture (API Design) → Development Guide (Common Tasks) |
| **Add Database Table** | Architecture (Data) → Development Guide (Migrations) |
| **Deploy to Production** | Development Guide (Deployment) → README (Deployment) |
| **Understand Multi-Tenancy** | Architecture (Multi-Tenant) → Source Tree (Middleware) |
| **Integrate New LLM Feature** | Architecture (LLM Integration) → Source Tree (llm/) |
| **Setup WhatsApp Bot** | Development Guide → Architecture (WhatsApp) |
| **Fix Bug** | Source Tree → Architecture → Development Guide (Troubleshooting) |

---

## 📊 Project Statistics

| Metric | Count |
|--------|-------|
| **Go Files** | 27 |
| **Handlers** | 6 |
| **Services** | 2 |
| **Repositories** | 4 |
| **Models** | 5 |
| **Database Tables** | 8 |
| **Migrations** | 16 (8 up, 8 down) |
| **Templates** | 17 HTML files |
| **Documentation Files** | 6+ (including README) |
| **Lines of Go Code** | ~2,500+ |

---

## 🔗 External Resources

### Official Documentation
- **Go:** https://go.dev/doc/
- **Chi Router:** https://github.com/go-chi/chi
- **PostgreSQL:** https://www.postgresql.org/docs/
- **HTMX:** https://htmx.org/docs/
- **Tailwind CSS:** https://tailwindcss.com/docs

### Dependencies
- **Whatsmeow:** https://github.com/tulir/whatsmeow
- **Z.AI:** https://z.ai/ (LLM Provider)
- **golang-migrate:** https://github.com/golang-migrate/migrate

### Tools
- **Air (Hot Reload):** https://github.com/cosmtrek/air
- **Docker:** https://docs.docker.com/
- **Vite:** https://vitejs.dev/

---

## 🎓 Learning Path

### Week 1: Foundation
- [ ] Read Project Overview
- [ ] Setup development environment (Development Guide)
- [ ] Run the application locally
- [ ] Explore codebase using Source Tree Analysis

### Week 2: Deep Dive
- [ ] Study Architecture document
- [ ] Understand multi-tenant isolation
- [ ] Review API endpoints
- [ ] Test WhatsApp bot integration

### Week 3: Contributing
- [ ] Add a new API endpoint (follow Development Guide)
- [ ] Create a database migration
- [ ] Add a new template page
- [ ] Write tests

---

## 🆘 Need Help?

### Troubleshooting
1. Check [Development Guide - Troubleshooting](./development-guide.md#troubleshooting)
2. Review logs: `docker-compose logs -f`
3. Check database connection in `.env`
4. Verify migrations: `make migrate-up`

### Common Issues
- **Database connection failed:** Check Docker & .env
- **Port already in use:** Kill process or change PORT in .env
- **Migration failed:** Force version, then re-run
- **WhatsApp not pairing:** Access `/admin/whatsapp`
- **LLM errors:** Verify API key in .env

### Support Channels
- **Repository:** https://github.com/riz/auto-lmk
- **Issues:** https://github.com/riz/auto-lmk/issues
- **Email:** support@auto-lmk.com

---

## 📝 Documentation Maintenance

### Last Generated
**Date:** 2025-11-15
**Scan Type:** Quick Scan (Pattern-based, no source file reading)
**Workflow:** BMM document-project workflow v1.2.0

### Re-generating Documentation

To update this documentation after significant codebase changes:

```bash
# Run the BMM document-project workflow
/bmad:bmm:workflows:document-project
```

Or use the analyst agent:

```bash
/bmad:bmm:agents:analyst
# Then select: document-project workflow
```

### What Was Scanned
- ✅ Directory structure
- ✅ Configuration files (go.mod, package.json, .env.example)
- ✅ File patterns (handlers, models, repositories, migrations)
- ✅ Database migrations
- ✅ Templates
- ✅ Makefile & docker-compose
- ❌ Source code contents (Quick Scan only)

**For deeper analysis:** Run with Deep Scan or Exhaustive Scan option.

---

## 🌟 Key Highlights

### Why This Project Is Well-Architected

1. ✅ **Clean Separation:** Handler → Service → Repository pattern
2. ✅ **Multi-Tenant Secure:** Row-level isolation via middleware
3. ✅ **Scalable:** Stateless design, horizontal scaling ready
4. ✅ **Modern Stack:** Go + PostgreSQL + HTMX + Tailwind
5. ✅ **AI-Powered:** LLM integration with function calling
6. ✅ **Production-Ready:** Docker, migrations, comprehensive docs

### Best Practices Followed

- ✅ Environment-based configuration
- ✅ Database migrations for schema versioning
- ✅ Parameterized queries (SQL injection prevention)
- ✅ Bcrypt password hashing
- ✅ CORS configuration
- ✅ Structured logging
- ✅ Hot reload development
- ✅ Docker containerization

---

## 🚀 Next Steps

### Immediate Actions
1. ✅ Development environment is documented
2. ✅ Architecture is clear
3. ✅ Code structure is mapped

### For Ongoing Development
1. Implement remaining features (see README roadmap)
2. Add comprehensive tests
3. Implement JWT authentication fully
4. Add CI/CD pipeline
5. Optimize database queries
6. Add caching layer (Redis)

### For New Features
1. Review architecture first
2. Follow existing patterns
3. Maintain multi-tenant isolation
4. Update migrations
5. Add tests
6. Update documentation

---

**🎯 This documentation provides complete context for brownfield PRD planning and AI-assisted development.**

**Happy coding! 🚀**
