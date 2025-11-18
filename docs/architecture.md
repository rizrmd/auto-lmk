# Architecture - Auto LMK Admin Tenant Management

> **Generated:** 2025-11-15
> **Project:** Auto LMK - Multi-Tenant Car Sales Platform with AI WhatsApp Bot
> **Version:** 2.0 (Brownfield Enhancement)
> **Author:** Yopi
> **Type:** Backend Monolith with Server-Rendered Frontend
> **Enhancement Focus:** Admin Interface + Sales Car Upload via WhatsApp

---

## Executive Summary

Dokumen ini mendefinisikan arsitektur teknis untuk **brownfield enhancement** Auto LMK yang menambahkan:
1. **Admin tenant management** interface untuk mengelola sales team dan monitor conversations
2. **Sales car upload via WhatsApp** - fitur novel yang memungkinkan sales upload mobil dari lapangan via chat AI

**Konteks Brownfield:**
- ✅ Stack teknologi sudah established (Go + Chi + PostgreSQL + HTMX + Tailwind v4)
- ✅ WhatsApp bot AI sudah fully implemented
- ✅ LLM integration dengan function calling sudah berjalan
- ✅ Multi-tenant architecture sudah solid
- 🆕 Perlu menambahkan handlers, templates, dan fitur novel car upload

**Keputusan Arsitektur Utama:**
- **Mobile-First Responsive Design** - Nyaman di smartphone dan laptop
- **Local File Storage** - Foto mobil disimpan di filesystem lokal
- **Server-Side Validation** - Security first, client optional
- **Inline Error Display** - User-friendly error near form fields
- **WebSocket for Pairing Status** - Real-time WhatsApp pairing updates (bukan polling)
- **Alpine.js State Management** - Best practices untuk UI state
- **Novel Pattern: AI-Powered Conversational Upload** - Sales upload mobil via WhatsApp chat

---

## Project Initialization

**Tidak ada project initialization command** - ini brownfield enhancement pada codebase existing.

**Existing Setup:**
```bash
# Project sudah initialized
# Development setup:
docker-compose up -d postgres
make migrate-up
make dev

# Server running at http://localhost:8080
```

---

## Decision Summary

| Category | Decision | Version | Affects Epics | Rationale |
| -------- | -------- | ------- | ------------- | --------- |
| **Language** | Go | 1.25.3 | All | Existing (production-ready) |
| **HTTP Router** | Chi | v5.2.3 | All | Existing (lightweight, composable) |
| **Database** | PostgreSQL | 15-alpine | All | Existing (ACID, multi-tenant ready) |
| **Frontend Framework** | HTMX | Latest | All | Existing (server-rendered, fast) |
| **CSS Framework** | Tailwind CSS | 4.1.17 | All | Existing (utility-first, responsive) |
| **JavaScript** | Alpine.js | Latest | Epic 3, 5 | Existing (minimal interactivity) |
| **WhatsApp Library** | Whatsmeow | 2024-11-10 | Epic 3, 6 | Existing (reliable WA Web API) |
| **LLM Provider** | Z.AI (glm-4.6) | Latest | Epic 6 | Existing (function calling support) |
| **Architecture Pattern** | Clean Architecture | - | All | Existing (4-layer: Handler → Service → Repository → DB) |
| **Multi-Tenant Strategy** | Domain-based middleware | - | All | Existing (context-based tenant ID) |
| **Responsive Design** | Mobile-First | - | All | **NEW - Nyaman di smartphone & laptop** |
| **Image Storage** | Local Filesystem | - | Epic 6 | **NEW - Simple, no cloud cost** |
| **Form Validation** | Server-Side (primary) | - | Epic 1, 2, 5 | **NEW - Security first** |
| **Error Display** | Inline (near fields) | - | Epic 5 | **NEW - User-friendly** |
| **Pairing Status Update** | WebSocket | - | Epic 3 | **NEW - Real-time, efficient** |
| **State Management** | Alpine.js best practices | - | Epic 3, 4, 5 | **NEW - Reactive x-data components** |
| **WhatsApp Image Flow** | Whatsmeow → Local upload | - | Epic 6 | **NEW - Reliable, controlled** |
| **Logging** | Zerolog (structured) | v1.34.0 | All | Existing (production-grade) |
| **Hot Reload** | Air | Latest | Dev only | Existing (fast development) |
| **Containerization** | Docker Compose | Latest | All | Existing (dev & production) |

---

## Project Structure

