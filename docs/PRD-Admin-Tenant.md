# Auto LMK - Product Requirements Document
# Admin Tenant Management & WhatsApp Bot Testing Interface

**Author:** Yopi
**Date:** 2025-11-15
**Version:** 1.0
**Status:** Brownfield Enhancement

---

## Executive Summary

Platform Auto LMK saat ini telah mengimplementasikan WhatsApp bot AI-powered yang lengkap dengan integrasi LLM Z.AI (glm-4.6), function calling untuk pencarian mobil, dan multi-tenant WhatsApp client management. Namun, fungsionalitas admin tenant untuk mengelola dan menguji bot ini belum lengkap.

PRD ini fokus pada **melengkapi interface admin tenant** agar dealership dapat:
1. Melakukan pairing WhatsApp dengan mudah
2. Mengelola sales team untuk internal communication
3. Melihat dan memonitor conversation history
4. Menguji chatbot AI secara end-to-end

### What Makes This Special

Ini adalah **brownfield enhancement** yang memanfaatkan:
- ✅ WhatsApp bot yang sudah fully implemented (internal/llm/bot.go, internal/service/whatsapp_service.go)
- ✅ LLM function calling yang sudah berjalan (searchCars, getCarDetails, sendCarImages)
- ✅ Multi-tenant architecture yang sudah solid
- ✅ Database schema yang sudah lengkap (8 tables)

Yang perlu dilakukan adalah **mengaktifkan route yang sudah ada** dan **membangun UI admin** untuk testing dan management.

---

## Project Classification

**Technical Type:** Backend API + Server-Rendered Admin UI (Brownfield Enhancement)
**Domain:** Automotive SaaS / Multi-Tenant Platform
**Complexity:** Medium (leveraging existing implementation)

### Current State

**Existing Implementation (Production-Ready):**
- WhatsApp bot dengan AI (100% complete)
- Car CRUD API (100% complete)
- Multi-tenant isolation (100% complete)
- Database schema (100% complete)
- LLM integration dengan 3 functions (100% complete)

**Missing Components (Target of This PRD):**
- Sales management handlers (commented out: cmd/api/main.go:224-228)
- Conversation viewing handlers (commented out: cmd/api/main.go:231-234)
- Admin UI templates untuk WhatsApp management
- Admin UI untuk sales team management
- Admin UI untuk conversation monitoring

---

## Success Criteria

### Primary Goals

**✅ Success Achieved When:**

1. **WhatsApp Pairing Testing**
   - Admin dapat initiate pairing dari UI
   - QR code ditampilkan dengan baik
   - Status pairing real-time visible
   - Success/failure notification jelas

2. **Sales Team Management**
   - Admin dapat menambahkan sales team member
   - Admin dapat melihat daftar sales team
   - Admin dapat menghapus sales team member
   - Sales team phone numbers tersimpan di database

3. **Conversation Monitoring**
   - Admin dapat melihat list conversations
   - Admin dapat melihat detail conversation (message history)
   - Dapat membedakan customer vs sales conversations
   - Timestamp dan status message terlihat jelas

4. **End-to-End Bot Testing**
   - Admin dapat mengirim test message ke bot
   - Response bot terlihat di UI
   - Function calling bot dapat dimonitor (car search, details, images)
   - Error handling terlihat jika ada masalah

### Business Metrics

**Immediate Impact:**
- **Time to Test:** < 5 menit dari fresh install ke testing bot
- **Pairing Success Rate:** > 95% (dengan error messaging yang jelas)
- **Admin Task Completion:** Semua admin tasks dapat diselesaikan via UI (no need for database access)

**Long-term Value:**
- Dealership dapat self-serve untuk setup WhatsApp bot
- Reduce support tickets untuk pairing issues
- Increase confidence dalam AI chatbot capability

---

## Product Scope

### MVP - Minimum Viable Product

**Phase 1: Core Admin Functionality (Target: 1 Week)**

#### 1. Sales Team Management
- ✅ Repository: `internal/repository/sales_repository.go` (already exists)
- 🆕 Handler: `internal/handler/sales_handler.go` (implement)
  - `POST /api/sales` - Add sales team member
  - `GET /api/sales` - List sales team
  - `DELETE /api/sales/{id}` - Remove sales team member
