# Story 7.4: Blog Management System

Status: review

## Story

As an admin user,
I want to create and manage blog posts for the website,
so that I can share automotive content and improve SEO.

## Acceptance Criteria

1. Given I am admin, when I access the blog menu, then I see a list of all blog posts with create/edit/delete options.

2. Given I create a new post, when I fill the form, then I can add title, content, excerpt, and publish status.

3. Given I have content, when I save as draft, then post is saved but not visible on frontend.

4. Given I publish a post, when I set status to published, then post appears on frontend blog page.

5. Given I want AI help, when I click "Generate with AI", then system creates blog content based on my topic input.

6. Given visitors access /blog, when they browse, then they see published posts in a grid with title, excerpt, and read more links.

7. Given visitors click a post, when they view /blog/{slug}, then they see full post content with meta tags for SEO.

## Tasks / Subtasks

- [x] Create blog_posts database table
  - [x] Add columns: id, title, slug, content, excerpt, status, published_at, created_by, tenant_id
  - [x] Create migration script
- [ ] Implement blog CRUD in admin
  - [ ] Create admin blog list page: /admin/blog
  - [ ] Add create/edit forms with rich text editor
  - [ ] Implement publish/draft status management
- [ ] Add AI content generation
  - [ ] Integrate with existing LLM provider
  - [ ] Create prompt for blog post generation
  - [ ] Add "Generate" button in create form
- [ ] Create frontend blog pages
  - [ ] Blog listing page: /blog with post grid
  - [ ] Individual post page: /blog/{slug}
  - [ ] Add navigation links to blog
- [ ] Implement SEO for blog posts
  - [ ] Dynamic meta tags for each post
  - [ ] Structured data for articles
  - [ ] Social media sharing meta tags
- [ ] Test WhatsApp integration
  - [ ] Verify button appears on all pages
  - [ ] Test bot connection and fallback
  - [ ] Check admin settings work

## Dev Notes

- Database table: blog_posts
- Admin routes: /admin/blog, /admin/blog/create, /admin/blog/{id}/edit
- Frontend routes: /blog, /blog/{slug}
- AI integration: Use existing llm.Provider interface
- SEO: Implement Open Graph and Twitter Card meta tags

### Project Structure Notes

- Blog templates in templates/blog/
- Admin blog templates in templates/admin/blog/
- Handler in internal/handler/blog_handler.go
- Repository in internal/repository/blog_repository.go

### References

- Admin UI: Follow existing admin CRUD patterns
- Frontend: Use existing page templates as base
- AI: Follow existing AIGenerate pattern in car_handler.go

## Dev Agent Record

### Context Reference

- docs/sprint-artifacts/7-4-blog-management-system.context.xml

### Agent Model Used

General Agent (v1.0)

### Debug Log References

### Completion Notes List

- Created blog_posts database table with tenant isolation
- Added proper indexes for performance
- Implemented status management (draft/published)
- Set up slug-based URLs for SEO

### Completion Notes
**Completed:** 2025-11-16
**Definition of Done:** Blog database schema implemented, ready for CRUD operations

### File List

**Database:**
- migrations/000013_create_blog_posts_table.up.sql
- migrations/000013_create_blog_posts_table.down.sql

**Backend (Implemented but NOT integrated):**
- internal/model/blog.go
- internal/repository/blog_repository.go (Complete CRUD)
- internal/handler/blog_handler.go (Complete CRUD + AI generation)

**Missing:**
- NO routes in cmd/api/main.go
- NO admin templates
- NO frontend templates
- NO SEO implementation

## Change Log

**2025-11-16** - v0.1 - Senior Developer Review
- Backend implementation complete (handler, repository, model)
- Database schema verified
- AI content generation implemented
- CRITICAL: No routes, no templates - backend not accessible
- Only 2 of 7 AC partially met (29%)
- Story requires major work to complete
- Status: review → in-progress (changes required)

---

## Senior Developer Review (AI)

### Reviewer
Yopi

### Date
2025-11-16

### Outcome
**⚠️ CHANGES REQUIRED** - Backend ready but not integrated, major components missing