```
auto-lmk/
├── cmd/
│   └── api/
│       └── main.go                         [MODIFY] Uncomment sales & conversation routes
│
├── internal/
│   ├── handler/
│   │   ├── car_handler.go                 [EXISTS] ✓ Car CRUD (reference pattern)
│   │   ├── whatsapp_handler.go            [EXISTS] ✓ WA pairing/status/test
│   │   ├── page_handler.go                [EXISTS] ✓ Frontend template rendering
│   │   ├── sales_handler.go               [CREATE] Epic 1 - Sales team CRUD
│   │   └── conversation_handler.go        [CREATE] Epic 2 - List & detail conversations
│   │
│   ├── service/
│   │   ├── whatsapp_service.go            [MODIFY] Epic 6 - Add photo upload handling
│   │   └── car_service.go                 [CREATE] Epic 6 - Optional business logic layer
│   │
│   ├── repository/
│   │   ├── sales_repository.go            [EXISTS] ✓ Sales CRUD methods
│   │   ├── conversation_repository.go     [EXISTS] ✓ Conversation & message queries
│   │   ├── car_repository.go              [MODIFY] Epic 6 - Add bulk photo save method
│   │   └── tenant_repository.go           [EXISTS] ✓ Tenant lookup
│   │
│   ├── llm/
│   │   ├── bot.go                         [MODIFY] Epic 6 - Add uploadCar function & role-based prompts
│   │   └── adapter.go                     [EXISTS] ✓ Z.AI client
│   │
│   ├── middleware/
│   │   ├── tenant.go                      [EXISTS] ✓ Domain-based tenant extraction
│   │   ├── cors.go                        [EXISTS] ✓ CORS handling
│   │   └── websocket.go                   [CREATE] Epic 3 - WebSocket pairing status broadcast
│   │
│   ├── model/
│   │   └── *.go                           [EXISTS] ✓ Domain models
│   │
│   └── util/
│       └── file.go                        [CREATE] Epic 6 - Photo upload helpers
│
├── templates/
│   ├── admin/
│   │   ├── layout.html                    [MODIFY] Epic 4 - Add nav menu items
│   │   ├── dashboard.html                 [MODIFY] Epic 4 - Add quick stats widgets
│   │   ├── whatsapp.html                  [MODIFY] Epic 3 - WebSocket pairing UI
│   │   ├── sales.html                     [CREATE] Epic 1 - Sales management page
│   │   └── conversations.html             [CREATE] Epic 2 - Conversation list & detail
│   │
│   ├── components/
│   │   ├── toast.html                     [CREATE] Epic 5 - Success notification
│   │   ├── modal.html                     [CREATE] Epic 5 - Confirmation dialog
│   │   └── error-inline.html              [CREATE] Epic 5 - Validation errors
│   │
│   └── layouts/
│       └── base.html                      [EXISTS] ✓ Base layout
│
├── static/
│   ├── uploads/
│   │   └── cars/
│   │       └── {tenant_id}/               [CREATE] Epic 6 - Tenant-isolated storage
│   │           ├── temp/                  Pending photos (10 min expiry)
│   │           └── {car_id}/              Permanent car photos
│   │
│   ├── css/
│   │   └── output.css                     [EXISTS] ✓ Tailwind compiled output
│   │
│   └── js/
│       ├── alpine-state.js                [CREATE] Epic 5 - Alpine.js state management
│       └── websocket-client.js            [CREATE] Epic 3 - WA pairing WebSocket client
│
├── migrations/
│   └── 000001-000009_*.sql                [EXISTS] ✓ Database schema
│
├── pkg/
│   └── security/
│       └── jwt.go                         [EXISTS] ✓ JWT helpers (future use)
│
├── docker-compose.yml                     [EXISTS] ✓
├── Dockerfile                             [EXISTS] ✓
├── Makefile                               [EXISTS] ✓
├── .air.toml                              [EXISTS] ✓ Hot reload config
├── tailwind.config.js                     [MODIFY] Epic 4 - Responsive breakpoints
├── vite.config.js                         [EXISTS] ✓
├── go.mod                                 [EXISTS] ✓
├── go.sum                                 [EXISTS] ✓
└── README.md                              [EXISTS] ✓
```

---

## Epic to Architecture Mapping

| Epic | Components | Files Created/Modified | Layer |
|------|-----------|------------------------|-------|
| **Epic 1: Sales Team Management** | Sales Handler, Sales Template | `sales_handler.go`, `sales.html`, `main.go` (routes) | Presentation |
| **Epic 2: Conversation Monitoring** | Conversation Handler, Conversation Template | `conversation_handler.go`, `conversations.html`, `main.go` (routes) | Presentation |
| **Epic 3: WhatsApp Management Enhancement** | WebSocket Middleware, Enhanced Template, WS Client | `websocket.go`, `whatsapp.html`, `websocket-client.js` | Middleware + Presentation |
| **Epic 4: Admin Dashboard Integration** | Dashboard & Layout Updates | `dashboard.html`, `layout.html` | Presentation |
| **Epic 5: Error Handling & Validation** | Component Templates, Alpine State | `toast.html`, `modal.html`, `error-inline.html`, `alpine-state.js` | Presentation |
| **Epic 6: Sales Car Upload via WhatsApp** | Bot Enhancement, Service Layer, Photo Storage | `bot.go`, `whatsapp_service.go`, `car_repository.go`, `static/uploads/` | Business Logic + Data |

---

## Technology Stack Details

### Core Technologies

**Backend:**
- **Language:** Go 1.25.3
- **HTTP Router:** Chi v5.2.3 (lightweight, composable middleware)
- **Database:** PostgreSQL 15-alpine (ACID compliance, row-level multi-tenancy)
- **DB Driver:** lib/pq v1.10.9 (pure Go PostgreSQL driver)
- **Migrations:** golang-migrate/migrate (schema versioning)