- 🆕 Template: `templates/admin/sales.html`
  - List view dengan table
  - Add form (phone number, name)
  - Delete confirmation

#### 2. Conversation Monitoring
- ✅ Repository: `internal/repository/conversation_repository.go` (already exists)
- 🆕 Handler: `internal/handler/conversation_handler.go` (implement)
  - `GET /api/conversations` - List conversations (with pagination)
  - `GET /api/conversations/{id}` - Get conversation detail with messages
  - `GET /api/conversations/{id}/messages` - Get messages (for HTMX partial load)
- 🆕 Template: `templates/admin/conversations.html`
  - List view dengan filters (customer/sales, date range)
  - Detail view dengan message history
  - Visual distinction untuk inbound/outbound messages

#### 3. WhatsApp Management UI Enhancement
- ✅ Handler: `internal/handler/whatsapp_handler.go` (already complete)
- ✅ Routes: `/admin/whatsapp/*` (already exist)
- 🆕 Template: Enhance `templates/admin/whatsapp.html`
  - Better QR code display
  - Real-time status updates (HTMX polling)
  - Test message interface yang user-friendly
  - Disconnect dengan confirmation
  - Error state handling

#### 4. Admin Dashboard Integration
- 🔧 Update: `templates/admin/dashboard.html`
  - Add navigation links ke Sales dan Conversations
  - Show quick stats (total conversations, connected status)
  - Recent conversations widget
- 🔧 Update: `templates/admin/layout.html`
  - Add menu items untuk Sales dan Conversations
  - Active state untuk current page

#### 5. Route Activation
- 🔧 Update: `cmd/api/main.go`
  - Uncomment sales routes (lines 224-228)
  - Uncomment conversation routes (lines 231-234)
  - Wire up handlers correctly

### Growth Features (Post-MVP)

**Phase 2: Enhanced Monitoring & Analytics (Future)**

- 📊 Conversation analytics dashboard
  - Total conversations per day/week/month
  - Response time metrics
  - Most searched car brands/models
  - Function call statistics
- 🔍 Advanced conversation search
  - Search by phone number
  - Search by message content
  - Filter by date range, type (customer/sales)
- 📧 Notification system
  - Email alerts untuk new conversations
  - WhatsApp disconnection alerts
  - Daily summary reports
- 👤 User authentication & roles
  - JWT-based authentication (framework already in place)
  - Role-based access (admin, viewer)
  - Audit logs untuk admin actions

### Vision (Future)

**Phase 3: Advanced Features (v2.0)**

- 🤖 Bot behavior customization
  - Customize system prompt per tenant
  - Enable/disable specific functions
  - Configure response tone/style
- 📱 Mobile-optimized admin panel
- 🔗 CRM integration hooks
  - Export conversations to Salesforce/HubSpot
  - Webhook notifications untuk new leads
- 🧪 A/B testing framework
  - Test different prompts
  - Compare bot performance metrics
- 📈 Advanced AI insights
  - Customer intent analysis
  - Popular car features extraction
  - Price sensitivity detection

---

## Domain-Specific Requirements

### Automotive Sales Context

**Customer Interaction Patterns:**
- Customers biasanya bertanya dalam Bahasa Indonesia informal
- Typical queries: budget, brand preference, fuel type, transmission
- Expectations: Fast response dengan foto mobil
- Trust factor: Need to feel like talking to knowledgeable salesperson

**Dealership Operations:**
- Sales team bisa jumlahnya 5-20 orang per dealership
- Conversations bisa mencapai 100-500 per bulan per tenant
- Peak hours: evenings and weekends
- Need untuk distinguish customer leads vs internal sales communication

**Multi-Tenant Isolation:**
- CRITICAL: Setiap tenant hanya bisa lihat conversations mereka sendiri
- Phone number collision possible (customer bisa chat dengan multiple dealerships)
- WhatsApp pairing harus per-tenant (one WhatsApp number per dealership)

This section shapes all functional and non-functional requirements below.

---

## Backend-Specific Requirements

### API Specification

**All endpoints tenant-scoped via middleware:** `internal/middleware/tenant.go`

#### Sales Team Management API

```
POST /api/sales
Body: {
  "phone_number": "+628123456789",
  "name": "John Doe",
  "role": "Sales Executive" (optional)
}
Response: 201 Created
{
  "id": 1,
  "tenant_id": 1,
  "phone_number": "+628123456789",
  "name": "John Doe",
  "role": "Sales Executive",
  "created_at": "2025-11-15T10:00:00Z"
}
```

