# auto-lmk - Epic Breakdown

**Author:** Yopi
**Date:** 2025-11-15
**Project Level:** Brownfield Enhancement
**Target Scale:** Multi-Tenant SaaS (5-10 tenants in MVP)

---

## Overview

This document provides the complete epic and story breakdown for auto-lmk, decomposing the requirements from the [PRD](./PRD-Admin-Tenant.md) into implementable stories.

**Living Document Notice:** This is the initial version. It will be updated after UX Design and Architecture workflows add interaction and technical details to stories.

### Epics Summary

This epic breakdown transforms the PRD's 24 functional requirements into 6 cohesive epics containing 28 detailed user stories.

**Epic 1: Sales Team Management** (3 stories)
- Enable admins to manage sales team members who have elevated privileges in the WhatsApp bot
- Foundation for role-based bot capabilities

**Epic 2: Conversation Monitoring** (4 stories)
- Provide admins visibility into all customer and sales WhatsApp conversations
- Includes filtering, searching, and detail viewing capabilities

**Epic 3: WhatsApp Management Enhancement** (4 stories)
- Streamline WhatsApp pairing with QR code display and real-time status updates
- Add test message interface for bot verification

**Epic 4: Admin Dashboard Integration** (3 stories)
- Unify all admin functions with consistent navigation and at-a-glance insights
- Display quick stats and recent activity widgets

**Epic 5: Error Handling & Validation** (3 stories)
- Ensure production-ready quality with comprehensive validation and error recovery
- Cross-cutting concern applied across all features

**Epic 6: Sales Car Upload via WhatsApp** ⭐ (7 stories)
- Game-changing feature: Sales can upload cars directly from field via WhatsApp chat
- AI-powered parsing of car details with role-based access control
- Eliminates admin overhead for inventory management

**Epic 7: Public Website Enhancement** (6 stories)
- Comprehensive public website features for car catalog
- Advanced search, analytics, WhatsApp integration, blog, branding, and showroom info
- Complete customer experience from browsing to contact

---

## Functional Requirements Inventory

**FR-1: Sales Team Management**
- FR-1.1: Add Sales Team Member
- FR-1.2: List Sales Team Members
- FR-1.3: Delete Sales Team Member

**FR-2: Conversation Monitoring**
- FR-2.1: List Conversations
- FR-2.2: View Conversation Detail
- FR-2.3: Filter Conversations
- FR-2.4: Search Conversation by Phone

**FR-3: WhatsApp Management UI Enhancement**
- FR-3.1: Enhanced Pairing Interface
- FR-3.2: Connection Status Display
- FR-3.3: Test Message Interface
- FR-3.4: Disconnect WhatsApp

**FR-4: Admin Dashboard Integration**
- FR-4.1: Navigation Menu Update
- FR-4.2: Dashboard Quick Stats
- FR-4.3: Recent Conversations Widget

**FR-5: Error Handling & Validation**
- FR-5.1: Form Validation
- FR-5.2: API Error Handling
- FR-5.3: Network Error Handling

**FR-6: Sales Car Upload via WhatsApp (Role-Based)**
- FR-6.1: Role Detection - Bot detects if sender is Sales or Customer
- FR-6.2: Upload Car Photos (Sales Only) - Sales upload multiple photos via WhatsApp
- FR-6.3: Submit Car Details via Text - Sales submit details in structured format
- FR-6.4: AI Parse Car Details - Bot parse text into structured data (brand, model, price, etc)
- FR-6.5: Auto-save to Database - Save to cars, car_photos, car_specs tables
- FR-6.6: Confirmation & Catalog Link - Bot confirms save and sends catalog URL
- FR-6.7: Reject Customer Upload - Bot politely rejects if customer tries to upload

**FR-7: Public Website Enhancement**
- FR-7.1: Advanced Search & Filtering - Real-time car search with multiple filters
- FR-7.2: Search Analytics System - Track and analyze customer search behavior
- FR-7.3: WhatsApp Integration - Floating WA button with bot and fallback
- FR-7.4: Blog Management System - Admin blog creation with AI assistance
- FR-7.5: Branding Customization - Tenant-specific logo, colors, and text
- FR-7.6: Showroom Information - Address, location, and contact details

**Total:** 30 Functional Requirements (17 original + 7 new for Sales Upload + 6 new for Public Website)

---

## FR Coverage Map

| Epic | FRs Covered | Count |
|------|-------------|-------|
| **Epic 1: Sales Team Management** | FR-1.1, FR-1.2, FR-1.3 | 3 |
| **Epic 2: Conversation Monitoring** | FR-2.1, FR-2.2, FR-2.3, FR-2.4 | 4 |
| **Epic 3: WhatsApp Management Enhancement** | FR-3.1, FR-3.2, FR-3.3, FR-3.4 | 4 |
| **Epic 4: Admin Dashboard Integration** | FR-4.1, FR-4.2, FR-4.3 | 3 |
| **Epic 5: Error Handling & Validation** | FR-5.1, FR-5.2, FR-5.3 | 3 |
| **Epic 6: Sales Car Upload via WhatsApp** | FR-6.1, FR-6.2, FR-6.3, FR-6.4, FR-6.5, FR-6.6, FR-6.7 | 7 |
| **Epic 7: Public Website Enhancement** | FR-7.1, FR-7.2, FR-7.3, FR-7.4, FR-7.5, FR-7.6 | 6 |
| **Total** | **All FRs** | **30** |

---

## Epic 7: Public Website Enhancement

**Goal:** Create an engaging and user-friendly public website that showcases the car catalog effectively, attracting more potential customers and improving the overall user experience.

**Value:** The public website is the first point of contact for customers. A modern, fast, and intuitive interface will increase engagement, reduce bounce rates, and drive more inquiries through the WhatsApp bot.

**FRs Covered:** FR-7.1

---

### Story 7.1: Improve Public Catalog Homepage UI

**As a** website visitor,
**I want** an improved homepage UI for the car catalog,
**so that** I can easily browse available cars with better visual appeal and navigation.

**Acceptance Criteria:**

**Given** I visit the public homepage (`/`)
**When** the page loads
**Then** I see:
- Hero section with compelling background and search functionality
- Featured cars grid (3-6 cars) with photos, prices, and key specs
- Quick filter buttons ("All Cars", "Under 200M", "Automatic", "Manual")
- Call-to-action buttons ("Browse All Cars", "Contact Sales")

**And** search functionality works:
- Real-time search as I type (HTMX-powered)
- Filter buttons update results instantly
- Empty state shows helpful message when no results

**And** each car card displays:
- High-quality thumbnail photo
- Brand, model, year
- Formatted price (Rp 185.000.000)
- Key specs (transmission, fuel type)
- "View Details" link to individual car page

**And** page performance meets standards:
- Loads in < 2 seconds
- Mobile responsive
- SEO optimized with proper meta tags
- Consistent with existing design system

**Prerequisites:** Existing car data in database, basic homepage template exists

**Technical Notes:**
- Template: `templates/pages/home.html`
- Handler: Update `page_handler.go` for homepage data
- API: New `GET /api/cars/search` endpoint
- Repository: Extend `car_repository.go` with search methods
- Styling: Tailwind CSS consistent with admin interface
- HTMX: For dynamic search and filtering
- Images: Use car_photos table with thumbnail generation

---

## Epic 1: Sales Team Management

**Goal:** Enable tenant admins to manage sales team members who have special privileges in WhatsApp bot (can upload cars, internal communication tracking).

**Value:** Foundation for role-based bot behavior - sales team members are distinguished from regular customers in conversations and granted additional capabilities.

**FRs Covered:** FR-1.1, FR-1.2, FR-1.3

---

### Story 1.1: Implement Sales Handler with CRUD Endpoints

**As a** backend developer,
**I want** to create sales_handler.go with Create, List, Delete endpoints,
**So that** tenant admins can manage sales team members via API.

**Acceptance Criteria:**

**Given** the sales repository already exists at `internal/repository/sales_repository.go`
**When** I implement `internal/handler/sales_handler.go`
**Then** the handler exposes three endpoints:
- `POST /api/sales` - Create sales member with phone_number, name, role
- `GET /api/sales` - List all sales members for current tenant
- `DELETE /api/sales/{id}` - Delete sales member by ID

**And** all endpoints use `model.GetTenantID(r.Context())` for tenant isolation
**And** POST endpoint validates phone number format (E.164 or local Indonesia format)
**And** POST endpoint checks for duplicate phone numbers within tenant
**And** DELETE endpoint returns 404 if sales member not found or belongs to different tenant
**And** all responses return proper JSON with appropriate HTTP status codes (201, 200, 204, 400, 404, 500)
**And** error messages are in Bahasa Indonesia

**Prerequisites:** None (repository already exists)

**Technical Notes:**
- Follow existing handler pattern from `car_handler.go`
- Use `sales_repository.go` methods: Create(), List(), Delete(), IsSales()
- Phone validation regex: `^\+?[1-9]\d{1,14}$` or `^08\d{8,11}$`
- Tenant ID from context: `tenantID, err := model.GetTenantID(r.Context())`
- Return structured errors: `{"error": "Nomor sudah terdaftar"}`
- File location: `internal/handler/sales_handler.go` (~150-200 lines)

---

### Story 1.2: Build Sales Management UI Template

**As a** tenant admin,
**I want** a clean UI to manage sales team members,
**So that** I can easily add, view, and remove sales team members.

**Acceptance Criteria:**

**Given** I am logged in as tenant admin
**When** I navigate to `/admin/sales`
**Then** I see a page with:
- Page title: "Sales Team Management"
- "Tambah Sales" button at top right
- Table with columns: Name, Phone Number, Role, Created At, Actions (Delete icon)
- Empty state if no sales members: "Belum ada sales team. Klik 'Tambah Sales' untuk menambahkan."

**And** when I click "Tambah Sales" button
**Then** a form appears (modal or inline) with fields:
- Phone Number (required, with placeholder: "+628123456789")
- Name (required, with placeholder: "John Doe")
- Role (optional, with placeholder: "Sales Executive")
- Submit button: "Simpan"
- Cancel button: "Batal"

**And** when I submit the form with valid data
**Then** HTMX sends POST to `/api/sales`
**And** on success, form closes and table updates with new entry (via HTMX swap)
**And** success toast notification appears: "Sales member berhasil ditambahkan"

**And** when I click delete icon on a row
**Then** confirmation modal appears: "Hapus [Name] dari sales team?"
**And** when I confirm, HTMX sends DELETE to `/api/sales/{id}`
**And** on success, row is removed from table
**And** success toast notification appears: "Sales member berhasil dihapus"

**And** all error messages display inline near form fields
**And** page is responsive (mobile-friendly with Tailwind breakpoints)

**Prerequisites:** Story 1.1 (API endpoints must exist)

**Technical Notes:**
- Template file: `templates/admin/sales.html`
- Extend: `templates/admin/layout.html`
- Use HTMX attributes:
  - `hx-post="/api/sales"` on form
  - `hx-delete="/api/sales/{id}"` on delete button
  - `hx-target="#sales-table"` for table updates
  - `hx-swap="outerHTML"` for row replacement