**Frontend:**
- **Template Engine:** Go html/template (stdlib, server-rendered)
- **HTMX:** Latest (dynamic partial updates without SPA complexity)
- **CSS:** Tailwind CSS v4.1.17 (utility-first, mobile-first responsive)
- **Build Tool:** Vite 7.2.2 (fast CSS compilation)
- **JavaScript:** Alpine.js (minimal reactive state management)

**Integrations:**
- **WhatsApp:** Whatsmeow (WhatsApp Web API, multi-device support)
- **LLM:** Z.AI glm-4.6 (OpenAI-compatible, function calling, Bahasa Indonesia)

**DevOps:**
- **Container:** Docker + Docker Compose
- **Hot Reload:** Air (development auto-reload)
- **Logging:** Zerolog v1.34.0 (structured JSON logging)
- **Security:** golang.org/x/crypto v0.44.0 (Bcrypt, JWT)

### Integration Points

**1. Handler → Service → Repository Flow:**
```
HTTP Request
    → Chi Router (with tenant middleware)
        → Handler (extract tenant, validate input)
            → Service (business logic, orchestration)
                → Repository (SQL queries with tenant filter)
                    → PostgreSQL (data persistence)
```

**2. WhatsApp → LLM → Database Flow:**
```
WhatsApp Message Received (Whatsmeow)
    → WhatsAppService.ProcessIncomingMessage()
        → SalesRepository.IsSales() [check role]
            → Bot.ProcessMessage() [LLM with role-based prompt]
                → LLM Function Calling (searchCars / uploadCar / etc)
                    → Repository Methods (car search / create)
                        → Database Transaction
                            → WhatsApp Response
```

**3. WebSocket Pairing Status Flow:**
```
Admin Browser
    → WebSocket Connect (/ws/pairing/{tenant_id})
        → WebSocket Middleware (authenticate, subscribe)
            → Whatsmeow Pairing Event
                → Broadcast to connected clients
                    → UI Updates Real-Time
```

**4. Photo Upload Flow (Epic 6):**
```
Sales WhatsApp (Send Photo)
    → Whatsmeow Media Download
        → Save to: static/uploads/cars/{tenant_id}/temp/{timestamp}.jpg
            → Bot.AddPendingPhoto(phone, path) [in-memory map, 10min timer]

Sales WhatsApp (Send Text: "Toyota Avanza 2020 185juta AT Bensin")
    → Bot.ProcessMessage() [AI parses text]
        → LLM Function Call: uploadCar(brand, model, year, price, transmission, fuel)
            → Bot.GetPendingPhotos(phone)
                → CarRepository.Create() [get carID]
                    → Move photos: temp/ → {car_id}/
                        → CarRepository.AddPhotos(carID, photoPaths)
                            → Bot.ClearPendingPhotos(phone)
                                → Response: "✓ Mobil berhasil ditambahkan! Link: https://tenant.com/cars/123"
```

---

## Novel Pattern Design: AI-Powered Conversational Upload

### Pattern Name
**Conversational Multi-Step Upload with Temporal Context**

### Problem Solved
Bagaimana cara user upload data kompleks (mobil + multiple photos + specifications) via chat interface tanpa form tradisional?

### Core Innovation
1. **In-Memory Pending Context** - Associate uploaded photos dengan user session selama 10 menit
2. **LLM Function Calling** - Parse free-form text (e.g., "Toyota Avanza 2020 185juta AT Bensin") jadi structured data
3. **Role-Based Access** - Sales dapat uploadCar, customer tidak bisa (polite rejection)
4. **Auto-Expiration** - Pending photos auto-delete setelah 10 menit atau setelah success

### Component Architecture

```
┌─────────────────────────────────────────────────────────────┐
│               AI-Powered Upload Pattern                      │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐   ┌──────────────┐   ┌──────────────┐   │
│  │  WhatsApp    │   │   Pending    │   │     LLM      │   │
│  │   Message    │──→│   Photo      │←──│    Bot       │   │
│  │   Handler    │   │   Context    │   │ (Role-Aware) │   │
│  └──────┬───────┘   │   Manager    │   └──────┬───────┘   │
│         │            │              │          │            │
│         │            │ (In-Memory   │          │            │
│         │            │  Map with    │          │            │
│         │            │ Expiration)  │          │            │
│         │            └──────────────┘          │            │
│         │                                      │            │
│         ▼                                      ▼            │
│  ┌──────────────┐                      ┌──────────────┐   │
│  │    File      │                      │  Function    │   │
│  │   Storage    │                      │  Executor:   │   │
│  │   Manager    │                      │  uploadCar() │   │
│  └──────┬───────┘                      └──────┬───────┘   │
│         │                                      │            │
│         └──────────────┬───────────────────────┘           │
│                        ▼                                    │
│                 ┌──────────────┐                           │
│                 │  Database    │                           │
│                 │ Transaction  │                           │
│                 │  (Car +      │                           │
│                 │   Photos)    │                           │
│                 └──────────────┘                           │
└─────────────────────────────────────────────────────────────┘
```

### Data Structures

