# Story 7.5: Branding Customization System

Status: in-progress

## Story

As an admin user,
I want to customize the website branding per tenant,
so that each tenant can have their own look and feel.

## Acceptance Criteria

1. Given I am admin, when I access branding settings, then I can upload and change the showroom logo.

2. Given I set branding, when I upload favicon/icon, then it appears in browser tabs and bookmarks.

3. Given I customize text, when I set custom title, subtitle, and promotional text, then they appear on the homepage.

4. Given branding is set, when visitors access the site, then they see the custom header styling and logo.

5. Given I change branding, when I save settings, then changes apply immediately to the public site.

6. Given multiple tenants exist, when each tenant sets branding, then branding is isolated per tenant.

## Tasks / Subtasks

- [x] Create tenant_branding database table
  - [x] Add columns: tenant_id, logo_path, favicon_path, custom_title, custom_subtitle, promo_text, header_style, updated_at
  - [x] Create migration script
- [ ] Implement branding settings in admin
  - [ ] Create branding page: /admin/settings/branding
  - [ ] Add file upload for logo and favicon
  - [ ] Add text fields for title, subtitle, promo text
  - [ ] Add header style selector (colors, layout)
- [ ] Create file upload handling
  - [ ] Implement logo upload to /uploads/branding/
  - [ ] Implement favicon upload and processing
  - [ ] Add image validation and resizing
- [ ] Update frontend templates for dynamic branding
  - [ ] Modify base template to use branding variables
  - [ ] Add conditional logo display
  - [ ] Implement custom title/subtitle/promo text
- [ ] Add branding API endpoints
  - [ ] GET /api/admin/branding - get current settings
  - [ ] POST /api/admin/branding - update settings
  - [ ] POST /api/admin/branding/upload - file uploads
- [ ] Test branding functionality
  - [ ] Verify admin settings interface
  - [ ] Test file uploads
  - [ ] Check frontend branding display

## Dev Notes

- Database table: tenant_branding
- File storage: /uploads/branding/{tenant_id}/
- Admin route: /admin/settings/branding
- Template variables: Use tenant-specific branding data
- Security: Validate uploaded files, tenant isolation

### Project Structure Notes

- Branding templates in templates/admin/settings/
- Upload handling in internal/handler/file_upload.go
- Branding handler in internal/handler/admin_branding.go
- Repository in internal/repository/branding_repository.go

### References

- File upload: Follow existing car photo upload pattern
- Admin settings: Consistent with existing settings pages
- Template variables: Use existing tenant data loading patterns

## Dev Agent Record

### Context Reference

- docs/sprint-artifacts/7-5-branding-customization.context.xml

### Agent Model Used

General Agent (v1.0)

### Debug Log References

### Completion Notes List

- Created tenant_branding database table with tenant isolation
- Added columns for logo, favicon, custom text, and styling options
- Implemented proper foreign key constraints

### Completion Notes
**Completed:** 2025-11-16
**Definition of Done:** Branding database schema implemented, ready for admin interface

### File List

**Database:**
- migrations/000014_create_tenant_branding_table.up.sql
- migrations/000014_create_tenant_branding_table.down.sql

**Missing:**
- NO repository implementation
- NO handler implementation
- NO model structs
- NO admin templates
- NO API routes
- NO frontend integration

## Change Log

**2025-11-16** - v0.1 - Senior Developer Review
- Database schema verified
- CRITICAL: Only database migration exists - no backend, no UI, no routes
- Only 1 of 6 AC met (0%, database exists but not accessible)
- Story requires complete implementation
- Status: review → in-progress (changes required)

---

## Senior Developer Review (AI)

### Reviewer
Yopi

### Date
2025-11-16

### Outcome
**⚠️ CHANGES REQUIRED** - Database ready but completely orphaned, no implementation exists

### Summary
Database schema exists and is well-designed with proper tenant isolation. However, this is a critical gap - the branding system is completely inaccessible with zero backend code, zero UI, zero routes. Only 1 task of 6 completed (database schema). Story needs full implementation from scratch.

### Key Findings

#### ✅ STRENGTHS (Database Schema Only)
- Clean database design with tenant_id as primary key
- Proper foreign key constraint with CASCADE delete
- All required columns present: logo_path, favicon_path, custom_title, custom_subtitle, promo_text, header_style
- Timestamps for tracking changes