- Use Alpine.js for modal show/hide
- Tailwind classes for styling (consistent with existing admin pages)
- Form validation: HTML5 `required`, `pattern` for phone
- Toast notifications: Use existing notification system or Alpine.js component

---

### Story 1.3: Wire Sales Routes and Test End-to-End

**As a** developer,
**I want** to activate sales routes in main.go and verify complete workflow,
**So that** the sales management feature is fully functional.

**Acceptance Criteria:**

**Given** Story 1.1 and 1.2 are complete
**When** I uncomment sales routes in `cmd/api/main.go` (lines 224-228)
**Then** the following routes are active:
```go
r.Route("/sales", func(r chi.Router) {
    r.Post("/", salesHandler.Create)
    r.Get("/", salesHandler.List)
    r.Delete("/{id}", salesHandler.Delete)
})
```

**And** I initialize salesHandler in setupRouter:
```go
salesHandler := handler.NewSalesHandler(salesRepo)
```

**And** when I run the application and test the complete workflow:
1. Navigate to `/admin/sales` - Page loads successfully
2. Click "Tambah Sales" - Form appears
3. Submit with valid data - Success, table updates
4. Submit with duplicate phone - Error: "Nomor sudah terdaftar"
5. Submit with invalid phone - Error: "Format nomor tidak valid"
6. Click delete on entry - Confirmation appears
7. Confirm delete - Success, row removed
8. Verify tenant isolation - Sales from tenant A not visible to tenant B

**And** all operations complete in < 500ms (API response time)
**And** no console errors in browser
**And** multi-tenant isolation verified via manual testing

**Prerequisites:** Story 1.1, Story 1.2

**Technical Notes:**
- File to edit: `cmd/api/main.go`
- Lines to uncomment: 224-228
- Add handler initialization before route setup
- Test with at least 2 different tenant domains (e.g., tenant1.localhost, tenant2.localhost)
- Manual test checklist:
  - [ ] Create sales member works
  - [ ] List shows only tenant's sales
  - [ ] Delete works with confirmation
  - [ ] Duplicate phone rejected
  - [ ] Invalid phone rejected
  - [ ] Tenant isolation verified
- Consider adding basic unit tests for handler (optional for MVP)

---

## Epic 2: Conversation Monitoring

**Goal:** Enable tenant admins to monitor all WhatsApp conversations (customer inquiries & sales communication) with powerful filtering, search, and detail viewing capabilities.

**Value:** Visibility into customer interactions and sales team conversations, enabling better support, quality control, and lead tracking.

**FRs Covered:** FR-2.1, FR-2.2, FR-2.3, FR-2.4

---

### Story 2.1: Implement Conversation List Handler with Pagination

**As a** backend developer,
**I want** to create conversation list endpoint with pagination and filtering,
**So that** admins can browse conversations efficiently.

**Acceptance Criteria:**

**Given** the conversation repository exists at `internal/repository/conversation_repository.go`
**When** I implement `internal/handler/conversation_handler.go` with List method
**Then** the handler exposes:
- `GET /api/conversations?page=1&limit=20&type=customer|sales|all&sort=recent`

**And** the endpoint accepts query parameters:
- `page` (default: 1, min: 1)
- `limit` (default: 20, max: 100)
- `type` (default: "all", options: "customer", "sales", "all")
- `sort` (default: "recent", for future: "oldest", "most_messages")

**And** the response includes:
```json
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
  "limit": 20,
  "total_pages": 3
}
```

**And** conversations are filtered by tenant_id from context
**And** sorting is by last_message_at DESC (most recent first)
**And** type filter works correctly (customer: is_sales=false, sales: is_sales=true)
**And** pagination calculates total_pages correctly: `ceil(total / limit)`
**And** response time is < 200ms even with 1000+ conversations

**Prerequisites:** None (repository exists)

**Technical Notes:**
- File: `internal/handler/conversation_handler.go` (~200 lines)
- Use `conversationRepo.List(ctx, tenantID, page, limit, typeFilter)`
- Add repository method if not exists: `List(ctx context.Context, tenantID, page, limit int, typeFilter string) ([]Conversation, int, error)`
- SQL query needs:
  - JOIN with messages to get last_message and message_count
  - `WHERE tenant_id = $1` for isolation
  - `AND is_sales = $2` if type filter active
  - `ORDER BY last_message_at DESC`
  - `LIMIT $3 OFFSET $4` for pagination
- Use parameterized queries to prevent SQL injection
- Cache total count if performance becomes issue (future optimization)

---

### Story 2.2: Implement Conversation Detail Handler

**As a** backend developer,
**I want** to create conversation detail endpoint that returns full message history,
**So that** admins can view complete conversation threads.

**Acceptance Criteria:**

**Given** conversation exists in database
**When** I implement `GET /api/conversations/{id}`
**Then** the endpoint returns:
```json
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

**And** messages are ordered chronologically (oldest first: `ORDER BY created_at ASC`)
**And** endpoint enforces tenant isolation (404 if conversation belongs to different tenant)
**And** endpoint returns 404 if conversation ID not found
**And** messages are limited to 50 per request (for initial load, use pagination for older messages)
**And** response includes total_messages count for pagination UI

**Prerequisites:** Story 2.1

**Technical Notes:**
- Add method to handler: `Get(w http.ResponseWriter, r *http.Request)`
- Extract ID from URL: `chi.URLParam(r, "id")`
- Use repository methods:
  - `conversationRepo.GetByID(ctx, conversationID, tenantID)`
  - `conversationRepo.GetMessages(ctx, conversationID, limit, offset)`
- Verify tenant ownership before returning data
- For pagination support, add optional query params:
  - `?limit=50&offset=0` for loading older messages (HTMX infinite scroll)
- Message content may be long - ensure no truncation in database fetch

---

### Story 2.3: Build Conversation List UI with Filters

**As a** tenant admin,
**I want** to see all conversations in a filterable list,
**So that** I can quickly find specific customer or sales conversations.

**Acceptance Criteria:**

**Given** I am logged in as tenant admin
**When** I navigate to `/admin/conversations`
**Then** I see a page with:
- Page title: "Conversations"
- Filter buttons: "Semua" (active), "Customer", "Sales"
- Search box with placeholder: "Cari nomor telepon..."
- Table with columns: Phone Number, Type Badge, Last Message, Time, Message Count, Actions (View button)
- Pagination controls: "Previous" / "Next" (disabled appropriately)
- Empty state if no conversations: "Belum ada conversation."

**And** each row displays:
- Phone number (formatted: +62 812-3456-7890)
- Type badge: "Customer" (blue) or "Sales" (green)
- Last message (truncated to 50 chars with "...")
- Time: Relative time (e.g., "5 menit lalu", "2 jam lalu", "3 hari lalu")
- Message count badge (e.g., "15 pesan")
- "View" button

**And** when I click "Customer" filter
**Then** HTMX sends GET to `/api/conversations?type=customer`
**And** table updates to show only customer conversations
**And** "Customer" button is highlighted (active state)

**And** when I type in search box (debounced 500ms)
**Then** HTMX sends GET to `/api/conversations?search={phone}`
**And** table updates to show matching conversations
**And** if no results, show: "Tidak ditemukan"

**And** when I click "Next" pagination
**Then** HTMX sends GET to `/api/conversations?page=2`
**And** table updates with next page data
**And** "Previous" button becomes enabled

**And** page is responsive and mobile-friendly

**Prerequisites:** Story 2.1 (API endpoint)

**Technical Notes:**
- Template file: `templates/admin/conversations.html`
- HTMX attributes:
  - `hx-get="/api/conversations"` on filter buttons
  - `hx-get="/api/conversations?search={value}"` on search input
  - `hx-trigger="keyup changed delay:500ms"` for debounced search
  - `hx-target="#conversations-table"`
  - `hx-swap="outerHTML"`
- Alpine.js for:
  - Active filter state management
  - Search input state
- Tailwind utilities:
  - Badge: `bg-blue-100 text-blue-800` (customer), `bg-green-100 text-green-800` (sales)
  - Truncate: `truncate max-w-xs`
- Relative time: Use JavaScript or Go template helper (e.g., "5 minutes ago")
- Pagination: Track current page in URL query or Alpine state

---

### Story 2.4: Build Conversation Detail UI with Message Thread

**As a** tenant admin,
**I want** to view full conversation history with visual message thread,
**So that** I can understand complete customer interaction context.

**Acceptance Criteria:**

**Given** I am on conversation list page
**When** I click "View" button on a conversation
**Then** I navigate to `/admin/conversations/{id}` detail page

**And** the page displays:
- Header section:
  - Back button to list
  - Phone number (formatted)
  - Type badge (Customer/Sales)
  - Created date
- Message thread section:
  - Inbound messages (customer/sales):
    - Aligned left
    - Grey background (`bg-gray-100`)
    - Sender phone number shown
    - Timestamp below message
  - Outbound messages (bot):
    - Aligned right
    - Blue background (`bg-blue-100`)
    - "BOT" label shown
    - Timestamp below message
  - Messages ordered chronologically (oldest at top)
  - Auto-scroll to latest message on page load

**And** if conversation has > 50 messages
**Then** show "Load Older Messages" button at top
**And** when clicked, HTMX loads previous 50 messages
**And** messages are prepended to thread (infinite scroll upward)

**And** each message displays:
- Full message content (no truncation)
- Sender identification (phone or "BOT")
- Timestamp (formatted: "15 Nov 2025, 14:30")
- Direction indicator via alignment and color

**And** page layout is:
- Chat-like UI (similar to WhatsApp web)
- Mobile responsive
- Fixed header, scrollable message thread
- Message bubbles with max-width for readability

**Prerequisites:** Story 2.2 (API endpoint)

**Technical Notes:**
- Template file: `templates/admin/conversations_detail.html`
- Message thread structure:
```html
<div class="space-y-4 overflow-y-auto h-[600px]">
  {{range .Messages}}
    {{if eq .Direction "inbound"}}
      <!-- Left-aligned message -->
      <div class="flex justify-start">
        <div class="bg-gray-100 rounded-lg p-3 max-w-md">
          <p class="text-sm">{{.Content}}</p>
          <span class="text-xs text-gray-500">{{.Sender}} • {{.CreatedAt}}</span>
        </div>
      </div>
    {{else}}
      <!-- Right-aligned message -->
      <div class="flex justify-end">
        <div class="bg-blue-100 rounded-lg p-3 max-w-md">
          <p class="text-sm">{{.Content}}</p>
          <span class="text-xs text-gray-500">BOT • {{.CreatedAt}}</span>
        </div>
      </div>
    {{end}}
  {{end}}