**In-Memory Pending Photo Context:**
```go
type PendingPhotoContext struct {
    mu            sync.RWMutex
    photosByPhone map[string]*PhotoSession
}

type PhotoSession struct {
    TenantID   int
    Photos     []string  // File paths
    UploadedAt time.Time
    ExpiresAt  time.Time
}

// Methods:
func (c *PendingPhotoContext) Add(phone string, tenantID int, photoPath string) int
func (c *PendingPhotoContext) Get(phone string) []string
func (c *PendingPhotoContext) Clear(phone string)
func (c *PendingPhotoContext) StartExpirationTimer(phone string, duration time.Duration)
```

**LLM Function Definition:**
```json
{
  "name": "uploadCar",
  "description": "Upload mobil baru ke catalog (hanya untuk sales team)",
  "parameters": {
    "type": "object",
    "properties": {
      "brand": {"type": "string", "description": "Merek mobil (Toyota, Honda, dll)"},
      "model": {"type": "string", "description": "Model mobil (Avanza, Brio, dll)"},
      "year": {"type": "integer", "description": "Tahun produksi"},
      "price": {"type": "integer", "description": "Harga dalam Rupiah"},
      "transmission": {"type": "string", "enum": ["AT", "MT"]},
      "fuel_type": {"type": "string", "enum": ["Bensin", "Diesel", "Hybrid", "Electric"]}
    },
    "required": ["brand", "model", "year", "price", "transmission", "fuel_type"]
  }
}
```

### Sequence Diagrams

**Sequence 1: Sales Upload Photos**
```
Sales → WhatsApp: Send Photo 1
WhatsApp → Service: ProcessMessage(image, mediaURL)
Service → Whatsmeow: DownloadMedia(mediaURL)
Service → FileStorage: Save to temp/{timestamp}.jpg
Service → PendingContext: Add(phone, filePath)
PendingContext → Timer: Start 10min expiration
Service → Sales: "Foto 1 diterima! Upload lebih atau ketik detail."

Sales → WhatsApp: Send Photo 2
[Repeat process]
Service → Sales: "Foto 2 diterima! Total 2 foto."
```

**Sequence 2: Sales Submit Details**
```
Sales → WhatsApp: "Toyota Avanza 2020 185juta AT Bensin"
WhatsApp → Service: ProcessMessage(text)
Service → SalesRepo: IsSales(phone) → true
Service → Bot: ProcessMessage(text, isSales=true)
Bot → LLM: API Call with uploadCar function
LLM → Bot: Function Call uploadCar(brand="Toyota", model="Avanza", ...)
Bot → PendingContext: Get(phone) → [photo1.jpg, photo2.jpg]
Bot → CarRepo: Create(tenantID, car) → carID=123
Bot → FileStorage: Move temp/* → {carID}/
Bot → CarRepo: AddPhotos(carID, photos)
Bot → PendingContext: Clear(phone)
Bot → Sales: "✓ Mobil berhasil ditambahkan! Link: https://tenant.com/cars/123"
```

**Sequence 3: Customer Rejection**
```
Customer → WhatsApp: Send Photo
WhatsApp → Service: ProcessMessage(image, mediaURL)
Service → SalesRepo: IsSales(phone) → false
Service → Customer: "Terima kasih! Untuk upload mobil, hubungi sales team kami."
[Photo not saved, immediate rejection]
```

### Edge Cases & Handling

| Edge Case | Handling Strategy |
|-----------|------------------|
| **Sales upload foto, tidak kirim detail** | Auto-expire 10 menit, hapus temp files, user upload ulang |
| **Sales kirim detail tanpa foto** | Bot: "Silakan upload foto terlebih dahulu" |
| **LLM gagal parse text** | Bot: "Format tidak jelas. Contoh: Toyota Avanza 2020 185juta AT Bensin" |
| **Database save gagal** | Rollback, hapus temp files, error message |
| **Foto corrupt/invalid** | Validate saat download, reject dengan error |
| **Customer coba upload** | Immediate polite rejection |
| **Multiple sales upload simultan** | Each phone = separate context (map key = phone) |
| **Server restart saat pending** | Lost (acceptable, max 10 min loss, user upload ulang) |

### Pattern Benefits
- ✅ **Natural UX:** Sales tidak perlu buka admin panel, langsung dari WA
- ✅ **Time Saving:** Upload dari lapangan, tidak perlu kembali ke kantor
- ✅ **AI-Powered:** Free-form text parsing, tidak butuh format strict
- ✅ **Secure:** Role-based, customer tidak bisa upload
- ✅ **Simple:** In-memory context, no database overhead
- ✅ **Reliable:** Auto-expiration, no stale data

---

## Implementation Patterns (Consistency Rules)

### 1. NAMING CONVENTIONS

**Go Code:**
```go
// ✅ CORRECT
type SalesHandler struct { }          // PascalCase for types
func (h *SalesHandler) Create()       // PascalCase for exported
func formatPhoneNumber()              // camelCase for private
var tenantID int                      // camelCase for variables

// ❌ WRONG
type sales_handler struct { }         // Don't snake_case
func (h *SalesHandler) create()       // Don't lowercase exported
```