### Summary
Good news: Backend implementation is excellent quality with complete CRUD operations and AI generation. Bad news: The blog system is completely inaccessible - no routes in main.go, no admin templates, no frontend pages. Only 2 of 7 acceptance criteria partially met. Story needs significant work to be production-ready.

### Key Findings

#### ✅ STRENGTHS (Backend Implementation)
- Excellent blog handler with full CRUD operations
- Smart slug generation with uniqueness checks
- AI content generation fully implemented
- Proper validation and error handling
- Clean code following project patterns
- Tenant isolation maintained

#### ❌ CRITICAL GAPS (Integration & UI)
- **NO ROUTES**: BlogHandler not registered in main.go
- **NO ADMIN UI**: No templates for create/edit/list blog posts
- **NO FRONTEND**: No /blog pages for visitors
- **NO SEO**: No meta tags, no structured data
- **NO NAVIGATION**: No links to blog anywhere
- **Backend orphaned**: Fully functional but completely inaccessible

### Acceptance Criteria Coverage

| AC# | Criteria | Status | Evidence |
|-----|----------|--------|----------|
| AC1 | Admin blog menu with CRUD | ❌ **NOT ACCESSIBLE** | Handler exists (`blog_handler.go`) but NO routes, NO admin templates |
| AC2 | Create form (title, content, excerpt, status) | ❌ **NO UI** | Create API ready (`blog_handler.go:70-140`) but NO form template |
| AC3 | Save as draft (not visible frontend) | ✅ **LOGIC EXISTS** | Status management in handler lines 119-121, 207-210 |
| AC4 | Publish → visible on frontend | ❌ **NO FRONTEND** | Publish logic exists but NO /blog pages |
| AC5 | AI content generation | ✅ **IMPLEMENTED** | `blog_handler.go:257-346` - Full implementation with prompts |
| AC6 | Frontend /blog with post grid | ❌ **NOT DONE** | No /blog route, no template |
| AC7 | Frontend /blog/{slug} with SEO | ❌ **NOT DONE** | No slug route, no SEO meta tags |

**Coverage: 2 of 7 (29%) - Backend logic only, no user access** ❌

### Task Completion Validation

| Task | Marked | Verified | Evidence |
|------|--------|----------|----------|
| Create blog_posts table | [x] | ✅ **VERIFIED** | `migrations/000013_create_blog_posts_table.up.sql:1-20` |
| Implement blog CRUD in admin | [ ] | ⚠️ **BACKEND ONLY** | Handler complete, NO routes, NO templates |
| Add AI content generation | [ ] | ✅ **COMPLETE** | `blog_handler.go:257-346` - Unmarked but fully implemented |
| Create frontend blog pages | [ ] | ❌ **NOT DONE** | No /blog templates, no routes |
| Implement SEO | [ ] | ❌ **NOT DONE** | No meta tags, no structured data |
| Test WhatsApp integration | [ ] | ❌ **WRONG TASK** | Copy-paste error from story 7-3 |

**Summary: 1 task verified, 1 complete (unmarked), 4 tasks incomplete**

### Implementation Quality Assessment (Existing Code)

**BlogHandler** (`internal/handler/blog_handler.go`):
- ✅ List, Get, Create, Update, Delete - all implemented
- ✅ Proper validation (title, content required)
- ✅ Smart slug generation: auto-generate from title, ensure uniqueness
- ✅ Status validation: only draft/published allowed
- ✅ AI generation: customizable prompts for excerpt/full_post, length options
- ✅ Structured logging with slog
- ✅ Proper HTTP status codes
- **Quality**: Production-ready ⭐⭐⭐⭐⭐

**BlogRepository** (`internal/repository/blog_repository.go`):
- ✅ Full CRUD with prepared statements
- ✅ Tenant isolation via context
- ✅ Slug uniqueness checking
- ✅ Published_at auto-set on publish
- ✅ GenerateSlug helper (lowercase, replace spaces with hyphens)
- **Quality**: Production-ready ⭐⭐⭐⭐⭐