```
GET /api/sales
Response: 200 OK
{
  "sales": [
    {
      "id": 1,
      "phone_number": "+628123456789",
      "name": "John Doe",
      "role": "Sales Executive",
      "created_at": "2025-11-15T10:00:00Z"
    }
  ],
  "total": 1
}
```

```
DELETE /api/sales/{id}
Response: 204 No Content
```

#### Conversation Monitoring API

```
GET /api/conversations?page=1&limit=20&type=customer|sales&sort=recent
Response: 200 OK
{
  "conversations": [
    {
      "id": 1,
      "phone_number": "+6281234567890",
      "is_sales": false,
      "last_message": "Terima kasih!",
      "last_message_at": "2025-11-15T14:30:00Z",
      "message_count": 15,
      "created_at": "2025-11-15T14:00:00Z"
    }
  ],
  "total": 45,
  "page": 1,
  "total_pages": 3
}
```

```
GET /api/conversations/{id}
Response: 200 OK
{
  "conversation": {
    "id": 1,
    "phone_number": "+6281234567890",
    "is_sales": false,
    "created_at": "2025-11-15T14:00:00Z"
  },
  "messages": [
    {
      "id": 1,
      "sender": "+6281234567890",
      "content": "Ada Toyota budget 200 juta?",
      "direction": "inbound",
      "created_at": "2025-11-15T14:00:00Z"
    },
    {
      "id": 2,
      "sender": "BOT",
      "content": "Tentu! Saya menemukan beberapa Toyota...",
      "direction": "outbound",
      "created_at": "2025-11-15T14:00:05Z"
    }
  ],
  "total_messages": 15
}
```

#### WhatsApp Management API (Already Implemented)

```
GET /admin/whatsapp/status
POST /admin/whatsapp/pair
POST /admin/whatsapp/disconnect
POST /admin/whatsapp/test
GET /admin/whatsapp/qr/{tenant_id}
```

### Authentication & Authorization

**Current State:**
- ✅ JWT framework exists: `pkg/security/jwt.go`
- ❌ Not enforced on routes yet

**MVP Approach:**
- Use session-based auth atau basic auth untuk admin panel
- Tenant extraction via domain (existing middleware)
- No public access ke admin routes

**Future (Post-MVP):**
- Implement JWT authentication
- Role-based access (admin, viewer, sales)
- API keys untuk programmatic access

### Multi-Tenancy Architecture

**Existing Implementation (No Changes Needed):**
- ✅ Middleware: `internal/middleware/tenant.go` extracts tenant from domain
- ✅ Context: Tenant ID injected into request context
- ✅ Database: All queries filtered by tenant_id
- ✅ WhatsApp: Separate client instance per tenant

**Critical Requirements:**
- All new handlers MUST use `model.GetTenantID(r.Context())` untuk retrieve tenant
- All database queries MUST filter by tenant_id
- Never expose cross-tenant data
- Admin UI MUST clearly show which tenant is logged in

### Permissions & Roles

**MVP: Single Admin Role**
- All authenticated users can access all admin features
- Tenant isolation adalah primary security boundary

**Future Role Matrix:**

| Role | Sales Management | Conversations | WhatsApp Pairing | Analytics |
|------|-----------------|---------------|------------------|-----------|
| Super Admin | ✅ | ✅ | ✅ | ✅ |
| Admin | ✅ | ✅ | ✅ | ✅ |
| Sales Manager | ✅ | ✅ | ❌ | ✅ |
| Viewer | ❌ | ✅ (read-only) | ❌ | ✅ |

---

## User Experience Principles

### Admin Panel UX Goals

1. **Simplicity First**
   - Minimal clicks untuk common tasks
   - Clear visual hierarchy
   - No cluttered interface

2. **Real-time Feedback**
   - HTMX untuk partial updates
   - Loading states untuk async operations
   - Success/error notifications yang jelas

3. **Mobile-Friendly (Secondary)**
   - Responsive layout dengan Tailwind
   - Touch-friendly buttons dan forms
   - Readable pada screen kecil

4. **Consistency**
   - Reuse layout dari existing admin pages
   - Consistent table styles
   - Consistent form patterns