**Database:**
```sql
-- ✅ CORRECT
CREATE TABLE sales (              -- Lowercase, plural
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER,            -- snake_case columns
    phone_number VARCHAR(20),
    created_at TIMESTAMP
);

-- ❌ WRONG
CREATE TABLE Sales                -- Don't capitalize
CREATE TABLE sale                 -- Don't singular
```

**Files:**
```
✅ CORRECT:
internal/handler/sales_handler.go       -- snake_case
templates/admin/conversations.html      -- lowercase, plural

❌ WRONG:
internal/handler/SalesHandler.go        -- Don't PascalCase
templates/admin/Conversation.html       -- Don't singular/capitalize
```

**API Endpoints:**
```
✅ CORRECT:
POST /api/sales                    -- Lowercase, plural
GET /api/conversations/{id}        -- {id} placeholder

❌ WRONG:
POST /api/Sale                     -- Don't capitalize
GET /api/conversation/123          -- Don't literal ID
```

**HTML IDs & Classes:**
```html
✅ CORRECT:
<div id="sales-table">             <!-- kebab-case -->
<button class="btn-primary">       <!-- kebab-case custom -->
<div class="bg-blue-500">          <!-- Tailwind as-is -->

❌ WRONG:
<div id="salesTable">              <!-- Don't camelCase -->
```

---

### 2. STRUCTURE PATTERNS

**Handler Structure (MANDATORY for ALL handlers):**
```go
type XxxHandler struct {
    repo *repository.XxxRepository
}

func NewXxxHandler(repo *repository.XxxRepository) *XxxHandler {
    return &XxxHandler{repo: repo}
}

func (h *XxxHandler) Create(w http.ResponseWriter, r *http.Request) {
    // 1. CRITICAL: Extract tenant
    tenantID, err := model.GetTenantID(r.Context())
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // 2. Parse request
    var req CreateXxxRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // 3. Validate
    if err := req.Validate(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // 4. Call repository (with tenant filter)
    result, err := h.repo.Create(r.Context(), tenantID, &req)
    if err != nil {
        log.Error().Err(err).Int("tenant_id", tenantID).Msg("Create failed")
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }

    // 5. Return JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(result)
}
```

**Template Structure:**
```html
{{define "admin/sales"}}
<!DOCTYPE html>
<html>
{{template "admin/header" .}}
<body class="bg-gray-50">
    {{template "admin/navbar" .}}

    <main class="container mx-auto px-4 py-8">
        <!-- Page Title -->
        <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
            <h1 class="text-2xl font-bold">Sales Team</h1>
            <button class="btn-primary w-full md:w-auto">Tambah Sales</button>
        </div>

        <!-- Content -->
        <div id="sales-table" class="overflow-x-auto">
            <!-- HTMX swaps content here -->
        </div>
    </main>

    {{template "admin/footer" .}}
</body>
</html>
{{end}}
```

**Repository Pattern:**
```go
func (r *XxxRepository) Create(ctx context.Context, tenantID int, data *Xxx) (int, error) {
    // CRITICAL: ALWAYS filter by tenant_id
    query := `INSERT INTO xxx (tenant_id, field1, field2)
              VALUES ($1, $2, $3) RETURNING id`

    var id int
    err := r.db.QueryRowContext(ctx, query, tenantID, data.Field1, data.Field2).Scan(&id)
    return id, err
}

func (r *XxxRepository) List(ctx context.Context, tenantID int, filters map[string]interface{}) ([]Xxx, error) {
    // CRITICAL: ALWAYS filter by tenant_id
    query := `SELECT id, field1, field2 FROM xxx WHERE tenant_id = $1`
    // ... apply additional filters
}
```

---

### 3. FORMAT PATTERNS

**API Response Format:**
```json
✅ SUCCESS:
{
  "data": { ... },
  "message": "Success"
}

✅ ERROR:
{
  "error": "Nomor sudah terdaftar",
  "field": "phone_number"
}

✅ LIST (Paginated):
{
  "data": [...],
  "total": 45,
  "page": 1,
  "limit": 20,
  "total_pages": 3
}

❌ WRONG:
{
  "success": true,              // Don't use boolean
  "result": { ... }             // Use "data"
}
```

**Date/Time Format:**
```
API (JSON): "2025-11-15T14:30:00Z"        // ISO 8601 UTC
UI Display: "15 Nov 2025, 14:30"          // Formatted, WIB
UI Relative: "5 menit lalu"               // Relative time
```

**Error Messages (Bahasa Indonesia):**
```go
✅ CORRECT:
"Nomor sudah terdaftar"
"Format nomor tidak valid"
"Data tidak ditemukan"
"Terjadi kesalahan, silakan coba lagi"

❌ WRONG:
"Duplicate entry"
"Invalid phone format"
"Record not found"
```

---

### 4. COMMUNICATION PATTERNS

**HTMX Attributes (CONSISTENT):**
```html
✅ CORRECT:
<form hx-post="/api/sales"
      hx-target="#sales-table"
      hx-swap="outerHTML">

<button hx-delete="/api/sales/{id}"
        hx-confirm="Yakin ingin menghapus?"
        hx-target="closest tr"
        hx-swap="outerHTML swap:1s">

<div hx-get="/api/conversations"
     hx-trigger="load"
     hx-indicator="#spinner">

❌ WRONG:
<form hx-post="/api/sales" hx-target="body">     <!-- Too broad -->
<button hx-delete="/api/sales/{id}">             <!-- No confirmation -->
```