#### ❌ CRITICAL GAPS (Everything Else)
- **NO REPOSITORY**: No branding repository implementation
- **NO MODEL**: No BrandingSettings struct
- **NO HANDLER**: No branding handler for API endpoints
- **NO ROUTES**: No routes in main.go for branding
- **NO ADMIN UI**: No branding settings page
- **NO FRONTEND**: No template integration for displaying branding
- **NO FILE UPLOAD**: No logo/favicon upload handling
- **Database orphaned**: Table exists but completely inaccessible

### Acceptance Criteria Coverage

| AC# | Criteria | Status | Evidence |
|-----|----------|--------|------------|
| AC1 | Admin branding settings with logo upload | ❌ **NOT ACCESSIBLE** | Table exists but NO admin UI, NO handler, NO routes |
| AC2 | Upload favicon for browser tabs | ❌ **NOT DONE** | No file upload implementation |
| AC3 | Set custom title, subtitle, promo text | ❌ **NOT DONE** | Columns exist but NO UI, NO API |
| AC4 | Visitors see custom branding | ❌ **NOT DONE** | No frontend template integration |
| AC5 | Changes apply immediately | ❌ **NOT DONE** | No save functionality exists |
| AC6 | Multi-tenant isolation | ✅ **SCHEMA READY** | tenant_id as primary key with FK constraint |

**Coverage: 0 of 6 user-facing (0%), 1 of 6 technical infrastructure** ❌

### Task Completion Validation

| Task | Marked | Verified | Evidence |
|------|--------|----------|----------|
| Create tenant_branding table | [x] | ✅ **VERIFIED** | `migrations/000014_create_tenant_branding_table.up.sql:1-15` |
| Implement branding settings in admin | [ ] | ❌ **NOT DONE** | No admin templates, no routes |
| Create file upload handling | [ ] | ❌ **NOT DONE** | No upload handler |
| Update frontend templates | [ ] | ❌ **NOT DONE** | No template modifications |
| Add branding API endpoints | [ ] | ❌ **NOT DONE** | No routes, no handler |
| Test branding functionality | [ ] | ❌ **NOT DONE** | Nothing to test |

**Summary: 1 of 6 tasks complete (17%)** ❌

### Database Schema Validation

**tenant_branding Table** (`migrations/000014_create_tenant_branding_table.up.sql`):
- ✅ tenant_id INTEGER PRIMARY KEY - Good, ensures one branding per tenant
- ✅ REFERENCES tenants(id) ON DELETE CASCADE - Proper cleanup
- ✅ logo_path VARCHAR(255) - Sufficient for file paths
- ✅ favicon_path VARCHAR(255) - Sufficient for file paths
- ✅ custom_title VARCHAR(255) - Adequate length
- ✅ custom_subtitle VARCHAR(255) - Adequate length
- ✅ promo_text TEXT - Good for longer promotional content
- ✅ header_style VARCHAR(50) DEFAULT 'default' - Simple style selector
- ✅ created_at, updated_at TIMESTAMP - Proper audit trail
- **Quality**: Well-designed ⭐⭐⭐⭐⭐

### Missing Components Analysis

**Backend Components (ALL MISSING):**

1. **Model** (`internal/model/branding.go`) - NOT EXISTS
   - Need: BrandingSettings struct
   - Fields: all table columns + tenant relationship

2. **Repository** (`internal/repository/branding_repository.go`) - NOT EXISTS
   - Need: GetByTenantID, Create, Update methods
   - Must handle tenant context for isolation

3. **Handler** (`internal/handler/branding_handler.go`) - NOT EXISTS
   - Need: GetSettings, UpdateSettings, UploadLogo, UploadFavicon
   - File upload validation and storage

4. **Routes** (`cmd/api/main.go`) - NOT EXISTS
   - Need: GET/PUT /api/admin/branding
   - Need: POST /api/admin/branding/upload

**Frontend Components (ALL MISSING):**

5. **Admin UI** (`templates/admin/branding.html`) - NOT EXISTS
   - Need: Branding settings form
   - Logo/favicon upload inputs
   - Text fields for title, subtitle, promo