### Key Interactions

#### 1. WhatsApp Pairing Flow
```
[Admin clicks "Pair WhatsApp"]
  → [Enter phone number in form]
  → [Submit]
  → [QR code displayed in modal/section]
  → [Real-time status updates via HTMX polling]
  → [Success message when paired]
  → [Status changes to "Connected" with green indicator]
```

#### 2. Add Sales Team Member
```
[Admin clicks "Add Sales"]
  → [Form appears (phone, name, role)]
  → [Submit]
  → [Validation (phone format, uniqueness)]
  → [Success notification]
  → [Table updates with new entry via HTMX]
```

#### 3. View Conversation
```
[Admin sees list of conversations]
  → [Clicks on conversation row]
  → [Detail page loads / modal opens]
  → [Message history displayed chronologically]
  → [Visual distinction: customer (left), bot (right)]
  → [Can scroll untuk load older messages (HTMX infinite scroll)]
```

#### 4. Test Bot
```
[Admin on WhatsApp Management page]
  → [Enters test message: "Ada Toyota Avanza?"]
  → [Clicks "Send Test"]
  → [Loading indicator]
  → [Response displayed in UI]
  → [Can verify function was called (search executed)]
```

---

## Functional Requirements

### FR-1: Sales Team Management

**FR-1.1: Add Sales Team Member**
- **Actor:** Tenant Admin
- **Precondition:** Admin logged in, on Sales Management page
- **Flow:**
  1. Admin clicks "Tambah Sales" button
  2. Form muncul dengan fields: Phone Number (required), Name (required), Role (optional)
  3. Admin mengisi form
  4. System validates:
     - Phone number format (E.164 atau lokal Indonesia)
     - Phone number belum ada untuk tenant ini
     - Name tidak kosong
  5. System menyimpan ke table `sales` dengan tenant_id dari context
  6. System menampilkan success notification
  7. Table di-refresh via HTMX untuk show new entry
- **Postcondition:** Sales member tersimpan, visible di list
- **Exception:**
  - Duplicate phone number → Error: "Nomor sudah terdaftar"
  - Invalid format → Error: "Format nomor tidak valid"

**FR-1.2: List Sales Team Members**
- **Actor:** Tenant Admin
- **Precondition:** Admin logged in
- **Flow:**
  1. Admin navigates to `/admin/sales`
  2. System fetches sales members untuk tenant dari database
  3. System renders table dengan kolom: Name, Phone Number, Role, Created At, Actions
  4. Jika tidak ada data, tampilkan empty state dengan CTA "Tambah Sales"
- **Postcondition:** List ditampilkan
- **Exception:** Database error → Show error message dengan retry option

**FR-1.3: Delete Sales Team Member**
- **Actor:** Tenant Admin
- **Precondition:** Sales member exists
- **Flow:**
  1. Admin clicks "Delete" icon/button pada row
  2. Confirmation modal muncul: "Hapus [Name] dari sales team?"
  3. Admin confirms
  4. System deletes dari database (soft delete atau hard delete)
  5. Row di-remove dari table via HTMX
  6. Success notification
- **Postcondition:** Sales member removed dari list dan database
- **Exception:** Delete fails → Show error, row tetap ada

---

### FR-2: Conversation Monitoring

**FR-2.1: List Conversations**
- **Actor:** Tenant Admin
- **Precondition:** Admin logged in
- **Flow:**
  1. Admin navigates to `/admin/conversations`
  2. System fetches conversations untuk tenant dengan pagination (default 20 per page)
  3. System renders list dengan kolom:
     - Phone Number
     - Type badge (Customer / Sales)
     - Last Message (truncated)
     - Last Message Time
     - Message Count
     - Actions (View Details)
  4. Sort by last_message_at DESC (most recent first)
  5. Jika tidak ada conversations, show empty state
- **Postcondition:** Conversations displayed
- **UI Enhancement:**
  - Filter buttons: All / Customer / Sales
  - Search box untuk phone number
  - Pagination controls (Previous / Next)