**WebSocket Message Format:**
```json
// Server → Client
{
  "type": "pairing_status",
  "data": {
    "status": "connected",
    "phone_number": "+628123456789",
    "timestamp": "2025-11-15T14:30:00Z"
  }
}

// Client → Server
{
  "action": "subscribe",
  "tenant_id": 1
}
```

---

### 5. LIFECYCLE PATTERNS

**Loading States:**
```html
<button hx-post="/api/sales"
        hx-indicator="#spinner">
    Simpan
</button>
<div id="spinner" class="htmx-indicator">
    <svg class="animate-spin h-5 w-5" ...>...</svg>
    Menyimpan...
</div>
```

**Error Handling:**
```html
<div hx-get="/api/sales"
     hx-trigger="load"
     hx-on::error="alert('Gagal memuat data. Silakan refresh halaman.')">
</div>
```

**Success Notification (Alpine.js):**
```html
<div x-data="{ show: false }"
     @success.window="show = true; setTimeout(() => show = false, 3000)"
     x-show="show"
     x-transition
     class="fixed top-4 right-4 bg-green-500 text-white px-4 py-3 rounded-lg shadow-lg">
    ✓ Berhasil disimpan!
</div>
```

---

### 6. LOCATION PATTERNS

**File Upload Paths:**
```
Photo storage:
static/uploads/cars/{tenant_id}/{car_id}/photo1.jpg
static/uploads/cars/{tenant_id}/{car_id}/photo2.jpg

Temp storage (pending):
static/uploads/cars/{tenant_id}/temp/{timestamp}_{random}.jpg

✅ ALWAYS include tenant_id in path (security)
❌ NEVER store cross-tenant files together
```

**Static Assets:**
```
CSS: static/css/output.css           (Tailwind compiled)
JS:  static/js/alpine-state.js       (Alpine components)
     static/js/websocket-client.js   (WebSocket client)

Production: Use versioned URLs
Example: /static/css/output.css?v=20251115
```

---

### 7. CONSISTENCY PATTERNS (CRITICAL)

**Tenant Isolation (MANDATORY in ALL handlers):**
```go
// ✅ CORRECT - ALWAYS first step
tenantID, err := model.GetTenantID(r.Context())
if err != nil {
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
}

// Use tenantID in ALL database operations
result, err := h.repo.GetByID(ctx, id, tenantID)

// ❌ WRONG - SECURITY BREACH
result, err := h.repo.GetByID(ctx, id)  // Missing tenant filter!
```

**Logging Pattern:**
```go
log.Info().
    Int("tenant_id", tenantID).
    Str("action", "create_sales").
    Str("phone", phoneNumber).
    Msg("Sales member created successfully")

log.Error().
    Err(err).
    Int("tenant_id", tenantID).
    Str("operation", "delete_sales").
    Int("sales_id", salesID).
    Msg("Failed to delete sales member")
```

**Validation Pattern:**
```go
type CreateSalesRequest struct {
    PhoneNumber string `json:"phone_number"`
    Name        string `json:"name"`
    Role        string `json:"role"`
}

func (r *CreateSalesRequest) Validate() error {
    // Phone validation
    phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$|^08\d{8,11}$`)
    if !phoneRegex.MatchString(r.PhoneNumber) {
        return errors.New("Format nomor tidak valid")
    }

    // Name validation
    if strings.TrimSpace(r.Name) == "" {
        return errors.New("Nama tidak boleh kosong")
    }

    return nil
}
```

**Responsive Design (Mobile-First):**
```html
<!-- Stack on mobile, grid on desktop -->
<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">

<!-- Hide on mobile, show on desktop -->
<div class="hidden md:block">Table view</div>
<div class="block md:hidden">Card view</div>

<!-- Full width mobile, constrained desktop -->
<div class="w-full md:w-1/2 lg:w-1/3">

<!-- Responsive text size -->
<h1 class="text-xl md:text-2xl lg:text-3xl">

<!-- Touch-friendly buttons (min 44x44px) -->
<button class="px-6 py-3 text-base">Simpan</button>
```

---

## Data Architecture

### Database Schema

**Existing Tables (No Changes):**
- `tenants` - Root table, no tenant_id
- `users` - Admin users (tenant-scoped)
- `sales` - Sales team members (tenant-scoped)
- `cars` - Vehicle inventory (tenant-scoped)
- `car_specs` - Additional specs EAV (via cars)
- `car_photos` - Car images (via cars)
- `conversations` - WhatsApp conversations (tenant-scoped)
- `messages` - Chat messages (via conversations)

**Key Indexes (Existing):**
```sql
CREATE INDEX idx_sales_tenant_id ON sales(tenant_id);
CREATE INDEX idx_sales_phone ON sales(tenant_id, phone_number);
CREATE INDEX idx_conversations_tenant ON conversations(tenant_id);
CREATE INDEX idx_conversations_created ON conversations(tenant_id, created_at DESC);
CREATE INDEX idx_messages_conversation ON messages(conversation_id, created_at ASC);
```

**No Schema Changes Required** - Existing schema already supports all epics.

---

## API Contracts

### Sales Team Management

```
POST /api/sales
Body: {"phone_number": "+628123456789", "name": "John Doe", "role": "Sales Executive"}
Response 201: {"id": 1, "tenant_id": 1, "phone_number": "+628123456789", "name": "John Doe", ...}