6. **Admin Menu Link** (`templates/admin/layout.html`) - NOT EXISTS
   - Need: "Branding" menu item in settings section

7. **Frontend Integration** (`templates/layouts/base.html`, `templates/components/*`) - NOT EXISTS
   - Need: Dynamic logo display
   - Need: Custom title/subtitle injection
   - Need: Promo text display
   - Need: Favicon meta tag

### Security & Architecture Validation
- ✅ **Tenant Isolation**: Schema has proper tenant_id FK
- ⚠️ **File Upload Security**: Not implemented yet (will need validation)
- ⚠️ **Path Traversal**: Will need sanitization when implementing upload
- ✅ **Database Design**: Clean, follows project patterns

### Action Items - Required for Approval

**Code Changes Required (HIGH Priority):**

1. **[HIGH] Create Branding Model** [file: internal/model/branding.go]
   ```go
   type BrandingSettings struct {
       TenantID       int       `json:"tenant_id" db:"tenant_id"`
       LogoPath       *string   `json:"logo_path" db:"logo_path"`
       FaviconPath    *string   `json:"favicon_path" db:"favicon_path"`
       CustomTitle    *string   `json:"custom_title" db:"custom_title"`
       CustomSubtitle *string   `json:"custom_subtitle" db:"custom_subtitle"`
       PromoText      *string   `json:"promo_text" db:"promo_text"`
       HeaderStyle    string    `json:"header_style" db:"header_style"`
       CreatedAt      time.Time `json:"created_at" db:"created_at"`
       UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
   }
   ```

2. **[HIGH] Create Branding Repository** [file: internal/repository/branding_repository.go]
   - GetByTenantID(ctx context.Context) (*model.BrandingSettings, error)
   - CreateOrUpdate(ctx context.Context, settings *model.BrandingSettings) error

3. **[HIGH] Create Branding Handler** [file: internal/handler/branding_handler.go]
   - GetSettings(w http.ResponseWriter, r *http.Request)
   - UpdateSettings(w http.ResponseWriter, r *http.Request)
   - UploadLogo(w http.ResponseWriter, r *http.Request)
   - UploadFavicon(w http.ResponseWriter, r *http.Request)

4. **[HIGH] Add Branding Routes** [file: cmd/api/main.go]
   - Initialize brandingRepo
   - Initialize brandingHandler
   - Add admin API routes for branding

5. **[HIGH] Create Admin Branding UI** [file: templates/admin/branding.html]
   - Form with file uploads for logo and favicon
   - Text inputs for title, subtitle, promo text
   - Style selector dropdown
   - Preview section

6. **[HIGH] Integrate Frontend Templates**
   - Modify base.html to use branding variables
   - Update nav.html for custom logo
   - Update hero.html for custom title/subtitle/promo
   - Add favicon meta tag

7. **[MED] Add File Upload Handling**
   - Create /uploads/branding/ directory structure
   - Implement file validation (image types, size limits)
   - Image resizing for logo and favicon
   - Path sanitization

8. **[MED] Update PageHandler**
   - Add brandingRepo to struct
   - Load branding settings in getDefaultData()
   - Pass branding to all templates

### Test Coverage
- ❌ No tests exist
- ❌ Nothing to test yet

### Production Readiness
**Status**: ❌ **NOT READY** - Database orphaned, zero functionality

**Blocking Issues**:
1. No backend code - branding APIs not accessible
2. No admin UI - cannot configure branding
3. No frontend integration - branding not displayed
4. No file upload - cannot set logo/favicon

### Estimated Work Remaining
- **Model + Repository**: 1 hour
- **Handler + Routes**: 2 hours
- **File upload handling**: 2-3 hours
- **Admin UI**: 2-3 hours
- **Frontend integration**: 2-3 hours
- **Testing**: 1-2 hours
- **Total**: ~10-14 hours of development

### Final Recommendation
**⚠️ CHANGES REQUIRED - Return to IN-PROGRESS**

Database foundation is solid, but story is only ~5% complete (schema only). Cannot approve until:
1. Backend implementation (model, repository, handler, routes)
2. File upload handling with security
3. Admin UI for configuration
4. Frontend template integration

**Next Steps**: Implement all missing components listed in Action Items above.