</div>
```
- HTMX for "Load Older" button:
  - `hx-get="/api/conversations/{id}/messages?offset=50"`
  - `hx-swap="afterbegin"` to prepend messages
- Auto-scroll JavaScript:
```javascript
window.onload = () => {
  const thread = document.getElementById('message-thread');
  thread.scrollTop = thread.scrollHeight;
}
```
- Format timestamp: Use Go template function or JavaScript Date formatting

---

## Epic 3: WhatsApp Management Enhancement

**Goal:** Make WhatsApp pairing and testing effortless for tenant admins with real-time status updates, QR code display, and test message interface.

**Value:** Enable tenant admins to set up and verify WhatsApp bot functionality independently without technical support, reducing time-to-value.

**FRs Covered:** FR-3.1, FR-3.2, FR-3.3, FR-3.4

---

### Story 3.1: Enhance Pairing UI with QR Display and Flow

**As a** tenant admin,
**I want** an improved WhatsApp pairing interface with clear QR code display,
**So that** I can easily pair my WhatsApp number to the bot.

**Acceptance Criteria:**

**Given** I navigate to `/admin/whatsapp`
**And** WhatsApp is not currently connected
**When** the page loads
**Then** I see:
- Status indicator: "Disconnected" (red badge)
- Form with label: "Nomor WhatsApp"
- Phone number input field (placeholder: "+628123456789")
- "Mulai Pairing" button (primary blue)

**And** when I enter phone number and click "Mulai Pairing"
**Then** HTMX sends POST to `/admin/whatsapp/pair` with phone number
**And** form is replaced with QR code section showing:
  - Large QR code image (from `/admin/whatsapp/qr/{tenant_id}`)
  - Instructions: "Scan QR code ini dengan WhatsApp di HP Anda"
  - Steps:
    1. Buka WhatsApp > Menu > Linked Devices
    2. Tap "Link a Device"
    3. Scan QR code di layar ini
  - Status indicator changes to: "Pairing..." (yellow badge)
  - Cancel button: "Batalkan"

**And** QR code image refreshes every 30 seconds (HTMX polling)
**And** page polls status endpoint every 3 seconds (Story 3.2)
**And** if pairing timeout (2 minutes), show error: "Pairing timeout. Silakan coba lagi."
**And** "Cancel" button stops polling and returns to initial form

**Prerequisites:** WhatsApp handler endpoints already exist (internal/handler/whatsapp_handler.go)

**Technical Notes:**
- Template file: `templates/admin/whatsapp.html` (enhance existing)
- HTMX attributes:
  - Form: `hx-post="/admin/whatsapp/pair" hx-target="#pairing-section" hx-swap="outerHTML"`
  - QR refresh: `hx-get="/admin/whatsapp/qr/{tenant_id}" hx-trigger="every 30s" hx-swap="outerHTML"`
- Alpine.js for:
  - Form/QR section toggle
  - Timeout timer (2 min countdown)
- QR code styling: Large (300x300px), centered, white background, padding
- Instructions: Clear step-by-step in Bahasa Indonesia
- Cancel flow: Reset to initial state, stop all polling

---

### Story 3.2: Add Real-time Status Polling

**As a** tenant admin,
**I want** real-time WhatsApp connection status updates,
**So that** I know immediately when pairing succeeds or connection drops.

**Acceptance Criteria:**

**Given** I am on `/admin/whatsapp` page
**When** pairing is in progress (QR code shown)
**Then** HTMX polls GET `/admin/whatsapp/status` every 3 seconds
**And** on each poll, status indicator updates:
  - If `isConnected: true`: Change to "Connected" (green), hide QR, show connected view
  - If `isConnected: false` and `pairingStatus: "in_progress"`: Keep "Pairing..." (yellow)
  - If `pairingStatus: "failed"`: Show error message

**And** when status changes to "Connected"
**Then** polling stops
**And** QR code section is replaced with connected view showing:
  - Status: "Connected" (green badge)
  - Connected phone number (from status response)
  - "Disconnect" button (red, outlined)
  - Test message section (Story 3.3)
  - Last active time: "Terhubung sejak {timestamp}"

**And** when status is "Connected" on page load
**Then** slow polling continues (every 30 seconds) to detect disconnections
**And** if disconnection detected, status updates to "Disconnected" and UI resets

**And** polling stops when:
  - Pairing succeeds (connected)
  - User clicks "Cancel"
  - Pairing timeout (2 minutes)
  - User navigates away from page

**Prerequisites:** Story 3.1

**Technical Notes:**
- HTMX polling:
  - Fast polling during pairing: `hx-get="/admin/whatsapp/status" hx-trigger="every 3s" hx-target="#status-indicator"`
  - Slow polling when connected: `hx-trigger="every 30s"`
  - Stop polling: Remove HTMX attributes or set `hx-trigger="none"`
- Alpine.js for:
  - Polling state management (active/stopped)
  - Status transitions (disconnected → pairing → connected)
- Status API response:
```json
{
  "tenant_id": 1,
  "is_connected": true,
  "phone_number": "+628123456789",
  "pairing_status": "connected",
  "last_active": "2025-11-15T14:30:00Z"
}
```
- Visual indicators:
  - Disconnected: Red badge, `bg-red-100 text-red-800`
  - Pairing: Yellow badge, `bg-yellow-100 text-yellow-800`
  - Connected: Green badge, `bg-green-100 text-green-800`
- Handle edge case: Multiple tabs open (status sync across tabs via polling)

---

### Story 3.3: Build Test Message Interface

**As a** tenant admin,
**I want** to send test messages to the WhatsApp bot and see responses,
**So that** I can verify the bot is working correctly.

**Acceptance Criteria:**

**Given** WhatsApp status is "Connected"
**When** I view the WhatsApp management page
**Then** I see a "Test Message" section with:
  - Label: "Test WhatsApp Bot"
  - Textarea for message input (placeholder: "Ketik pesan test, misal: Ada Toyota Avanza?")
  - "Send Test Message" button (primary blue)
  - Response area (initially empty)

**And** when I type a test message and click "Send Test Message"
**Then** HTMX sends POST to `/admin/whatsapp/test` with message body
**And** button shows loading state: "Sending..." (disabled)
**And** after response (< 3 seconds), button returns to normal state

**And** response area displays:
  - Success message: "✅ Test message sent successfully!"
  - Details shown:
    - "Message sent to: {bot_phone_number}"
    - "Expected response: Check your WhatsApp for bot reply"
    - Timestamp of test
  - If response included in API: Show bot's actual response text

**And** if WhatsApp not connected
**Then** textarea and button are disabled
**And** show message: "⚠️ WhatsApp belum terhubung. Silakan pairing terlebih dahulu."

**And** if test fails
**Then** show error: "❌ Gagal mengirim test message: {error_details}"
**And** option to retry

**And** test message section is below connected status display
**And** layout is clean and user-friendly

**Prerequisites:** Story 3.2 (status must be connected)

**Technical Notes:**
- Part of `templates/admin/whatsapp.html`
- HTMX attributes:
  - `hx-post="/admin/whatsapp/test"`
  - `hx-target="#test-response"`
  - `hx-swap="innerHTML"`
  - `hx-indicator="#test-loading"` for loading state
- Alpine.js for:
  - Enable/disable based on connection status
  - Clear response area when sending new test
- Test endpoint (already exists): `POST /admin/whatsapp/test`
  - Request body: `{"message": "Ada Toyota Avanza?"}`
  - Response: `{"success": true, "message": "Test sent", "phone": "+628xxx"}`
- Response area styling:
  - Success: Green background, check icon
  - Error: Red background, X icon
  - Details: Monospace font for technical info
- Suggested test messages (shown as examples):
  - "Ada Toyota Avanza?"
  - "Mobil matic budget 150 juta?"
  - "Info detail mobil ID 5"

---

### Story 3.4: Add Disconnect Confirmation Modal

**As a** tenant admin,
**I want** to safely disconnect WhatsApp with confirmation,
**So that** I don't accidentally disconnect the bot.

**Acceptance Criteria:**

**Given** WhatsApp status is "Connected"
**When** I click the "Disconnect" button
**Then** a modal appears with:
  - Title: "Disconnect WhatsApp?"
  - Message: "Yakin ingin disconnect WhatsApp? Bot akan berhenti merespon pesan sampai Anda pairing ulang."
  - Warning icon (red)
  - "Ya, Disconnect" button (red, primary)
  - "Batalkan" button (grey, secondary)

**And** when I click "Ya, Disconnect"
**Then** HTMX sends POST to `/admin/whatsapp/disconnect`
**And** modal shows loading state
**And** on success:
  - Modal closes
  - Status updates to "Disconnected" (red badge)
  - QR section and test section are hidden
  - Pairing form is shown again
  - Success toast: "WhatsApp berhasil di-disconnect"

**And** when I click "Batalkan"
**Then** modal closes without action
**And** status remains "Connected"

**And** if disconnect fails
**Then** show error in modal: "Gagal disconnect: {error}"
**And** option to retry or close modal

**And** modal has:
  - Backdrop overlay (semi-transparent black)
  - Close on backdrop click or ESC key
  - Centered on screen
  - Mobile responsive

**Prerequisites:** Story 3.2

**Technical Notes:**
- Alpine.js for modal show/hide:
```javascript
<div x-data="{ showModal: false }">
  <button @click="showModal = true">Disconnect</button>
  <div x-show="showModal" @click.away="showModal = false">
    <!-- Modal content -->
  </div>
</div>
```
- HTMX on confirm button:
  - `hx-post="/admin/whatsapp/disconnect"`
  - `hx-target="#whatsapp-container"`
  - `hx-swap="outerHTML"`
- Disconnect endpoint (already exists): `POST /admin/whatsapp/disconnect`
  - Response: `{"success": true, "message": "Disconnected"}`
- Modal Tailwind styling:
  - Backdrop: `fixed inset-0 bg-black bg-opacity-50 z-40`
  - Modal: `fixed inset-0 flex items-center justify-center z-50`
  - Content: `bg-white rounded-lg p-6 max-w-md shadow-xl`
- Toast notification: Reuse existing toast system or Alpine component
- ESC key handler: Alpine.js `@keydown.escape.window="showModal = false"`

---

## Epic 4: Admin Dashboard Integration

**Goal:** Provide tenant admins with unified navigation and at-a-glance insights into bot activity, system status, and recent conversations.

**Value:** Single entry point for all admin tasks with quick visibility into system health and recent activity.

**FRs Covered:** FR-4.1, FR-4.2, FR-4.3

---

### Story 4.1: Update Navigation Menu with New Links

**As a** tenant admin,
**I want** updated sidebar navigation with all admin sections,
**So that** I can easily navigate between different management pages.

**Acceptance Criteria:**

**Given** I am logged in and on any admin page
**When** I view the left sidebar navigation
**Then** I see menu items in this order:
  1. 📊 Dashboard (existing)
  2. 🚗 Cars (existing)
  3. 👥 **Sales Team** (new)
  4. 💬 **Conversations** (new)
  5. 📱 WhatsApp (existing, possibly renamed to "WhatsApp Bot")

**And** each menu item shows:
  - Icon (emoji or Tailwind icon)
  - Label text
  - Active state indicator (highlighted when on that page)
  - Link to respective page

**And** when I click "Sales Team"
**Then** I navigate to `/admin/sales`
**And** "Sales Team" menu item is highlighted

**And** when I click "Conversations"
**Then** I navigate to `/admin/conversations`
**And** "Conversations" menu item is highlighted

**And** active state is:
  - Background: `bg-blue-100` (light blue)
  - Text: `text-blue-800` (dark blue)
  - Border left: `border-l-4 border-blue-600`

**And** sidebar is responsive:
  - Desktop: Always visible (fixed left)
  - Mobile: Collapsible hamburger menu

**Prerequisites:** Epic 1, Epic 2, Epic 3 (pages must exist)

**Technical Notes:**
- Template file: `templates/admin/layout.html`
- Update navigation section (~line 20-50)
- Menu structure:
```html
<nav class="w-64 bg-white border-r min-h-screen p-4">
  <ul class="space-y-2">
    <li>
      <a href="/admin/dashboard" class="{{if eq .CurrentPage "dashboard"}}bg-blue-100 text-blue-800 border-l-4 border-blue-600{{end}} flex items-center px-4 py-2 rounded">
        <span class="mr-3">📊</span>
        Dashboard
      </a>
    </li>
    <li>
      <a href="/admin/cars" class="{{if eq .CurrentPage "cars"}}bg-blue-100 text-blue-800 border-l-4 border-blue-600{{end}} flex items-center px-4 py-2 rounded">
        <span class="mr-3">🚗</span>
        Cars
      </a>
    </li>
    <li>
      <a href="/admin/sales" class="{{if eq .CurrentPage "sales"}}bg-blue-100 text-blue-800 border-l-4 border-blue-600{{end}} flex items-center px-4 py-2 rounded">
        <span class="mr-3">👥</span>
        Sales Team
      </a>
    </li>
    <li>
      <a href="/admin/conversations" class="{{if eq .CurrentPage "conversations"}}bg-blue-100 text-blue-800 border-l-4 border-blue-600{{end}} flex items-center px-4 py-2 rounded">
        <span class="mr-3">💬</span>
        Conversations
      </a>
    </li>
    <li>
      <a href="/admin/whatsapp" class="{{if eq .CurrentPage "whatsapp"}}bg-blue-100 text-blue-800 border-l-4 border-blue-600{{end}} flex items-center px-4 py-2 rounded">
        <span class="mr-3">📱</span>
        WhatsApp Bot
      </a>
    </li>
  </ul>
