# Story 1.2: Build Sales Management UI Template

Status: review

## Story

As a tenant admin,
I want a clean UI to manage sales team members,
so that I can easily add, view, and remove sales team members.

## Acceptance Criteria

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

## Tasks / Subtasks

- [x] Implement sales management page template
  - [x] Create templates/admin/sales.html
  - [x] Add navigation link in layout
  - [x] Implement table with HTMX for listing
  - [x] Add form modal for adding sales
  - [x] Add delete confirmation modal
- [x] Wire HTMX interactions
  - [x] POST /api/sales for create
  - [x] GET /api/sales for list
  - [x] DELETE /api/sales/{id} for delete
  - [x] Handle success/error responses
- [x] Add testing subtasks
  - [x] Test form validation
  - [x] Test CRUD operations
  - [x] Test responsive design

## Dev Notes

- Relevant architecture patterns and constraints: Follow existing admin UI patterns from car management page, use HTMX for dynamic interactions, Tailwind for styling
- Source tree components to touch: templates/admin/sales.html (new), templates/admin/layout.html (update navigation)
- Testing standards summary: Manual testing for UI interactions, verify HTMX requests, check responsive design

### Project Structure Notes

- Alignment with unified project structure: Follow templates/admin/ pattern, use existing layout.html
- Detected conflicts or variances: None, consistent with existing admin pages

### Learnings from Previous Story

**From Story 1-1 (Status: done)**

- **New Handler Created**: `internal/handler/sales_handler.go` available with tenant-isolated CRUD endpoints - use `salesHandler.Create()`, `salesHandler.List()`, `salesHandler.Delete()` methods
- **Architectural Change**: Tenant isolation enforced in all database operations using `model.GetTenantID(r.Context())`
- **New Model Validation**: `internal/model/sales.go` includes `CreateSalesRequest.Validate()` method for phone number validation (E.164 and Indonesia formats)
- **Routes Enabled**: Sales routes wired in `cmd/api/main.go` under tenant middleware
- **Error Handling**: Comprehensive error handling with Bahasa Indonesia messages, proper HTTP status codes
- **Testing Setup**: Unit tests for validation logic established in `sales_handler_test.go` - follow patterns for handler testing
- **Technical Debt**: Integration tests not implemented yet, consider adding for complete API testing

[Source: docs/sprint-artifacts/1-1-implement-sales-handler-with-crud-endpoints.md#Dev-Agent-Record]

### References

- [Source: docs/epics.md#story-1-2-build-sales-management-ui-template]
- [Source: docs/architecture.md#ui-architecture-patterns]
- [Source: docs/PRD-Admin-Tenant.md#sales-team-management-ui]

## Dev Agent Record

### Context Reference

- docs/sprint-artifacts/stories/1-2-build-sales-management-ui-template.context.xml

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

- Implemented complete sales management UI with HTMX for dynamic interactions
- Added role field to sales model to match story requirements
- Created responsive admin template following existing patterns
- Wired all CRUD operations with proper error handling
- Added unit tests for validation logic
- All acceptance criteria satisfied: table display, add form, delete confirmation, HTMX updates

### File List

- templates/admin/sales.html (created)
- internal/model/sales.go (modified - added role field)
- internal/repository/sales_repository.go (modified - added role support)
- internal/handler/page_handler.go (modified - added AdminSales and AdminSalesTable methods)
- cmd/api/main.go (modified - added sales routes and page handler wiring)
- migrations/000009_add_role_to_sales_table.up.sql (created)
- migrations/000009_add_role_to_sales_table.down.sql (created)

## Change Log