GET /api/sales
Response 200: {"data": [...], "total": 5}

DELETE /api/sales/{id}
Response 204: No Content
```

### Conversation Monitoring

```
GET /api/conversations?page=1&limit=20&type=customer|sales|all
Response 200: {"data": [...], "total": 45, "page": 1, "limit": 20, "total_pages": 3}

GET /api/conversations/{id}
Response 200: {"conversation": {...}, "messages": [...], "total_messages": 15}
```

### WhatsApp Management

```
GET /admin/whatsapp/status
Response 200: {"is_connected": true, "phone_number": "+628123456789", "pairing_status": "connected"}

POST /admin/whatsapp/pair
Body: {"phone_number": "+628123456789"}
Response 200: {"qr_code_url": "/admin/whatsapp/qr/{tenant_id}"}

WebSocket: /ws/pairing/{tenant_id}
Message: {"type": "pairing_status", "data": {"status": "connected", ...}}
```

---

## Security Architecture

### Multi-Tenant Isolation (CRITICAL)

**Enforcement Points:**
1. **Middleware:** `internal/middleware/tenant.go` extracts tenant from domain
2. **Context:** Tenant ID injected into `r.Context()`
3. **Handler:** MUST call `model.GetTenantID(r.Context())` first
4. **Repository:** MUST filter by `tenant_id` in ALL queries
5. **File Storage:** MUST use `{tenant_id}` in path

**Validation:**
```go
// MANDATORY in every handler
tenantID, err := model.GetTenantID(r.Context())
if err != nil {
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
}
```

### Input Validation

- **Server-Side (Primary):** All handlers validate using `Validate()` methods
- **Client-Side (Optional):** HTML5 `required`, `pattern` for UX only
- **SQL Injection Prevention:** Parameterized queries ONLY (no string concat)
- **XSS Prevention:** Go templates auto-escape HTML

### Authentication (MVP: Basic)

- Session-based auth for admin panel
- HTTPS required in production (Nginx SSL)
- Future: JWT authentication (framework already exists)

### Rate Limiting

- API: 100 req/min per tenant
- WhatsApp pairing: 5 attempts/hour per tenant
- Test messages: 20/hour per tenant

---

## Performance Considerations

### Response Time Targets

- API endpoints: < 200ms (p95)
- Page load: < 2s (full render)
- HTMX partial updates: < 500ms
- WebSocket latency: < 100ms

### Database Optimization

- All queries use indexes (tenant_id, created_at, phone_number)
- Pagination for large datasets (conversations, messages)
- Limit query results (max 100 messages per request)

### File Storage

- Local filesystem (fast read/write)
- Tenant isolation via directory structure
- Auto-cleanup temp files (10 min expiration)

### Caching Strategy (Future)

- Static assets: Browser cache (1 year)
- API responses: Consider Redis for read-heavy endpoints
- Photo thumbnails: CDN (future enhancement)

---

## Deployment Architecture

### Development

```
Docker Compose:
- PostgreSQL container
- Go app (hot reload with Air)
- Vite dev server (CSS watch)

Access: http://localhost:8080
```

### Production

```
Nginx (SSL termination, reverse proxy)
    → Go App (systemd service)
        → PostgreSQL (managed instance)

Multi-tenant routing:
tenant1.auto-lmk.com → Nginx → Go App (extracts tenant from domain)
tenant2.auto-lmk.com → Nginx → Go App (extracts tenant from domain)
```

### Environment Variables

```bash
# Database
DATABASE_URL=postgres://user:pass@localhost:5432/autolmk

# Server
PORT=8080
ENV=production

# LLM
ZAI_API_KEY=your_api_key
ZAI_API_URL=https://api.z.ai/v1

# WhatsApp
WHATSAPP_DATA_DIR=/var/lib/autolmk/whatsapp

# Storage
UPLOAD_DIR=/var/lib/autolmk/uploads
```

---

## Development Environment

### Prerequisites

- Go 1.25.3+
- PostgreSQL 15+
- Docker & Docker Compose
- Node.js 18+ (for Tailwind CSS build)

### Setup Commands

```bash
# Clone project
git clone https://github.com/riz/auto-lmk.git
cd auto-lmk

# Environment setup
cp .env.example .env
# Edit .env with your credentials

# Start database
docker-compose up -d postgres

# Run migrations
make migrate-up

# Install frontend dependencies
npm install

# Start development server (with hot reload)
make dev

# In another terminal: Watch CSS
npm run dev