**FR-2.2: View Conversation Detail**
- **Actor:** Tenant Admin
- **Precondition:** Conversation exists
- **Flow:**
  1. Admin clicks "View" atau conversation row
  2. System fetches conversation detail dan messages
  3. System renders detail page/modal dengan:
     - Header: Phone number, type badge, created date
     - Message thread (chronological, oldest to newest)
     - Visual styling:
       - Inbound messages (customer/sales) aligned left, grey background
       - Outbound messages (bot) aligned right, blue background
     - Timestamp untuk each message
     - Sender identification (phone number atau "BOT")
  4. Auto-scroll to latest message
  5. Option untuk load older messages jika > 50 messages (HTMX load more)
- **Postcondition:** Full conversation visible
- **Exception:**
  - Conversation not found → 404 page
  - Access dari tenant lain → 403 Forbidden

**FR-2.3: Filter Conversations**
- **Actor:** Tenant Admin
- **Precondition:** On conversations list page
- **Flow:**
  1. Admin clicks filter button (All / Customer / Sales)
  2. System re-fetches dengan filter parameter
  3. List updates via HTMX
  4. Active filter highlighted
- **Postcondition:** Filtered list displayed

**FR-2.4: Search Conversation by Phone**
- **Actor:** Tenant Admin
- **Precondition:** On conversations list page
- **Flow:**
  1. Admin types phone number di search box
  2. System searches (debounced, after typing stops untuk 500ms)
  3. Results update via HTMX
  4. Jika tidak ada results, show "Tidak ditemukan"
- **Postcondition:** Matching conversations displayed

---

### FR-3: WhatsApp Management UI Enhancement

**FR-3.1: Enhanced Pairing Interface**
- **Actor:** Tenant Admin
- **Precondition:** Admin on `/admin/whatsapp`
- **Flow:**
  1. Page loads, fetches status via GET `/admin/whatsapp/status`
  2. If not connected:
     - Show form untuk phone number entry
     - Button "Mulai Pairing"
  3. Admin enters phone number, clicks "Mulai Pairing"
  4. POST to `/admin/whatsapp/pair` dengan phone number
  5. QR code ditampilkan (via `/admin/whatsapp/qr/{tenant_id}`)
  6. Status polling starts (HTMX polls `/admin/whatsapp/status` every 3 seconds)
  7. When pairing success detected:
     - QR code hidden
     - Success message: "WhatsApp berhasil terhubung!"
     - Status indicator changes to "Connected" (green)
  8. If timeout (2 minutes) atau error:
     - Show error message
     - Option to retry
- **Postcondition:** WhatsApp paired atau error displayed

**FR-3.2: Connection Status Display**
- **Actor:** Tenant Admin
- **Precondition:** On WhatsApp page
- **Flow:**
  1. Page shows current status:
     - **Connected:** Green indicator, phone number shown, last active time
     - **Pairing in progress:** Yellow indicator, QR code shown
     - **Disconnected:** Red indicator, option to pair
  2. Real-time updates via HTMX polling
- **Postcondition:** Current status always visible

**FR-3.3: Test Message Interface**
- **Actor:** Tenant Admin
- **Precondition:** WhatsApp connected
- **Flow:**
  1. Test message section visible dengan:
     - Textarea untuk message input
     - "Send Test Message" button
  2. Admin types test message (e.g., "Ada Toyota Avanza?")
  3. Admin clicks button
  4. POST to `/admin/whatsapp/test` dengan message body
  5. Loading indicator shown
  6. Response displayed di UI:
     - "Message sent to: [bot phone number]"
     - "Expected response: Check your WhatsApp"
     - Atau show actual response jika API returns it
  7. Success notification
- **Postcondition:** Test message sent
- **Exception:**
  - Not connected → Error: "WhatsApp belum terhubung"
  - Send failed → Error with details

**FR-3.4: Disconnect WhatsApp**
- **Actor:** Tenant Admin
- **Precondition:** WhatsApp connected
- **Flow:**
  1. Admin clicks "Disconnect" button
  2. Confirmation modal: "Yakin ingin disconnect WhatsApp?"
  3. Admin confirms
  4. POST to `/admin/whatsapp/disconnect`
  5. Status updates to "Disconnected"
  6. UI updates untuk show pairing option
- **Postcondition:** WhatsApp disconnected

---

### FR-4: Admin Dashboard Integration

**FR-4.1: Navigation Menu Update**
- **Actor:** Tenant Admin
- **Precondition:** Admin logged in
- **Flow:**
  1. Admin sees left sidebar menu with:
     - Dashboard (existing)
     - Cars (existing)
     - **Sales Team** (new)
     - **Conversations** (new)
     - WhatsApp (existing, enhanced)
  2. Current page highlighted
  3. Click navigates to respective page