</nav>
```
- Each page template must pass `.CurrentPage` to layout
- Mobile hamburger: Alpine.js toggle
```javascript
<div x-data="{ mobileMenuOpen: false }">
  <button @click="mobileMenuOpen = !mobileMenuOpen" class="md:hidden">☰</button>
  <nav x-show="mobileMenuOpen" @click.away="mobileMenuOpen = false">
    <!-- Menu items -->
  </nav>
</div>
```

---

### Story 4.2: Add Quick Stats Cards to Dashboard

**As a** tenant admin,
**I want** to see at-a-glance statistics on the dashboard,
**So that** I can quickly understand system status and activity.

**Acceptance Criteria:**

**Given** I navigate to `/admin/dashboard`
**When** the page loads
**Then** I see 4 stat cards in a grid (2x2 on desktop, 1 column on mobile):

**Card 1: WhatsApp Status**
- Icon: 📱
- Title: "WhatsApp Bot"
- Value: "Connected" (green text) or "Disconnected" (red text)
- Subtitle: Phone number if connected, or "Not paired" if disconnected
- Link: "Manage" → `/admin/whatsapp`

**Card 2: Total Conversations**
- Icon: 💬
- Title: "Conversations"
- Value: "{count}" (e.g., "45")
- Subtitle: "Last 30 days"
- Link: "View All" → `/admin/conversations`

**Card 3: Sales Team**
- Icon: 👥
- Title: "Sales Team"
- Value: "{count} members" (e.g., "5 members")
- Subtitle: "Active sales team"
- Link: "Manage" → `/admin/sales`

**Card 4: Cars**
- Icon: 🚗
- Title: "Inventory"
- Value: "{count} cars" (e.g., "32 cars")
- Subtitle: "In catalog" (existing)
- Link: "Manage" → `/admin/cars`

**And** each card displays:
- White background with subtle shadow
- Rounded corners
- Padding and spacing for readability
- Hover effect (slight shadow increase)
- Click anywhere on card navigates to link

**And** stats are fetched from:
- WhatsApp: `GET /admin/whatsapp/status`
- Conversations: `GET /api/conversations/stats` (new endpoint)
- Sales: `GET /api/sales/stats` (new endpoint)
- Cars: Existing car count query

**And** if data fetch fails, show placeholder: "—" with retry option

**Prerequisites:** Story 1.3 (sales), Story 2.1 (conversations), Story 3.2 (whatsapp status)

**Technical Notes:**
- Template file: `templates/admin/dashboard.html` (update existing)
- Add stats section at top, before existing content
- Card grid layout:
```html
<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
  <!-- Card 1 -->
  <a href="/admin/whatsapp" class="bg-white rounded-lg shadow p-6 hover:shadow-lg transition">
    <div class="flex items-center mb-2">
      <span class="text-3xl mr-3">📱</span>
      <h3 class="text-gray-600 text-sm font-medium">WhatsApp Bot</h3>
    </div>
    <p class="text-2xl font-bold {{if .WhatsAppConnected}}text-green-600{{else}}text-red-600{{end}}">
      {{if .WhatsAppConnected}}Connected{{else}}Disconnected{{end}}
    </p>
    <p class="text-gray-500 text-sm mt-1">
      {{if .WhatsAppConnected}}{{.WhatsAppPhone}}{{else}}Not paired{{end}}
    </p>
  </a>
  <!-- Repeat for other cards -->
</div>
```
- Handler changes (`page_handler.go` or dashboard handler):
  - Fetch stats from repositories
  - Pass to template: `.WhatsAppConnected`, `.ConversationCount`, `.SalesCount`, `.CarCount`
- New repository methods (if not exist):
  - `conversationRepo.GetCount(ctx, tenantID, days int) (int, error)` (last N days)
  - `salesRepo.GetCount(ctx, tenantID) (int, error)`
- Cache stats for 5 minutes (optional optimization)

---

### Story 4.3: Build Recent Conversations Widget

**As a** tenant admin,
**I want** to see recent conversations on the dashboard,
**So that** I can quickly check latest customer activity.

**Acceptance Criteria:**

**Given** I am on `/admin/dashboard`
**When** the page loads (below stats cards)
**Then** I see "Recent Conversations" widget with:
  - Section title: "Recent Conversations"
  - List of 5 most recent conversations
  - "View All" link → `/admin/conversations`

**And** each conversation item displays:
  - Phone number (formatted)
  - Type badge (Customer/Sales) - small size
  - Last message preview (truncated to 40 chars)
  - Relative time (e.g., "5 menit lalu")
  - Click navigates to conversation detail

**And** conversation items are styled as:
  - White background
  - Border between items
  - Hover effect (background color change)
  - Compact spacing

**And** if no conversations exist
**Then** show empty state: "Belum ada conversation. Conversations akan muncul setelah customer atau sales mengirim pesan via WhatsApp."

**And** widget layout is:
  - Takes 60% width on desktop (with 40% for potential future widget)
  - Full width on mobile
  - Scrollable if content overflows

**Prerequisites:** Story 2.1 (conversation list endpoint)

**Technical Notes:**
- Part of `templates/admin/dashboard.html`
- Add widget below stats cards
- Widget structure:
```html
<div class="bg-white rounded-lg shadow p-6">
  <div class="flex justify-between items-center mb-4">
    <h2 class="text-lg font-semibold">Recent Conversations</h2>
    <a href="/admin/conversations" class="text-blue-600 text-sm hover:underline">View All →</a>
  </div>
  <div class="space-y-3">
    {{range .RecentConversations}}
    <a href="/admin/conversations/{{.ID}}" class="block p-3 hover:bg-gray-50 rounded border-b last:border-b-0">
      <div class="flex items-center justify-between">
        <div class="flex-1">
          <div class="flex items-center gap-2 mb-1">
            <span class="font-medium">{{.PhoneNumber}}</span>
            <span class="text-xs px-2 py-1 rounded {{if .IsSales}}bg-green-100 text-green-800{{else}}bg-blue-100 text-blue-800{{end}}">
              {{if .IsSales}}Sales{{else}}Customer{{end}}
            </span>
          </div>
          <p class="text-sm text-gray-600 truncate">{{.LastMessage}}</p>
        </div>
        <span class="text-xs text-gray-500 ml-4">{{.RelativeTime}}</span>
      </div>
    </a>
    {{end}}
  </div>
</div>
```
- Handler (`page_handler.go`):
  - Fetch: `conversationRepo.List(ctx, tenantID, 1, 5, "all")` (first 5)
  - Pass to template: `.RecentConversations`
- Relative time helper function (Go template or JavaScript):
  - < 1 hour: "{n} menit lalu"
  - < 24 hours: "{n} jam lalu"
  - < 7 days: "{n} hari lalu"
  - >= 7 days: "{date} {month}"

---

## Epic 5: Error Handling & Validation

**Goal:** Ensure robust user experience with comprehensive validation, error handling, and graceful degradation across all admin interfaces.

**Value:** Professional, production-ready admin panel that handles edge cases gracefully and provides clear feedback to users.

**FRs Covered:** FR-5.1, FR-5.2, FR-5.3

---

### Story 5.1: Implement Form Validation Framework

**As a** developer,
**I want** a consistent form validation approach across all admin forms,
**So that** users get immediate feedback on input errors before submission.

**Acceptance Criteria:**

**Given** any form in admin interface (sales, pairing, etc.)
**When** form is rendered
**Then** all required fields have:
  - HTML5 `required` attribute
  - Visual indicator (* asterisk or "Required" label)
  - Appropriate `type` attribute (tel, email, text, number)
  - `pattern` attribute for format validation (where applicable)

**And** client-side validation triggers:
  - On blur (when user leaves field)
  - On submit (before HTMX sends request)
  - Real-time for specific fields (e.g., phone format)

**And** validation error display:
  - Inline below/near the field
  - Red text color (`text-red-600`)
  - Icon indicator (❌ or ⚠️)
  - Clear message in Bahasa Indonesia

**And** form submit button is:
  - Disabled if validation fails (client-side)
  - Shows loading state during submission
  - Re-enabled after response (success or error)

**And** validation rules include:
  - **Phone number:** Regex `^\+?[1-9]\d{1,14}$` or `^08\d{8,11}$`
  - **Name:** Min 2 chars, max 100 chars, no special chars except spaces
  - **Required fields:** Non-empty, trimmed
  - **Max lengths:** Enforced via `maxlength` attribute

**And** server-side validation duplicates all client-side rules
**And** server errors override client validation messages

**Prerequisites:** None (applies to all forms)

**Technical Notes:**
- Alpine.js validation component:
```javascript
<div x-data="formValidation()">
  <form @submit.prevent="validateAndSubmit">
    <div class="mb-4">
      <label>Phone Number *</label>
      <input
        type="tel"
        x-model="form.phone"
        @blur="validatePhone"
        pattern="^\+?[1-9]\d{1,14}$|^08\d{8,11}$"
        required
        class="border rounded px-3 py-2 w-full"
        :class="{'border-red-500': errors.phone}"
      />
      <p x-show="errors.phone" x-text="errors.phone" class="text-red-600 text-sm mt-1"></p>
    </div>
  </form>
