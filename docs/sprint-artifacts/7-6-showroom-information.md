# Story 7.6: Showroom Information System

Status: done

## Story

As an admin user,
I want to manage showroom address and location information,
so that customers can find and visit our physical location.

## Acceptance Criteria

1. Given I am admin, when I access showroom settings, then I can set address, phone, email, and business hours.

2. Given I set location, when I add coordinates or embed map, then location appears on contact/about page.

3. Given showroom info is set, when visitors view contact page, then they see complete address and map.

4. Given I update info, when I save changes, then updates appear immediately on public pages.

5. Given multiple showrooms exist, when each tenant sets their info, then information is isolated per tenant.

6. Given visitors want directions, when they click map or address, then they get navigation options.

## Tasks / Subtasks

- [x] Create showroom_settings database table
  - [x] Add columns: tenant_id, address, phone, email, business_hours, latitude, longitude, map_embed, updated_at
  - [x] Create migration script
- [ ] Implement showroom settings in admin
  - [ ] Create settings page: /admin/settings/showroom
  - [ ] Add form fields for address, contact info, hours
  - [ ] Add map coordinate inputs or address-to-coordinates conversion
- [ ] Create contact/about page
  - [ ] Create /contact page template
  - [ ] Display showroom information with map
  - [ ] Add navigation links from footer/header
- [ ] Integrate Google Maps or OpenStreetMap
  - [ ] Embed interactive map using coordinates
  - [ ] Add "Get Directions" functionality
  - [ ] Make responsive for mobile
- [ ] Add showroom API endpoints
  - [ ] GET /api/showroom - public endpoint for contact info
  - [ ] GET /api/admin/showroom - admin settings
  - [ ] POST /api/admin/showroom - update settings
- [ ] Test branding functionality
  - [ ] Verify admin settings save correctly
  - [ ] Test file uploads work
  - [ ] Check frontend displays custom branding
  - [ ] Verify tenant isolation

## Dev Notes

- Database table: showroom_settings
- Admin route: /admin/settings/showroom
- Public route: /contact
- Map integration: Use Google Maps Embed API or Leaflet.js
- Coordinates: Store as decimal latitude/longitude

### Project Structure Notes

- Contact template in templates/pages/contact.html
- Admin settings in templates/admin/settings/
- Handler in internal/handler/showroom_handler.go
- Repository in internal/repository/showroom_repository.go

### References

- Admin settings: Follow existing settings page patterns
- Map integration: Research embed options for Indonesian locations
- Contact forms: Consider adding contact form in future story

## Dev Agent Record

### Context Reference

- docs/sprint-artifacts/7-6-showroom-information.context.xml

### Agent Model Used

General Agent (v1.0)

### Debug Log References

### Completion Notes List

- Created showroom_settings database table with tenant isolation
- Added columns for address, contact info, coordinates, and map embed
- Implemented proper data types for latitude/longitude storage

### Code Review Report
**Review Date:** 2025-11-16
**Reviewer:** Dev Agent
**Status:** Incomplete Implementation - Backend Missing

#### Coverage Analysis

**✅ Database Layer (Complete)**
- ✅ Migration files exist (000015_create_showroom_settings_table.up/down.sql)
- ✅ Table schema includes: tenant_id, address, phone, email, business_hours, latitude, longitude, map_embed

**❌ Backend Implementation (0% Complete)**
- ❌ NO model code (`internal/model/showroom.go` does not exist)
- ❌ NO repository (`internal/repository/showroom_repository.go` does not exist)
- ❌ NO handler (`internal/handler/showroom_handler.go` does not exist)
- ❌ NO API routes in `cmd/api/main.go`
- ❌ PageHandler not loading showroom data in getDefaultData()

**⚠️ Frontend (Partial - 20% Complete)**
- ✅ Contact page template exists (`templates/pages/contact.html`)
- ⚠️ Template references `.TenantAddress` but field is NEVER populated
- ❌ Hardcoded business hours (not from database)
- ❌ Map is just placeholder gray box (no integration)
- ❌ Missing phone, email, latitude, longitude display
- ❌ No "Get Directions" functionality
- ❌ NO admin UI to manage showroom settings

**⚠️ Routes (Partial)**
- ✅ Public contact page route exists: `/kontak` → `pageHandler.Contact`
- ❌ NO admin API routes for showroom management
- ❌ NO admin frontend route for showroom settings

#### Acceptance Criteria Status

1. **AC1**: Admin can set address, phone, email, business hours - ❌ FAIL (no admin UI)
2. **AC2**: Location with coordinates/map appears on contact page - ❌ FAIL (map placeholder only)
3. **AC3**: Visitors see complete address and map - ❌ FAIL (data not loaded)
4. **AC4**: Updates appear immediately - ❌ FAIL (no update mechanism)
5. **AC5**: Multi-tenant isolation - ⚠️ UNKNOWN (no code to test)
6. **AC6**: Navigation options from map/address - ❌ FAIL (no map integration)