- **Postcondition:** Easy navigation between sections

**FR-4.2: Dashboard Quick Stats**
- **Actor:** Tenant Admin
- **Precondition:** On dashboard (`/admin/dashboard`)
- **Flow:**
  1. Dashboard shows cards:
     - **WhatsApp Status:** Connected/Disconnected dengan indicator
     - **Total Conversations:** Count (last 30 days)
     - **Sales Team:** Count of members
     - **Cars:** Count (existing)
  2. Quick action buttons:
     - "Manage WhatsApp"
     - "View Conversations"
- **Postcondition:** Overview visible at a glance

**FR-4.3: Recent Conversations Widget**
- **Actor:** Tenant Admin
- **Precondition:** On dashboard
- **Flow:**
  1. Widget shows 5 most recent conversations
  2. For each: Phone number, last message preview, time ago
  3. "View All" link navigates to `/admin/conversations`
- **Postcondition:** Quick access to recent activity

---

### FR-5: Error Handling & Validation

**FR-5.1: Form Validation**
- All forms validate input before submission
- Client-side validation (HTML5 + Alpine.js)
- Server-side validation (Go handlers)
- Error messages displayed inline near fields

**FR-5.2: API Error Handling**
- All API calls wrapped dengan error handling
- HTTP error codes properly mapped:
  - 400 Bad Request → Validation errors
  - 403 Forbidden → Access denied (tenant isolation)
  - 404 Not Found → Resource tidak ditemukan
  - 500 Internal Server Error → Generic error dengan retry option
- Error messages user-friendly dalam Bahasa Indonesia

**FR-5.3: Network Error Handling**
- HTMX timeout handling (10 seconds)
- Retry mechanism untuk critical operations
- Loading states untuk semua async operations
- Graceful degradation jika JavaScript disabled (form masih bisa submit)

---

## Non-Functional Requirements

### Performance

**PR-1: Response Time**
- API endpoints: < 200ms (p95)
- Page load: < 2s (full page render)
- HTMX partial updates: < 500ms
- QR code generation: < 1s

**PR-2: Database Queries**
- All queries use indexes (tenant_id, created_at, phone_number)
- Pagination untuk large datasets (conversations, messages)
- Limit query results (max 100 messages per request)

**PR-3: WhatsApp Status Polling**
- Polling interval: 3 seconds (only saat pairing)
- Stop polling after pairing success atau timeout
- Tidak over-poll untuk avoid rate limiting

### Security

**SR-1: Multi-Tenant Isolation (CRITICAL)**
- ✅ Middleware enforcement (existing): `internal/middleware/tenant.go`
- ✅ Context-based tenant ID (existing)
- 🔒 All new handlers MUST call `model.GetTenantID(r.Context())`
- 🔒 All repository methods MUST filter by tenant_id
- 🔒 Unit tests untuk verify tenant isolation

**SR-2: Input Validation**
- Phone number format validation (prevent injection)
- SQL injection prevention via parameterized queries (existing pattern)
- XSS prevention via HTML escaping di templates (Go templates default)
- CSRF protection untuk forms (add CSRF middleware)

**SR-3: Authentication (MVP: Basic)**
- Session-based auth atau basic auth untuk admin panel
- HTTPS required di production (Nginx SSL termination)
- Secure cookies (HttpOnly, Secure flags)

**SR-4: Rate Limiting**
- API rate limiting: 100 requests per minute per tenant
- WhatsApp pairing: max 5 attempts per hour per tenant
- Test message: max 20 per hour per tenant

### Scalability

**SC-1: Database Scalability**
- Conversations table bisa grow hingga 100K+ rows per tenant
- Proper indexes untuk prevent slow queries:
  - Index on (tenant_id, created_at)
  - Index on (tenant_id, phone_number)
  - Index on (tenant_id, is_sales)
- Consider archiving old conversations (> 6 months) di future

**SC-2: Concurrent Users**
- Support 50 concurrent admin users (5 tenants × 10 admins each)
- No locking issues pada conversations (read-heavy)
- Sales table small (< 100 rows per tenant)