</div>

<script>
function formValidation() {
  return {
    form: { phone: '', name: '' },
    errors: {},
    validatePhone() {
      const phoneRegex = /^\+?[1-9]\d{1,14}$|^08\d{8,11}$/;
      if (!phoneRegex.test(this.form.phone)) {
        this.errors.phone = 'Format nomor tidak valid';
      } else {
        delete this.errors.phone;
      }
    },
    validateAndSubmit() {
      this.validatePhone();
      // Validate other fields...
      if (Object.keys(this.errors).length === 0) {
        // Allow HTMX to submit
      } else {
        // Prevent submission
        return false;
      }
    }
  }
}
</script>
```
- Error message templates (Bahasa Indonesia):
  - Phone: "Format nomor tidak valid (contoh: +628123456789 atau 08123456789)"
  - Required: "Field ini wajib diisi"
  - Min length: "Minimal {n} karakter"
  - Max length: "Maksimal {n} karakter"
  - Duplicate: "Data sudah ada"
- Apply to all forms: Sales add, WhatsApp pairing, any future forms
- Consider creating reusable validation Alpine component

---

### Story 5.2: Add API Error Handling Middleware and Response Format

**As a** developer,
**I want** consistent error handling across all API endpoints,
**So that** clients receive predictable error responses.

**Acceptance Criteria:**

**Given** any API endpoint is called
**When** an error occurs (validation, database, logic error)
**Then** the response includes:
  - Appropriate HTTP status code (400, 403, 404, 500)
  - JSON error object:
```json
{
  "error": "Nomor sudah terdaftar",
  "code": "DUPLICATE_PHONE",
  "details": {
    "field": "phone_number",
    "value": "+628123456789"
  }
}
```

**And** HTTP status codes are mapped correctly:
  - **400 Bad Request:** Validation errors, malformed input
  - **403 Forbidden:** Tenant isolation violation, unauthorized access
  - **404 Not Found:** Resource not found (but exists for different tenant)
  - **409 Conflict:** Duplicate resource (e.g., phone number already exists)
  - **500 Internal Server Error:** Unexpected errors, database failures

**And** all error messages are in Bahasa Indonesia
**And** error responses never expose sensitive info (stack traces, SQL queries)
**And** errors are logged server-side with full context (for debugging)

**And** validation errors include field-level details:
```json
{
  "error": "Validation failed",
  "code": "VALIDATION_ERROR",
  "details": {
    "fields": {
      "phone_number": "Format nomor tidak valid",
      "name": "Field ini wajib diisi"
    }
  }
}
```

**Prerequisites:** None (applies to all handlers)

**Technical Notes:**
- Create error handling middleware: `internal/middleware/error_handler.go`
```go
type ErrorResponse struct {
    Error   string                 `json:"error"`
    Code    string                 `json:"code,omitempty"`
    Details map[string]interface{} `json:"details,omitempty"`
}

func WriteError(w http.ResponseWriter, status int, message, code string, details map[string]interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteStatus(status)
    json.NewEncoder(w).Encode(ErrorResponse{
        Error:   message,
        Code:    code,
        Details: details,
    })
}

// Helper functions
func BadRequest(w http.ResponseWriter, message string) {
    WriteError(w, 400, message, "BAD_REQUEST", nil)
}

func NotFound(w http.ResponseWriter, message string) {
    WriteError(w, 404, message, "NOT_FOUND", nil)
}

func Forbidden(w http.ResponseWriter, message string) {
    WriteError(w, 403, message, "FORBIDDEN", nil)
}

func InternalError(w http.ResponseWriter, message string) {
    WriteError(w, 500, message, "INTERNAL_ERROR", nil)
}
```
- Use in handlers:
```go
if err := validatePhone(phone); err != nil {
    middleware.BadRequest(w, "Format nomor tidak valid")
    return
}

if salesExists {
    middleware.WriteError(w, 409, "Nomor sudah terdaftar", "DUPLICATE_PHONE", map[string]interface{}{
        "field": "phone_number",
        "value": phone,
    })
    return
}
```
- Error code constants: Define in `internal/model/errors.go`
- Logging: Use `pkg/logger` to log errors with context (tenant ID, request ID, user, timestamp)
- Recovery middleware: Catch panics and return 500 with generic message (don't expose panic details to client)

---

### Story 5.3: Add Network Error Recovery and Loading States

**As a** tenant admin,
**I want** clear feedback when operations are in progress and graceful error recovery,
**So that** I understand system state and can retry failed operations.

**Acceptance Criteria:**

**Given** any HTMX request is triggered (form submit, pagination, filter)
**When** request is in progress
**Then** I see a loading indicator:
  - Button: Shows "Loading..." or spinner, disabled state
  - Section: Shows skeleton loader or spinner overlay
  - Global: Optional loading bar at top of page

**And** when request succeeds
**Then** loading indicator is removed
**And** success feedback is shown (toast, inline message, or UI update)

**And** when request fails (network error, timeout, server error)
**Then** loading indicator is removed
**And** error message is displayed:
  - Inline near the action (for forms)
  - Toast notification (for global actions)
  - Error message in Bahasa Indonesia
  - "Retry" button or link

**And** HTMX timeout is set to 10 seconds (configurable)
**And** on timeout, show: "Koneksi timeout. Silakan coba lagi."

**And** if JavaScript is disabled
**Then** forms still work via standard POST (graceful degradation)
**And** full page reloads replace HTMX partial updates

**And** network error scenarios handled:
  - **Timeout:** Show retry option
  - **Offline:** Show "Tidak ada koneksi internet" message
  - **500 error:** Show "Terjadi kesalahan server" with retry
  - **403/404:** Show specific error from server

**Prerequisites:** All HTMX-enabled pages

**Technical Notes:**
- HTMX global config (add to layout.html):
```html
<script>
  document.body.addEventListener('htmx:configRequest', (event) => {
    event.detail.timeout = 10000; // 10 second timeout
  });

  document.body.addEventListener('htmx:timeout', (event) => {
    showToast('Koneksi timeout. Silakan coba lagi.', 'error');
  });

  document.body.addEventListener('htmx:responseError', (event) => {
    const status = event.detail.xhr.status;
    let message = 'Terjadi kesalahan';

    if (status === 0) {
      message = 'Tidak ada koneksi internet';
    } else if (status >= 500) {
      message = 'Terjadi kesalahan server';
    } else if (status === 403) {
      message = 'Akses ditolak';
    } else if (status === 404) {
      message = 'Data tidak ditemukan';
    }

    showToast(message, 'error');
  });
</script>
```
- Loading indicators with HTMX:
```html
<button
  hx-post="/api/sales"
  hx-indicator="#loading"
  class="btn-primary"
>
  <span class="htmx-indicator">Saving...</span>
  <span>Save</span>
</button>
<div id="loading" class="htmx-indicator">Loading...</div>
```
- CSS for htmx-indicator:
```css
.htmx-indicator {
  display: none;
}
.htmx-request .htmx-indicator {
  display: inline;
}
.htmx-request.htmx-indicator {
  display: inline;
}
```
- Toast notification component (Alpine.js):
```javascript
<div x-data="toast()" x-show="show" x-transition>
  <div :class="type === 'error' ? 'bg-red-500' : 'bg-green-500'"
       class="fixed top-4 right-4 px-6 py-3 rounded shadow-lg text-white">
    <p x-text="message"></p>
  </div>
</div>

<script>
function toast() {
  return {
    show: false,
    message: '',
    type: 'success',
    showToast(msg, t = 'success') {
      this.message = msg;
      this.type = t;
      this.show = true;
      setTimeout(() => { this.show = false }, 3000);
    }
  }
}