# Access: http://localhost:8080
```

---

## Architecture Decision Records (ADRs)

### ADR-001: Mobile-First Responsive Design

**Decision:** Implement mobile-first responsive design untuk semua admin interfaces.

**Rationale:**
- Admin sering bekerja di lapangan (butuh akses via smartphone)
- Sales team lebih banyak pakai mobile daripada laptop
- Tailwind CSS v4 sudah support responsive utilities

**Consequences:**
- Design dimulai dari mobile breakpoint, kemudian scale up
- Testing harus dilakukan di smartphone dan desktop
- Component library harus responsive-aware

**Alternatives Considered:**
- Desktop-first: Rejected (tidak sesuai user behavior)
- Mobile-only: Rejected (desktop tetap digunakan di kantor)

---

### ADR-002: WebSocket untuk WhatsApp Pairing Status

**Decision:** Gunakan WebSocket untuk real-time pairing status updates (bukan polling).

**Rationale:**
- Real-time experience lebih baik untuk pairing flow
- Lebih efisien daripada polling setiap 3 detik
- Go sudah punya gorilla/websocket library yang mature

**Consequences:**
- Perlu implement WebSocket middleware
- Perlu JavaScript WebSocket client
- Connection management (reconnect logic)

**Alternatives Considered:**
- HTMX polling: Rejected (overhead 3 detik polling)
- Server-Sent Events (SSE): Rejected (WebSocket lebih flexible)

---

### ADR-003: Local File Storage untuk Car Photos

**Decision:** Simpan foto mobil di local filesystem, bukan cloud storage.

**Rationale:**
- Simple, no external dependencies
- No cloud storage costs (MVP)
- Cukup untuk 10-100 tenants dengan 50-500 cars each
- Easy backup (rsync / tar)

**Consequences:**
- File path structure: `static/uploads/cars/{tenant_id}/{car_id}/`
- Backup strategy: Include uploads directory
- Future migration: Can move to S3/GCS later

**Alternatives Considered:**
- AWS S3: Rejected (cost, complexity untuk MVP)
- Database BLOB: Rejected (poor performance)

---

### ADR-004: In-Memory Pending Photo Context

**Decision:** Gunakan in-memory map dengan auto-expiration untuk pending photos (bukan database).

**Rationale:**
- Temporary data, tidak perlu persist
- Fast access (no DB roundtrip)
- Simple implementation dengan Go timers
- Auto-cleanup dengan expiration

**Consequences:**
- Lost saat server restart (acceptable - user upload ulang)
- Max 10 menit retention
- Memory usage: ~1KB per pending session

**Alternatives Considered:**
- Database table: Rejected (overhead untuk temporary data)
- Redis: Rejected (extra dependency untuk MVP)

---

### ADR-005: Server-Side Validation Only

**Decision:** Validasi form di server-side (primary), client-side optional (UX enhancement).

**Rationale:**
- Security: Client-side validation bisa di-bypass
- Consistency: Single source of truth di server
- Simpler: No duplicate validation logic

**Consequences:**
- All validation errors require server roundtrip
- Client-side (HTML5) hanya untuk UX improvement
- HTMX handles error display seamlessly

**Alternatives Considered:**
- Client + Server: Rejected (duplicate logic, maintenance overhead untuk MVP)

---

### ADR-006: AI-Powered Conversational Upload Pattern

**Decision:** Implement novel pattern untuk car upload via WhatsApp chat dengan LLM parsing.

**Rationale:**
- Unique value proposition: Upload dari lapangan tanpa admin panel
- Natural UX: Sales familiar dengan WhatsApp
- AI strength: Parse free-form text jadi structured data
- Role-based: Sales only, customer rejected

**Consequences:**
- Requires in-memory context management
- Requires LLM function calling implementation
- Requires clear user instructions via bot prompts

**Alternatives Considered:**
- Web form upload: Rejected (sales tidak mau buka admin panel)
- WhatsApp bot dengan strict format: Rejected (poor UX)

---

## Validation Checklist

- ✅ Decision table has specific versions
- ✅ Every epic mapped to architecture components
- ✅ Source tree is complete and specific (not generic)
- ✅ No placeholder text remains
- ✅ All 24 FRs from PRD have architectural support
- ✅ All NFRs (performance, security, scalability) addressed
- ✅ Implementation patterns cover all potential agent conflicts
- ✅ Novel pattern (Epic 6) fully documented
- ✅ Mobile-first responsive design strategy defined
- ✅ WebSocket implementation specified
- ✅ File storage strategy clear
- ✅ Tenant isolation enforcement rules CRITICAL and clear

---

## Next Steps

**After Architecture Completion:**

1. **Test Design** (next required workflow)
   - Run: `/bmad:bmm:workflows:test-design`
   - Agent: TEA (Test Engineer Agent)
   - Output: Test strategy untuk 28 stories

2. **Sprint Planning** (after test design)
   - Run: `/bmad:bmm:workflows:sprint-planning`
   - Agent: SM (Scrum Master)
   - Output: Sprint breakdown dengan story assignments

3. **Implementation** (development phase)
   - Run: `/bmad:bmm:agents:dev`
   - Execute stories sprint by sprint
   - Follow architecture patterns EXACTLY

---

_Architecture document generated by BMAD Architecture Workflow v1.0_
_Date: 2025-11-15_
_For: Yopi_
_Project: auto-lmk (Brownfield Enhancement)_

**🎯 Arsitektur ini menjadi contract untuk semua AI agents. Implementasi HARUS follow patterns ini untuk consistency!**