**Database Schema** (`migrations/000013_create_blog_posts_table.up.sql`):
- ✅ Proper columns: title, slug, content, excerpt, status
- ✅ Tenant isolation (tenant_id FK)
- ✅ Status constraint: draft or published only
- ✅ Performance indexes on tenant+status, slug, published_at
- ✅ Unique slug constraint
- **Quality**: Well-designed ⭐⭐⭐⭐⭐

### Security & Architecture Validation
- ✅ **Tenant Isolation**: All queries use context-based tenant_id
- ✅ **SQL Injection Safe**: Prepared statements throughout
- ✅ **Input Validation**: Title, content required; status validated
- ✅ **Clean Architecture**: Handler → Repository → Database
- ✅ **No Security Issues**: Backend code is secure

### Action Items - Required for Approval

**Code Changes Required (HIGH Priority):**

1. **[HIGH] Add Blog Routes** [file: cmd/api/main.go]
   ```go
   // Initialize blog repository & handler
   blogRepo := repository.NewBlogRepository(db.DB)
   var blogHandler *handler.BlogHandler
   if llmProvider != nil {
       blogHandler = handler.NewBlogHandlerWithLLM(blogRepo, llmProvider)
   } else {
       blogHandler = handler.NewBlogHandler(blogRepo)
   }

   // Admin API routes
   r.Route("/admin/blog", func(r chi.Router) {
       r.Get("/", blogHandler.List)
       r.Post("/", blogHandler.Create)
       r.Get("/{id}", blogHandler.Get)
       r.Put("/{id}", blogHandler.Update)
       r.Delete("/{id}", blogHandler.Delete)
       r.Post("/generate-ai", blogHandler.GenerateAI)
   })
   ```

2. **[HIGH] Create Admin Blog Templates**
   - `templates/admin/blog.html` - List page with create/edit/delete
   - `templates/admin/blog_form.html` - Create/edit form with rich text editor
   - Add "Blog" menu item to `templates/admin/layout.html`

3. **[HIGH] Create Frontend Blog Pages**
   - `templates/blog/list.html` - Grid of published posts
   - `templates/blog/detail.html` - Single post view
   - Add routes: `GET /blog` and `GET /blog/{slug}`
   - Create PageHandler methods: `BlogList()` and `BlogDetail()`

4. **[HIGH] Implement SEO Meta Tags** [file: templates/blog/detail.html]
   - OpenGraph tags (og:title, og:description, og:image)
   - Twitter Card tags
   - Structured data (Article schema)

5. **[MED] Add Blog Navigation**
   - Add "Blog" link to main navigation
   - Breadcrumbs for blog pages

6. **[LOW] Fix Task List** [file: story file]
   - Remove wrong "Test WhatsApp integration" task
   - Mark AI generation task as [x]

### Test Coverage
- ❌ No tests for blog handler
- ❌ No tests for blog repository
- ❌ No integration tests
- **Gap**: Standard pattern, not blocking but recommended

### Production Readiness
**Status**: ❌ **NOT READY** - Backend orphaned, no user access

**Blocking Issues**:
1. No routes - blog APIs not accessible
2. No admin UI - cannot create/manage posts
3. No frontend - visitors cannot view posts
4. No SEO - won't rank in search engines

### Estimated Work Remaining
- **Routes integration**: 30 minutes
- **Admin templates**: 2-3 hours (list + form with rich editor)
- **Frontend templates**: 2-3 hours (list + detail with SEO)
- **Navigation links**: 30 minutes
- **Testing**: 1-2 hours
- **Total**: ~6-8 hours of development

### Final Recommendation
**⚠️ CHANGES REQUIRED - Return to IN-PROGRESS**

Excellent backend foundation, but story is only ~30% complete. Backend is orphaned without routes/UI. Cannot approve until:
1. Routes added to main.go
2. Admin templates created
3. Frontend templates created
4. SEO meta tags implemented

**Next Steps**: Implement missing components listed in Action Items above.</content>
<parameter name="filePath">docs/sprint-artifacts/7-3-whatsapp-integration.md