window.showToast = (msg, type) => {
  // Trigger Alpine toast
  Alpine.store('toast').showToast(msg, type);
}
</script>
```
- Graceful degradation: Ensure all forms have proper `action` and `method` attributes for non-JS fallback

---

## Epic 6: Sales Car Upload via WhatsApp

**Goal:** Empower sales team to add cars to catalog directly via WhatsApp chat with AI-powered text parsing, photo upload, and role-based access control.

**Value:** Streamline car inventory management by allowing sales to upload directly from field (photo at dealership → instant catalog entry), reducing admin overhead and time-to-market.

**FRs Covered:** FR-6.1, FR-6.2, FR-6.3, FR-6.4, FR-6.5, FR-6.6, FR-6.7

---

### Story 6.1: Enhance Bot with Role Detection Logic

**As a** bot developer,
**I want** the LLM bot to detect if sender is sales or customer,
**So that** different capabilities are available based on role.

**Acceptance Criteria:**

**Given** a WhatsApp message is received
**When** bot processes the message in `ProcessMessage` method
**Then** bot checks: `isSales := salesRepo.IsSales(tenantID, senderPhone)`
**And** bot stores role in conversation context

**And** bot's system prompt changes based on role:
  - **Customer:** "Anda adalah asisten penjualan mobil yang ramah. Bantu customer mencari mobil yang sesuai..."
  - **Sales:** "Anda adalah asisten untuk sales team. Anda dapat membantu upload mobil baru ke catalog selain menjawab pertanyaan tentang inventory..."

**And** available functions change based on role:
  - **Customer:** `searchCars`, `getCarDetails`, `sendCarImages`
  - **Sales:** All customer functions + `uploadCar`

**And** if customer tries to upload car (detected from keywords like "upload", "tambah mobil")
**Then** bot responds politely: "Maaf, fitur upload mobil hanya tersedia untuk sales team. Silakan hubungi sales kami untuk menambahkan mobil."

**And** if sales uploads car
**Then** bot processes upload workflow (Stories 6.2-6.6)

**Prerequisites:** Sales repository with IsSales method (already exists)

**Technical Notes:**
- File to modify: `internal/llm/bot.go`
- Update `ProcessMessage` method signature to accept `isSales bool` (already has this parameter)
- Update `buildSystemPrompt` method:
```go
func (b *Bot) buildSystemPrompt(isSales bool) string {
    basePrompt := `Anda adalah asisten penjualan mobil yang ramah dan profesional di Auto LMK.`

    if isSales {
        return basePrompt + `

Anda sedang membantu SALES TEAM. Anda memiliki kemampuan tambahan:
- Upload mobil baru ke catalog
- Lihat dan kelola inventory

Saat sales ingin upload mobil, minta mereka:
1. Upload foto mobil (bisa multiple)
2. Ketik detail: "Brand Model Tahun Harga Transmisi BahanBakar"
   Contoh: "Toyota Avanza 2020 185juta AT Bensin"
        `
    }

    return basePrompt + `

Anda membantu CUSTOMER mencari mobil. Anda dapat:
- Mencari mobil berdasarkan brand, budget, transmisi, dll
- Menampilkan detail dan foto mobil
- Memberikan rekomendasi

Jika customer tanya tentang upload atau tambah mobil, jelaskan bahwa fitur itu untuk sales team.
    `
}
```
- Update `GetAvailableFunctions`:
```go
func (b *Bot) GetAvailableFunctions(isSales bool) []Function {
    baseFunctions := []Function{
        {Name: "searchCars", ...},
        {Name: "getCarDetails", ...},
        {Name: "sendCarImages", ...},
    }

    if isSales {
        baseFunctions = append(baseFunctions, Function{
            Name: "uploadCar",
            Description: "Upload mobil baru ke catalog dengan foto dan detail",
            Parameters: map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "brand": map[string]string{"type": "string", "description": "Brand mobil (Toyota, Honda, dll)"},
                    "model": map[string]string{"type": "string", "description": "Model mobil (Avanza, Civic, dll)"},
                    "year": map[string]string{"type": "integer", "description": "Tahun produksi"},
                    "price": map[string]string{"type": "integer", "description": "Harga dalam rupiah"},
                    "transmission": map[string]string{"type": "string", "description": "AT atau MT"},
                    "fuel_type": map[string]string{"type": "string", "description": "Bensin atau Diesel"},
                    "description": map[string]string{"type": "string", "description": "Deskripsi tambahan (optional)"},
                },
                "required": []string{"brand", "model", "year", "price", "transmission", "fuel_type"},
            },
        })
    }

    return baseFunctions
}
```
- Keyword detection for rejection (if not using function calling):
  - Keywords: "upload", "tambah mobil", "post mobil", "masukin mobil"
  - Response template: Store in bot configuration

---

### Story 6.2: Add uploadCar Function to LLM Bot

**As a** bot developer,
**I want** to implement uploadCar function that LLM can call,
**So that** sales can upload car data via conversational AI.

**Acceptance Criteria:**

**Given** sales sends message with car details
**And** LLM decides to call `uploadCar` function
**When** bot receives function call
**Then** bot executes `executeUploadCar` method with parsed parameters:
  - brand (string, required)
  - model (string, required)
  - year (int, required)
  - price (int, required, in rupiah)
  - transmission (string, required: "AT" or "MT")
  - fuel_type (string, required: "Bensin" or "Diesel")
  - description (string, optional)
  - photo_urls ([]string, from previous message context)

**And** bot validates all required fields are present
**And** bot checks for uploaded photos in pending uploads (from Story 6.3)
**And** if photos missing, bot responds: "Silakan upload foto mobil terlebih dahulu"

**And** bot calls car repository to save:
  1. Create car record in `cars` table
  2. Save photos to `car_photos` table
  3. Save additional specs to `car_specs` table (if needed)

**And** on success, bot returns:
```json
{
  "success": true,
  "car_id": 123,
  "message": "Mobil berhasil ditambahkan ke catalog!",
  "catalog_url": "https://tenant.auto-lmk.com/cars/123"
}
```

**And** on failure, bot returns error:
```json
{
  "success": false,
  "error": "Harga harus berupa angka"
}
```

**And** LLM uses function result to generate confirmation message to sales

**Prerequisites:** Story 6.1 (role detection), Story 6.3 (photo upload handling)

**Technical Notes:**
- Add to `internal/llm/bot.go`:
```go
func (b *Bot) executeUploadCar(ctx context.Context, tenantID int, args map[string]interface{}) (interface{}, error) {
    // Parse arguments
    brand := args["brand"].(string)
    model := args["model"].(string)
    year := int(args["year"].(float64))
    price := int(args["price"].(float64))
    transmission := args["transmission"].(string)
    fuelType := args["fuel_type"].(string)
    description := ""
    if desc, ok := args["description"]; ok {
        description = desc.(string)
    }

    // Validate
    if transmission != "AT" && transmission != "MT" {
        return nil, errors.New("Transmisi harus AT atau MT")
    }
    if fuelType != "Bensin" && fuelType != "Diesel" {
        return nil, errors.New("Bahan bakar harus Bensin atau Diesel")
    }

    // Get pending photos (from context storage)
    photos := b.GetPendingPhotos(senderPhone) // New method to implement
    if len(photos) == 0 {
        return nil, errors.New("Silakan upload foto mobil terlebih dahulu")
    }

    // Save to database via car repository
    car := &model.Car{
        TenantID:     tenantID,
        Brand:        brand,
        Model:        model,
        Year:         year,
        Price:        price,
        Transmission: transmission,
        FuelType:     fuelType,
        Description:  description,
    }

    carID, err := b.carRepo.Create(ctx, car)
    if err != nil {
        return nil, fmt.Errorf("Gagal menyimpan mobil: %v", err)
    }

    // Save photos
    for _, photoPath := range photos {
        b.carRepo.AddPhoto(ctx, carID, photoPath)
    }

    // Clear pending photos
    b.ClearPendingPhotos(senderPhone)

    // Return success
    return map[string]interface{}{
        "success": true,
        "car_id": carID,
        "message": "Mobil berhasil ditambahkan ke catalog!",
        "catalog_url": fmt.Sprintf("https://%s.auto-lmk.com/cars/%d", getTenantDomain(tenantID), carID),
    }, nil
}
```
- Update `executeFunction` switch statement to handle "uploadCar"
- Add photo context storage:
  - In-memory map: `pendingPhotos map[string][]string` (phone → photo URLs)
  - Or use conversation metadata table
  - Photos expire after 10 minutes if not used
- Car repository method: `Create(ctx, car)` already exists
- Photo repository method: `AddPhoto(ctx, carID, photoPath)` may need to be added

---

### Story 6.3: Implement Photo Upload Handler for WhatsApp

**As a** bot developer,
**I want** to handle photo uploads from WhatsApp messages,
**So that** sales can send car photos before submitting details.

**Acceptance Criteria:**

**Given** sales sends a WhatsApp message with photo attachment(s)
**When** WhatsAppService receives the message event
**Then** service detects message type is "image"
**And** service downloads image from WhatsApp servers

**And** service saves image to local storage:
  - Path: `/static/uploads/cars/{tenant_id}/{timestamp}_{random}.jpg`
  - Generate thumbnail: 300x300px for listing view
  - Store original: Full resolution for detail view

**And** service stores photo URL in pending uploads context:
  - Key: sender phone number
  - Value: array of photo URLs
  - Expires: 10 minutes (cleared after successful upload or timeout)

**And** service responds to sender:
  - If first photo: "Foto diterima! Upload lebih banyak foto atau ketik detail mobil."
  - If additional photo: "Foto {n} diterima! Total {total} foto."
  - If 5+ photos: "Maksimum 5 foto. Silakan ketik detail mobil sekarang."

**And** if customer (non-sales) tries to upload photo
**Then** bot responds: "Terima kasih! Untuk upload mobil, silakan hubungi sales team kami."
**And** photo is not saved

**And** supported formats: JPG, PNG (max 5MB each)
**And** if format invalid, respond: "Format tidak didukung. Gunakan JPG atau PNG."

**Prerequisites:** Story 6.1 (role detection)

**Technical Notes:**
- Update `internal/service/whatsapp_service.go`:
```go
func (s *WhatsAppService) ProcessIncomingMessage(ctx context.Context, tenantID int, senderPhone, messageText string, messageType string, mediaURL string) error {
    // Existing role check
    isSales, _ := s.salesRepo.IsSales(tenantID, senderPhone)

    // Handle image upload
    if messageType == "image" {
        if !isSales {
            s.waClient.SendMessage(tenantID, senderPhone, "Terima kasih! Untuk upload mobil, silakan hubungi sales team kami.")
            return nil
        }

        // Download image from WhatsApp
        imageData, err := s.waClient.DownloadMedia(mediaURL)
        if err != nil {
            return err
        }

        // Validate format and size
        if !isValidImage(imageData) {
            s.waClient.SendMessage(tenantID, senderPhone, "Format tidak didukung. Gunakan JPG atau PNG (max 5MB).")
            return nil
        }

        // Save to disk
        photoPath, err := saveCarPhoto(tenantID, imageData)
        if err != nil {
            return err
        }

        // Store in pending context
        count := s.bot.AddPendingPhoto(senderPhone, photoPath)

        // Send confirmation
        var response string
        if count == 1 {
            response = "Foto diterima! Upload lebih banyak foto atau ketik detail mobil.\nContoh: Toyota Avanza 2020 185juta AT Bensin"
        } else if count < 5 {
            response = fmt.Sprintf("Foto %d diterima! Total %d foto.", count, count)
        } else {
            response = "Maksimum 5 foto tercapai. Silakan ketik detail mobil sekarang."
        }
        s.waClient.SendMessage(tenantID, senderPhone, response)

        return nil
    }

    // Continue with existing text message processing...
}
```
- Add to `internal/llm/bot.go`:
```go
type Bot struct {
    // ... existing fields
    pendingPhotos map[string][]string // phone → photo URLs
    photosMutex   sync.RWMutex
}

func (b *Bot) AddPendingPhoto(phone, photoURL string) int {
    b.photosMutex.Lock()
    defer b.photosMutex.Unlock()

    if b.pendingPhotos == nil {
        b.pendingPhotos = make(map[string][]string)
    }

    if len(b.pendingPhotos[phone]) >= 5 {
        return len(b.pendingPhotos[phone])
    }

    b.pendingPhotos[phone] = append(b.pendingPhotos[phone], photoURL)

    // Set expiration timer (10 minutes)
    go func() {
        time.Sleep(10 * time.Minute)
        b.ClearPendingPhotos(phone)
    }()

    return len(b.pendingPhotos[phone])
}

func (b *Bot) GetPendingPhotos(phone string) []string {
    b.photosMutex.RLock()
    defer b.photosMutex.RUnlock()
    return b.pendingPhotos[phone]
}

func (b *Bot) ClearPendingPhotos(phone string) {
    b.photosMutex.Lock()
    defer b.photosMutex.Unlock()
    delete(b.pendingPhotos, phone)
}
```
- Image processing helpers:
```go
func isValidImage(data []byte) bool {
    // Check magic bytes for JPG/PNG
    // Check size < 5MB
    return true
}