**SC-3: WhatsApp Client Management**
- Multi-tenant WhatsApp client already implemented
- One WhatsApp session per tenant
- Memory footprint: ~10MB per active session

### Integration

**IN-1: WhatsApp Integration (Existing)**
- ✅ Whatsmeow library: https://github.com/tulir/whatsmeow
- ✅ Multi-device support
- ✅ QR code pairing
- ✅ Message sending/receiving
- ✅ Media support (images)

**IN-2: LLM Integration (Existing)**
- ✅ Z.AI provider: https://z.ai/
- ✅ Model: glm-4.6
- ✅ Function calling support
- ✅ OpenAI-compatible API

**IN-3: Frontend Integration**
- HTMX untuk dynamic updates
- Alpine.js untuk minimal interactivity
- Tailwind CSS v4 untuk styling
- No heavy JavaScript framework needed

---

## Implementation Planning

### Epic Breakdown

Karena requirements cukup lengkap dan context limit 200k, PRD ini harus di-breakdown menjadi epics dan stories.

**Proposed Epic Structure:**

1. **Epic 1: Sales Team Management**
   - Story 1.1: Implement sales handler (Create, List, Delete)
   - Story 1.2: Build sales management UI (templates)
   - Story 1.3: Wire up routes dan test end-to-end

2. **Epic 2: Conversation Monitoring**
   - Story 2.1: Implement conversation handler (List, Detail)
   - Story 2.2: Build conversation list UI dengan filters
   - Story 2.3: Build conversation detail UI dengan message thread
   - Story 2.4: Add pagination dan search

3. **Epic 3: WhatsApp Management UI Enhancement**
   - Story 3.1: Enhance WhatsApp template (QR display, status polling)
   - Story 3.2: Improve test message interface
   - Story 3.3: Add disconnect confirmation modal
   - Story 3.4: Add error handling dan validation

4. **Epic 4: Admin Dashboard Integration**
   - Story 4.1: Update dashboard dengan quick stats
   - Story 4.2: Add recent conversations widget
   - Story 4.3: Update navigation menu
   - Story 4.4: Ensure responsive layout

5. **Epic 5: Testing & Refinement**
   - Story 5.1: End-to-end testing (pairing → chat → view conversation)
   - Story 5.2: Multi-tenant isolation testing
   - Story 5.3: Performance testing (load 1000 conversations)
   - Story 5.4: Security audit (input validation, tenant leakage)

**Next Step:** Run `/bmad:bmm:workflows:create-epics-and-stories` untuk generate detailed epic dan story breakdown.

---

## Integration Points with Existing Code

### Leveraging Existing Implementation

**✅ No Changes Needed:**
1. **Database Schema** (`migrations/`)
   - Table `sales` sudah ada
   - Table `conversations` dan `messages` sudah ada
   - Indexes sudah proper

2. **Repositories** (`internal/repository/`)
   - `sales_repository.go` sudah lengkap (Create, List, Delete, IsSales)
   - `conversation_repository.go` sudah lengkap (GetOrCreate, GetMessages, AddMessage)

3. **Middleware** (`internal/middleware/`)
   - `tenant.go` sudah handle tenant extraction dari domain

4. **WhatsApp Service** (`internal/service/whatsapp_service.go`)
   - `ProcessIncomingMessage` sudah orchestrate bot + database
   - Tidak perlu diubah

5. **LLM Bot** (`internal/llm/bot.go`)
   - Function calling sudah fully implemented
   - Tidak perlu diubah

**🔧 Need Integration:**
1. **Handlers** (`internal/handler/`)
   - Create `sales_handler.go` (use existing repository)
   - Create `conversation_handler.go` (use existing repository)
   - Enhance `whatsapp_handler.go` (minimal changes untuk response format)

2. **Routes** (`cmd/api/main.go`)
   - Uncomment lines 224-228 (sales routes)
   - Uncomment lines 231-234 (conversation routes)
   - Wire up new handlers

3. **Templates** (`templates/admin/`)
   - Create `sales.html`
   - Create `conversations.html`
   - Enhance `whatsapp.html`
   - Update `layout.html` (navigation)
   - Update `dashboard.html` (quick stats widget)

### Data Flow Example: View Conversation