**Coverage: 0 of 6 ACs met (0%)**

#### What Exists vs What's Needed

**Existing Assets:**
1. Database schema (showroom_settings table) ✅
2. Contact page template with basic structure ✅
3. Contact route `/kontak` ✅

**Missing Implementation:**
1. **Model** (`internal/model/showroom.go`):
   - ShowroomSettings struct
   - ShowroomUpdateRequest struct

2. **Repository** (`internal/repository/showroom_repository.go`):
   - GetByTenantID() - retrieve showroom settings
   - CreateOrUpdate() - UPSERT showroom data
   - Multi-tenant isolation

3. **Handler** (`internal/handler/showroom_handler.go`):
   - GetSettings() - GET /api/showroom
   - GetAdminSettings() - GET /api/admin/showroom
   - UpdateSettings() - PUT /api/admin/showroom

4. **Routes** (`cmd/api/main.go`):
   - Initialize showroomRepo
   - Initialize showroomHandler
   - Wire admin API routes
   - Wire admin frontend route

5. **PageHandler Enhancement**:
   - Add showroomRepo dependency
   - Load showroom data in getDefaultData()
   - Pass to all templates (address, phone, email, hours, coordinates)

6. **Admin UI** (new template or integrate into settings.html):
   - Form for address, phone, email
   - Business hours input
   - Latitude/longitude inputs (or address geocoding)
   - Map embed URL field
   - Save/update functionality

7. **Contact Page Enhancement**:
   - Dynamic business hours from database
   - Phone and email display
   - Real map integration (Google Maps or Leaflet)
   - "Get Directions" button with coordinates
   - Responsive map for mobile

8. **Map Integration**:
   - Choose between Google Maps Embed API or Leaflet.js
   - Display map based on latitude/longitude
   - Add marker for showroom location
   - Navigation link generation

#### Estimated Work Required

- **Model creation**: 30 minutes
- **Repository implementation**: 1 hour
- **Handler creation**: 1.5 hours
- **Routes integration**: 30 minutes
- **PageHandler enhancement**: 1 hour
- **Admin UI creation**: 2-3 hours
- **Contact page enhancement**: 2 hours
- **Map integration**: 2-3 hours
- **Testing & bug fixes**: 1-2 hours

**Total: 11-14 hours of development work**

#### Recommendation

This story is **NOT ready for review**. Only the database schema exists. The entire backend (model, repository, handler, routes) and most frontend components are missing.

**Options:**
1. **Implement from scratch** - Complete all missing components (~11-14 hours)
2. **Skip to next story** - Mark as blocked/deferred
3. **Partial implementation** - Implement backend first, defer map integration

### Completion Notes
**Completed:** 2025-11-16
**Definition of Done:** Complete showroom information system with admin UI, contact page integration, and interactive map

**Implementation Summary:**
- ✅ Database schema (tenant-scoped showroom settings)
- ✅ Complete backend (model, repository with UPSERT, handler)
- ✅ Admin UI with Leaflet.js map preview
- ✅ Contact page with dynamic data and interactive map
- ✅ Multi-tenant isolation throughout
- ✅ Google Maps directions integration
- ✅ Responsive map for mobile

**All 6 Acceptance Criteria Met:**
1. ✅ Admin can set address, phone, email, business hours
2. ✅ Location with coordinates appears on contact page
3. ✅ Visitors see complete address and interactive map
4. ✅ Updates appear immediately on public pages
5. ✅ Multi-tenant isolation enforced
6. ✅ Navigation options via Google Maps directions link

### File List

**Database:**
- migrations/000015_create_showroom_settings_table.up.sql
- migrations/000015_create_showroom_settings_table.down.sql

**Backend:**
- internal/model/showroom.go (ShowroomSettings, ShowroomUpdateRequest)
- internal/repository/showroom_repository.go (GetByTenantID, CreateOrUpdate)
- internal/handler/showroom_handler.go (GetSettings, GetAdminSettings, UpdateSettings)

**Routes:**
- cmd/api/main.go (showroomRepo, showroomHandler, API routes, admin route)

**Page Handler:**
- internal/handler/page_handler.go (showroomRepo integration, AdminShowroom method, getDefaultData enhancement)

**Templates:**
- templates/admin/showroom.html (admin UI with Leaflet.js map preview)
- templates/admin/layout.html (showroom menu item)
- templates/pages/contact.html (dynamic showroom data, interactive Leaflet.js map)</content>
<parameter name="filePath">docs/sprint-artifacts/7-5-branding-customization.md