func saveCarPhoto(tenantID int, imageData []byte) (string, error) {
    // Create directory if not exists
    dir := fmt.Sprintf("static/uploads/cars/%d", tenantID)
    os.MkdirAll(dir, 0755)

    // Generate filename
    filename := fmt.Sprintf("%d_%s.jpg", time.Now().Unix(), generateRandomString(8))
    path := filepath.Join(dir, filename)

    // Save original
    ioutil.WriteFile(path, imageData, 0644)

    // Generate thumbnail (optional for MVP, can be done later)
    // createThumbnail(path, 300, 300)

    return "/"+path, nil
}
```
- WhatsApp media download: Use whatsmeow library's DownloadMediaMessage method

---

### Story 6.4: Implement Text Parser for Car Details

**As a** bot developer,
**I want** LLM to intelligently parse car details from natural language text,
**So that** sales can submit information flexibly without strict format.

**Acceptance Criteria:**

**Given** sales sends text message with car details
**When** message contains keywords related to car upload (brand names, price indicators, "upload", etc.)
**And** sales has pending photos uploaded
**Then** LLM processes message and extracts structured data via function calling

**And** LLM can parse various formats:
  - Structured: "Toyota Avanza 2020 185juta AT Bensin"
  - Natural: "Avanza tahun 2020, harga 185 juta, matic, bensin"
  - Conversational: "Ini Avanza 2020, dijual 185jt, transmisi matic"

**And** LLM normalizes values:
  - Price: "185juta" → 185000000, "185jt" → 185000000, "185 juta" → 185000000
  - Transmission: "matic" → "AT", "manual" → "MT", "otomatis" → "AT"
  - Fuel: "bensin" → "Bensin", "diesel" → "Diesel", "solar" → "Diesel"
  - Year: Validates 1990-2025 range

**And** if required fields missing, LLM asks follow-up:
  - Missing brand/model: "Brand dan model apa?"
  - Missing price: "Berapa harganya?"
  - Missing year: "Tahun berapa?"
  - Missing transmission: "Transmisi matic atau manual?"

**And** LLM confirms before calling uploadCar function:
  - "Konfirmasi: Toyota Avanza 2020, Rp 185.000.000, AT, Bensin. Benar?"
  - Sales responds "Ya" or "Benar" → Call uploadCar
  - Sales responds "Tidak" or corrects → Re-parse

**Prerequisites:** Story 6.1 (role detection), Story 6.2 (uploadCar function), Story 6.3 (photo upload)

**Technical Notes:**
- LLM parsing happens automatically via function calling (no explicit parser needed)
- System prompt enhancement (in `buildSystemPrompt`):
```
Saat sales ingin upload mobil:
1. Cek apakah ada foto yang sudah di-upload (gunakan context)
2. Extract informasi dari pesan:
   - Brand: Toyota, Honda, Daihatsu, Suzuki, Mitsubishi, Nissan, dll
   - Model: Avanza, Civic, Xenia, Ertiga, Pajero, dll
   - Tahun: 1990-2025
   - Harga: Parse dari "185juta", "185jt", "185 juta" → 185000000
   - Transmisi: "matic"/"otomatis"/"AT" → AT, "manual"/"MT" → MT
   - Bahan Bakar: "bensin" → Bensin, "diesel"/"solar" → Diesel
3. Jika info lengkap, konfirmasi dulu sebelum call uploadCar
4. Jika info kurang, tanya field yang missing