```
User Request: GET /admin/conversations/123
    ↓
1. Chi Router (cmd/api/main.go:231)
    ↓
2. TenantExtractor Middleware
    ↓ (injects tenant_id = 5 into context)
3. ConversationHandler.Get (internal/handler/conversation_handler.go)
    ↓ calls
4. ConversationRepository.Get(ctx, conversationID=123)
    ↓ executes SQL
5. SELECT * FROM conversations WHERE id = $1 AND tenant_id = $2
    ↓ (automatically filtered by tenant_id from context)
6. ConversationRepository.GetMessages(ctx, conversationID=123, limit=50)
    ↓ executes SQL
7. SELECT * FROM messages WHERE conversation_id = $1 ORDER BY created_at ASC LIMIT $2
    ↓
8. Handler renders template: conversations_detail.html
    ↓
9. HTML response to user
```

**Key Point:** Tenant isolation terjadi automatically via middleware + repository pattern yang sudah ada. New handlers hanya perlu follow existing pattern.

---

## References

- **Brownfield Documentation:**
  - [Project Overview](./project-overview.md)
  - [Architecture](./architecture.md)
  - [Source Tree Analysis](./source-tree-analysis.md)
  - [Development Guide](./development-guide.md)
  - [Documentation Index](./index.md)

- **Source Code Audit (2025-11-15):**
  - Identified missing handlers: Sales, Conversations
  - Confirmed WhatsApp bot 100% functional
  - Confirmed LLM integration 100% functional
  - Identified commented routes in main.go

- **Existing Implementation:**
  - Repository: https://github.com/riz/auto-lmk
  - Main branch: `main`
  - Current version: v1.0 (production-ready)

---

## Next Steps

### Immediate Actions

1. ✅ **PRD Complete** - This document
2. 🔄 **Epic Breakdown** - Run: `/bmad:bmm:workflows:create-epics-and-stories`
3. 🔄 **Architecture Review** - Run: `/bmad:bmm:workflows:create-architecture` (verify integration points)
4. 🔄 **Sprint Planning** - Run: `/bmad:bmm:workflows:sprint-planning`
5. 🔄 **Implementation** - Run: `/bmad:bmm:agents:dev` untuk execute stories

### Development Sequence

**Week 1: MVP Implementation**
- Day 1-2: Sales Team Management (Epic 1)
- Day 3-4: Conversation Monitoring (Epic 2)
- Day 5: WhatsApp UI Enhancement (Epic 3)
- Day 6: Dashboard Integration (Epic 4)
- Day 7: Testing & Refinement (Epic 5)

**Week 2: Testing & Deployment**
- End-to-end testing
- Security testing (tenant isolation)
- Performance testing
- Documentation update
- Deploy to staging
- User acceptance testing
- Deploy to production

---

## Acceptance Criteria (Overall)

**✅ PRD Considered Complete When:**

1. **Tenant admin dapat setup WhatsApp bot dari awal sampai test chat dalam < 5 menit**
   - Pair WhatsApp via QR code
   - Add sales team members
   - Send test message
   - Receive AI response
   - View conversation history

2. **All CRUD operations berfungsi dengan baik:**
   - Create sales member ✅
   - List sales members ✅
   - Delete sales member ✅
   - List conversations ✅
   - View conversation detail ✅
   - Filter conversations ✅

3. **UI user-friendly dan responsive:**
   - Works di desktop dan mobile
   - HTMX updates smooth tanpa page reload
   - Error messages jelas
   - Loading states visible

4. **Multi-tenant isolation verified:**
   - Tenant A tidak bisa lihat data Tenant B
   - All queries filtered by tenant_id
   - No cross-tenant data leakage

5. **WhatsApp bot berfungsi end-to-end:**
   - Customer sends message → Bot responds
   - Function calling works (search cars, get details, send images)
   - Conversation tersimpan di database
   - Admin bisa view conversation

---

_PRD ini merupakan brownfield enhancement yang memanfaatkan WhatsApp bot AI-powered yang sudah fully implemented. Focus adalah pada melengkapi admin interface untuk enable testing dan management, bukan membangun bot dari scratch._

_Dengan melengkapi missing handlers dan UI templates, Auto LMK akan menjadi platform yang production-ready dengan full self-service capability untuk dealerships._

**Created through collaborative discovery between Yopi and AI Product Manager (John).**

---

**Auto LMK - Enabling dealerships to test and leverage AI-powered WhatsApp bot with confidence.** 🚗💬🤖