Contoh parsing:
- "Toyota Avanza 2020 185juta AT Bensin" → brand:Toyota, model:Avanza, year:2020, price:185000000, transmission:AT, fuel:Bensin
- "Avanza 2020 matic 185jt bensin" → (sama, extract brand dari model knowledge)
```
- Price parsing regex (if needed in validation):
  - `/(\d+)\s*(juta|jt|jutaan)/i` → multiply by 1,000,000
  - `/(\d+)\s*(jt)/i` → multiply by 1,000,000
  - `/(\d{3,})/` → use as-is if > 100k
- Validation after LLM extraction:
  - Year: 1990 <= year <= 2025
  - Price: 10,000,000 <= price <= 10,000,000,000 (10 juta - 10 miliar)
  - Transmission: "AT" or "MT"
  - Fuel: "Bensin" or "Diesel"
- Confirmation flow:
  - LLM generates confirmation message
  - Store parsed data in conversation context
  - Wait for "ya"/"benar"/"betul" response
  - On confirmation, call uploadCar with stored data
- Alternative: If LLM directly calls uploadCar, skip confirmation (simpler, trust LLM)

---

### Story 6.5: Integrate Auto-save to Database Logic

**As a** bot developer,
**I want** uploadCar function to save complete car data to database,
**So that** uploaded cars immediately appear in catalog.

**Acceptance Criteria:**

**Given** uploadCar function is called with valid car data and photo URLs
**When** function executes database save
**Then** the following database operations happen in a transaction:

**1. Create car record:**
```sql
INSERT INTO cars (tenant_id, brand, model, year, price, transmission, fuel_type, description, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
RETURNING id;
```

**2. Save car photos:**
```sql
INSERT INTO car_photos (car_id, photo_url, is_primary, created_at)
VALUES
  ($1, $2, true, NOW()),  -- First photo is primary
  ($1, $3, false, NOW()), -- Subsequent photos
  ...;
```

**3. Save car specs (optional, if additional details provided):**
```sql
INSERT INTO car_specs (car_id, spec_key, spec_value)
VALUES
  ($1, 'color', $2),
  ($1, 'mileage', $3),
  ...;
```

**And** if any database operation fails
**Then** entire transaction is rolled back
**And** error is returned to LLM: "Gagal menyimpan mobil: {error}"

**And** if all operations succeed
**Then** transaction is committed
**And** car is immediately visible on:
  - Public catalog: `https://tenant.auto-lmk.com/cars`
  - Admin car list: `/admin/cars`
  - WhatsApp bot search results

**And** success response includes:
  - car_id (for reference)
  - catalog_url (for sharing with customers)

**Prerequisites:** Story 6.2 (uploadCar function)

**Technical Notes:**
- Already covered in Story 6.2, but emphasize transaction handling:
```go
func (r *CarRepository) CreateWithPhotos(ctx context.Context, car *model.Car, photoURLs []string) (int, error) {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return 0, err
    }
    defer tx.Rollback()

    // Insert car
    var carID int
    err = tx.QueryRowContext(ctx, `
        INSERT INTO cars (tenant_id, brand, model, year, price, transmission, fuel_type, description, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
        RETURNING id
    `, car.TenantID, car.Brand, car.Model, car.Year, car.Price, car.Transmission, car.FuelType, car.Description).Scan(&carID)
    if err != nil {
        return 0, err
    }

    // Insert photos
    for i, photoURL := range photoURLs {
        isPrimary := i == 0
        _, err = tx.ExecContext(ctx, `
            INSERT INTO car_photos (car_id, photo_url, is_primary, created_at)
            VALUES ($1, $2, $3, NOW())
        `, carID, photoURL, isPrimary)
        if err != nil {
            return 0, err
        }
    }

    // Commit transaction
    if err = tx.Commit(); err != nil {
        return 0, err
    }

    return carID, nil
}
```
- Use in `executeUploadCar`:
```go
carID, err := b.carRepo.CreateWithPhotos(ctx, car, photos)
if err != nil {
    return nil, fmt.Errorf("Gagal menyimpan mobil: %v", err)
}
```
- Update car repository interface if `CreateWithPhotos` doesn't exist
- Consider adding audit log: Who uploaded (sales phone), when, from WhatsApp

---

### Story 6.6: Build Confirmation Message Generator

**As a** sales team member,
**I want** clear confirmation when car upload succeeds,
**So that** I know the car is live and can share catalog link with customers.

**Acceptance Criteria:**

**Given** uploadCar function successfully saves car to database
**When** LLM receives success response from function
**Then** LLM generates confirmation message to sales including:
  - Success indicator: "✅ Mobil berhasil ditambahkan ke catalog!"
  - Car summary: "{Brand} {Model} {Year} - Rp {FormattedPrice}"
  - Car ID: "ID Mobil: #{car_id}"
  - Catalog link: "Lihat di catalog: {catalog_url}"
  - Photo count: "{n} foto tersimpan"
  - Call to action: "Silakan cek catalog atau share link ke customer!"

**And** message is formatted nicely:
```
✅ *Mobil berhasil ditambahkan!*

📋 Detail:
• Brand: Toyota Avanza 2020
• Harga: Rp 185.000.000
• Transmisi: Matic (AT)
• Bahan Bakar: Bensin
• Foto: 3 foto tersimpan

🆔 ID Mobil: #123

🔗 Lihat di catalog:
https://tenant.auto-lmk.com/cars/123

Silakan cek catalog atau share link ke customer!
```

**And** catalog URL is clickable in WhatsApp (valid HTTPS URL)
**And** price is formatted with thousand separators: "Rp 185.000.000"

**And** if upload fails
**Then** LLM generates error message:
```
❌ Gagal menambahkan mobil

Error: {error_message}

Silakan coba lagi atau hubungi admin.
```

**Prerequisites:** Story 6.5 (auto-save logic)

**Technical Notes:**
- LLM handles message generation automatically based on function response
- Enhance system prompt with formatting instructions:
```
Saat uploadCar berhasil, buat pesan konfirmasi dengan format:
✅ *Mobil berhasil ditambahkan!*

📋 Detail:
• Brand: {brand} {model} {year}
• Harga: Rp {price_formatted}
• Transmisi: {transmission_display}
• Bahan Bakar: {fuel_type}
• Foto: {photo_count} foto

🆔 ID Mobil: #{car_id}

🔗 Lihat di catalog:
{catalog_url}

Silakan cek catalog atau share link ke customer!

Format harga dengan titik separator (185.000.000).
Transmisi: AT → "Matic (AT)", MT → "Manual (MT)".
```
- Helper function for price formatting (if needed in function response):
```go
func formatPrice(price int) string {
    // 185000000 → "185.000.000"
    str := strconv.Itoa(price)
    var result []rune
    for i, r := range str {
        if i > 0 && (len(str)-i)%3 == 0 {
            result = append(result, '.')
        }
        result = append(result, r)
    }
    return string(result)
}
```
- Return formatted price in function response:
```go
return map[string]interface{}{
    "success": true,
    "car_id": carID,
    "brand": brand,
    "model": model,
    "year": year,
    "price_formatted": formatPrice(price),
    "transmission_display": transmissionDisplay(transmission),
    "fuel_type": fuelType,
    "photo_count": len(photos),
    "catalog_url": catalogURL,
}, nil
```
- WhatsApp formatting: Use *bold*, _italic_, ~strikethrough~ (if supported)

---

### Story 6.7: Add Customer Upload Rejection Logic

**As a** bot developer,
**I want** to politely reject car upload attempts from customers,
**So that** catalog integrity is maintained and only sales can add cars.

**Acceptance Criteria:**

**Given** customer (non-sales) sends message that looks like car upload attempt
**When** bot detects upload-related keywords OR customer uploads photo
**Then** bot responds with polite rejection message

**And** rejection scenarios:
  1. **Customer uploads photo:**
     - Response: "Terima kasih! Untuk upload mobil, silakan hubungi sales team kami. Saya dapat membantu Anda mencari mobil yang tersedia. Ada yang bisa saya bantu?"

  2. **Customer uses upload keywords ("upload mobil", "tambah mobil", "posting mobil"):**
     - Response: "Maaf, fitur upload mobil hanya tersedia untuk sales team. Apakah Anda ingin mencari mobil? Saya bisa bantu!"

  3. **Customer asks about uploading:**
     - Customer: "Gimana cara upload mobil?"
     - Response: "Untuk menambahkan mobil ke catalog, Anda perlu menjadi bagian dari sales team kami. Silakan hubungi admin untuk informasi lebih lanjut. \n\nApakah ada yang bisa saya bantu untuk pencarian mobil?"

**And** rejection messages are:
  - Polite and professional
  - Redirect to available features (search cars)
  - Offer alternative help
  - In Bahasa Indonesia

**And** customer photos are not saved to disk
**And** customer upload attempts are logged (for audit/analytics)

**And** after rejection, bot continues normal customer service flow

**Prerequisites:** Story 6.1 (role detection)

**Technical Notes:**
- Already partially implemented in Story 6.1 and Story 6.3
- Enhance keyword detection (if not using LLM for intent):
```go
var uploadKeywords = []string{
    "upload mobil", "tambah mobil", "posting mobil", "masukin mobil",
    "input mobil", "daftar mobil", "cara upload", "gimana upload",
}

func detectUploadIntent(message string) bool {
    lowerMsg := strings.ToLower(message)
    for _, keyword := range uploadKeywords {
        if strings.Contains(lowerMsg, keyword) {
            return true
        }
    }
    return false
}
```
- In `ProcessIncomingMessage` (whatsapp_service.go):
```go
if !isSales && (messageType == "image" || detectUploadIntent(messageText)) {
    response := "Terima kasih! Untuk upload mobil, silakan hubungi sales team kami. Saya dapat membantu Anda mencari mobil yang tersedia. Ada yang bisa saya bantu?"
    s.waClient.SendMessage(tenantID, senderPhone, response)

    // Log attempt
    s.logger.Info("Customer upload attempt blocked",
        "phone", senderPhone,
        "tenant_id", tenantID,
        "message_type", messageType)

    return nil
}
```
- LLM system prompt (customer version) already handles this:
```
Jika customer bertanya tentang upload atau tambah mobil, jelaskan bahwa:
- Fitur upload hanya untuk sales team
- Mereka bisa hubungi admin atau sales untuk info lebih lanjut
- Tawarkan bantuan mencari mobil yang tersedia
```
- Response templates stored in configuration or system prompt
- Logging for analytics: Track how many customers ask about uploading (might indicate interest in becoming sellers/partners)

---

## FR Coverage Matrix

| FR Code | FR Description | Epic | Stories |
|---------|----------------|------|---------|
| FR-1.1 | Add Sales Team Member | Epic 1 | 1.1, 1.2, 1.3 |
| FR-1.2 | List Sales Team Members | Epic 1 | 1.1, 1.2, 1.3 |
| FR-1.3 | Delete Sales Team Member | Epic 1 | 1.1, 1.2, 1.3 |
| FR-2.1 | List Conversations | Epic 2 | 2.1, 2.3 |
| FR-2.2 | View Conversation Detail | Epic 2 | 2.2, 2.4 |
| FR-2.3 | Filter Conversations | Epic 2 | 2.3 |
| FR-2.4 | Search Conversation by Phone | Epic 2 | 2.3 |
| FR-3.1 | Enhanced Pairing Interface | Epic 3 | 3.1 |
| FR-3.2 | Connection Status Display | Epic 3 | 3.2 |
| FR-3.3 | Test Message Interface | Epic 3 | 3.3 |
| FR-3.4 | Disconnect WhatsApp | Epic 3 | 3.4 |
| FR-4.1 | Navigation Menu Update | Epic 4 | 4.1 |
| FR-4.2 | Dashboard Quick Stats | Epic 4 | 4.2 |
| FR-4.3 | Recent Conversations Widget | Epic 4 | 4.3 |
| FR-5.1 | Form Validation | Epic 5 | 5.1 |
| FR-5.2 | API Error Handling | Epic 5 | 5.2 |
| FR-5.3 | Network Error Handling | Epic 5 | 5.3 |
| FR-6.1 | Role Detection | Epic 6 | 6.1 |
| FR-6.2 | Upload Car Photos (Sales Only) | Epic 6 | 6.3 |
| FR-6.3 | Submit Car Details via Text | Epic 6 | 6.4 |
| FR-6.4 | AI Parse Car Details | Epic 6 | 6.4 |
| FR-6.5 | Auto-save to Database | Epic 6 | 6.2, 6.5 |
| FR-6.6 | Confirmation & Catalog Link | Epic 6 | 6.6 |
| FR-6.7 | Reject Customer Upload | Epic 6 | 6.7 |

**Coverage Status:** ✅ All 24 FRs mapped to stories

---

## Summary

### Epic Breakdown Complete

**Total Epics:** 6
**Total Stories:** 28
**Total FRs Covered:** 24/24 (100%)

### Epic Summary

| Epic | Stories | Estimated Effort | Priority | Dependencies |
|------|---------|-----------------|----------|--------------|
| **Epic 1: Sales Team Management** | 3 stories | 1-2 hari | HIGH (Foundation) | None |
| **Epic 2: Conversation Monitoring** | 4 stories | 2-3 hari | HIGH | Epic 1 (for role context) |
| **Epic 3: WhatsApp Management Enhancement** | 4 stories | 1-2 hari | HIGH | None (parallel dengan Epic 1-2) |
| **Epic 4: Admin Dashboard Integration** | 3 stories | 1 hari | MEDIUM | Epic 1, 2, 3 (integration) |
| **Epic 5: Error Handling & Validation** | 3 stories | 1-2 hari | MEDIUM (incremental) | Applied across all epics |
| **Epic 6: Sales Car Upload via WhatsApp** | 7 stories | 2-3 hari | HIGH | Epic 1 (needs sales data) |

**Total Estimated Effort:** 8-13 hari development (with 1-2 developers)

### Implementation Sequencing

**Week 1: Core Admin Interface & WhatsApp Bot Enhancement**

**Day 1-2:** Epic 1 (Sales Team Management)
- Story 1.1: Sales Handler (3-4 jam)
- Story 1.2: Sales UI (3-4 jam)
- Story 1.3: Route activation & testing (2-3 jam)
- **Deliverable:** Admin dapat manage sales team

**Day 2-3:** Epic 6 (Sales Car Upload) - Start setelah Epic 1 selesai
- Story 6.1: Role detection enhancement (2-3 jam)
- Story 6.2: uploadCar function (3-4 jam)
- Story 6.3: Photo upload handler (3-4 jam)
- Story 6.4: Text parser (2-3 jam)
- Story 6.5: Database transaction (2 jam)
- Story 6.6: Confirmation messages (2 jam)
- Story 6.7: Customer rejection (1 jam)
- **Deliverable:** Sales dapat upload mobil via WhatsApp

**Day 3-4:** Epic 2 (Conversation Monitoring) - Parallel dengan Epic 6
- Story 2.1: List handler (3-4 jam)
- Story 2.2: Detail handler (2-3 jam)
- Story 2.3: List UI dengan filters (3-4 jam)
- Story 2.4: Detail UI dengan message thread (3-4 jam)
- **Deliverable:** Admin dapat monitor conversations

**Day 4-5:** Epic 3 (WhatsApp Management) - Parallel dengan Epic 2
- Story 3.1: Enhanced pairing UI (3-4 jam)
- Story 3.2: Status polling (2-3 jam)
- Story 3.3: Test message interface (2 jam)
- Story 3.4: Disconnect confirmation (1-2 jam)
- **Deliverable:** Pairing dan testing workflow smooth

**Day 5:** Epic 4 (Dashboard Integration)
- Story 4.1: Navigation menu (1-2 jam)
- Story 4.2: Quick stats cards (2-3 jam)
- Story 4.3: Recent conversations widget (2 jam)
- **Deliverable:** Unified admin dashboard

**Day 1-6 (Incremental):** Epic 5 (Error Handling)
- Story 5.1: Form validation framework (2-3 jam, apply to all forms)
- Story 5.2: API error handling (2-3 jam, add middleware)
- Story 5.3: Network error recovery (2-3 jam, HTMX config)
- **Deliverable:** Robust error handling across all features

**Day 6-7:** End-to-End Testing & Refinement
- Multi-tenant isolation testing
- Performance testing (1000+ conversations)
- Security audit
- Bug fixes dan polish

### Critical Path

```
Epic 1 (Sales Team) → Epic 6 (Sales Upload)
                    ↓
                Epic 4 (Dashboard Integration)
```

**Parallel Tracks:**
- Epic 2 (Conversations) dapat dikerjakan parallel dengan Epic 6
- Epic 3 (WhatsApp) dapat dikerjakan parallel dengan Epic 2
- Epic 5 (Error Handling) incremental across all epics

### Key Highlights

**Leveraging Existing Code:**
- ✅ Repositories sudah ada (sales, conversation, car)
- ✅ WhatsApp handler sudah ada
- ✅ LLM bot sudah ada
- ✅ Database schema sudah lengkap
- 🆕 Yang perlu: Handlers baru (2 files), Templates (3 files), Routes activation

**New Capabilities (Epic 6):**
- 🌟 Sales dapat upload mobil via WhatsApp (game changer untuk dealership)
- 🌟 AI parsing untuk car details (flexible input)
- 🌟 Role-based bot behavior (customer vs sales)
- 🌟 Auto-catalog update (instant visibility)

**Quality Assurance (Epic 5):**
- Form validation di client dan server
- Consistent error responses
- Network resilience dengan retry
- Graceful degradation tanpa JavaScript

### Success Metrics

**After MVP completion, admins dapat:**
1. ✅ Manage sales team (add, view, delete) dalam < 1 menit
2. ✅ Pair WhatsApp bot dalam < 5 menit
3. ✅ Test bot dan lihat response dalam < 30 detik
4. ✅ View semua conversations dengan filter dan search
5. ✅ Sales upload mobil baru via WhatsApp dalam < 2 menit
6. ✅ Monitor activity via dashboard at-a-glance

**Technical Metrics:**
- API response time: < 200ms
- Page load: < 2s
- WhatsApp upload: < 3s (photo + details)
- Multi-tenant isolation: 100% (no data leakage)

### Next Steps After Epic Breakdown

1. **Validate dengan stakeholders** - Review epic structure dan priorities
2. **UX Design (Optional)** - Run `/bmad:bmm:workflows:create-ux-design` untuk design mockups
3. **Architecture** - Run `/bmad:bmm:workflows:create-architecture` untuk technical spec
4. **Sprint Planning** - Run `/bmad:bmm:workflows:sprint-planning` untuk organize stories into sprints
5. **Implementation** - Run `/bmad:bmm:agents:dev` untuk start development

### Notes for Developers

**Tech Stack:**
- Backend: Go + Chi Router
- Frontend: HTMX + Alpine.js + Tailwind CSS v4
- Database: PostgreSQL
- WhatsApp: Whatsmeow
- LLM: Z.AI (glm-4.6)

**Coding Standards:**
- Follow existing handler patterns (car_handler.go)
- Use tenant extraction middleware
- Parameterized queries (SQL injection prevention)
- Error messages dalam Bahasa Indonesia
- HTMX untuk dynamic updates
- Alpine.js untuk client-side interactivity
- Tailwind untuk consistent styling

**Testing Checklist:**
- [ ] Multi-tenant isolation verified
- [ ] All forms validated (client + server)
- [ ] Error states handled gracefully
- [ ] Mobile responsive
- [ ] Performance targets met
- [ ] Security audit passed

---

**✅ Epic breakdown complete dengan 100% FR coverage.**

**Ready untuk architecture workflow dan implementation planning.**

{{fr_coverage_matrix}}

---

## Summary

{{epic_breakdown_summary}}

---

_For implementation: Use the `create-story` workflow to generate individual story implementation plans from this epic breakdown._

_This document will be updated after UX Design and Architecture workflows to incorporate interaction details and technical decisions._
