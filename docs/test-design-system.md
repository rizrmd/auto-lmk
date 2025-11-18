# System-Level Test Design

**Date:** Sat Nov 15 2025
**Author:** Yopi
**Status:** Draft

---

## Executive Summary

**Scope:** System-level testability review for Auto LMK brownfield enhancement

**Risk Summary:**

- Total risks identified: 8
- High-priority risks (≥6): 4
- Critical categories: Security, Performance, Reliability

**Coverage Summary:**

- Unit: 60% - Business logic and data validation
- Integration: 30% - API endpoints and database operations
- E2E: 10% - Critical user journeys and WhatsApp flows

**Total effort**: 40 hours (~5 days)

---

## Testability Assessment

### Controllability: CONCERNS
- **Strengths:** Clean architecture with repository pattern enables API seeding and database reset. Multi-tenant middleware provides tenant isolation for parallel testing.
- **Concerns:** External dependencies (WhatsApp, LLM) not easily mockable. No chaos engineering setup for error condition testing.
- **Recommendations:** Implement dependency injection for WhatsApp and LLM clients to enable mocking.

### Observability: CONCERNS
- **Strengths:** Structured logging with Zerolog. Server-rendered HTML provides deterministic UI state.
- **Concerns:** No application metrics or APM integration. Limited visibility into LLM function calling performance.
- **Recommendations:** Add Prometheus metrics for API response times and error rates.

### Reliability: PASS
- **Strengths:** Stateless HTTP handlers, ACID database transactions, tenant-scoped operations prevent cross-contamination.
- **Assessment:** Components are loosely coupled with clear boundaries. Auto-cleanup patterns in place for file uploads.

---

## Architecturally Significant Requirements (ASRs)

| Risk ID | Category | Description | Probability | Impact | Score | Mitigation |
| ------- | -------- | ----------- | ----------- | ------ | ----- | ---------- |
| ASR-001 | SEC | Multi-tenant data isolation | 3 | 3 | 9 | Tenant middleware validation, parameterized queries |
| ASR-002 | PERF | API response time <200ms | 2 | 3 | 6 | Database indexing, query optimization |
| ASR-003 | REL | WhatsApp connection reliability | 3 | 2 | 6 | Connection monitoring, auto-reconnect |
| ASR-004 | SEC | LLM function calling security | 2 | 3 | 6 | Input validation, rate limiting |
| ASR-005 | PERF | Page load <2s | 2 | 2 | 4 | Asset optimization, caching strategy |
| ASR-006 | REL | File upload security | 1 | 3 | 3 | Path validation, size limits |
| ASR-007 | MAINT | Code maintainability | 2 | 1 | 2 | Test coverage targets, CI quality gates |
| ASR-008 | REL | WebSocket real-time updates | 1 | 2 | 2 | Connection handling, error recovery |

---

## Test Levels Strategy

- **Unit: 60%** - Business logic in services/repositories, input validation, data transformation. Fast feedback, isolated testing.
- **Integration: 30%** - API endpoints with database, WhatsApp service integration, LLM function calling. Validates component interaction.
- **E2E: 10%** - Critical WhatsApp flows (pairing, car upload), admin user journeys. Highest confidence for production-like scenarios.

---

## NFR Testing Approach

- **Security:** Playwright E2E for multi-tenant isolation testing, API tests for RBAC validation, OWASP ZAP for vulnerability scanning
- **Performance:** k6 load testing for API endpoints (50-100 concurrent users), Lighthouse for page performance metrics
- **Reliability:** Playwright for error handling UI, API tests for retry logic, chaos testing for external service failures
- **Maintainability:** GitHub Actions CI for test coverage (80%+), code duplication scanning, npm audit for dependencies

---

## Test Environment Requirements

- **Local Development:** Docker Compose with PostgreSQL, in-memory test data
- **CI/CD:** Ephemeral environments with test databases, mock external services
- **Staging:** Production-like environment with real WhatsApp (limited rate), test LLM API
- **Production Monitoring:** APM tools for performance metrics, error tracking

---

## Testability Concerns

### Critical Concerns (Blockers)
1. **External Service Dependencies:** WhatsApp and LLM APIs cannot be easily mocked, requiring integration testing in CI
2. **Multi-Tenant Isolation:** Must verify no data leakage between tenants in all tests
3. **AI Parsing Accuracy:** LLM function calling for car upload requires consistent parsing validation

### Minor Concerns
1. **File Upload Testing:** Local filesystem storage requires cleanup between tests
2. **WebSocket Testing:** Real-time pairing status requires special test setup
3. **Mobile Responsiveness:** Admin UI requires mobile testing for dealership field usage

---

## Recommendations for Sprint 0

### Framework Setup
1. Initialize Playwright for E2E testing with mobile viewport support
2. Set up k6 for performance testing with custom metrics
3. Configure test database isolation (separate schema per test run)
4. Implement mock adapters for WhatsApp and LLM services

### CI Pipeline Configuration
1. Multi-stage pipeline: unit → integration → e2e → performance
2. Parallel test execution with tenant isolation
3. Artifact collection (screenshots, HAR files, performance reports)
4. Quality gates: coverage 80%+, no critical vulnerabilities, performance baselines

### Test Data Strategy
1. Factory pattern for test data generation (users, cars, conversations)
2. Auto-cleanup fixtures for database and file system
3. Seed data for consistent test scenarios
4. Mock external API responses for deterministic testing

---

## Risk Assessment

### High-Priority Risks (Score ≥6)

| Risk ID | Category | Description | Probability | Impact | Score | Mitigation | Owner | Timeline |
| ------- | -------- | ----------- | ----------- | ------ | ----- | ---------- | ----- | -------- |
| R-001 | SEC | Multi-tenant data leakage | 3 | 3 | 9 | Tenant middleware validation | Dev Team | Sprint 1 |
| R-002 | PERF | API performance degradation | 2 | 3 | 6 | Database optimization | Dev Team | Sprint 1 |
| R-003 | REL | WhatsApp connection failures | 3 | 2 | 6 | Retry logic, monitoring | Dev Team | Sprint 1 |
| R-004 | SEC | LLM prompt injection | 2 | 3 | 6 | Input sanitization | Dev Team | Sprint 1 |

### Medium-Priority Risks (Score 3-4)

| Risk ID | Category | Description | Probability | Impact | Score | Mitigation | Owner |
| ------- | -------- | ----------- | ----------- | ------ | ----- | ---------- | ----- |
| R-005 | PERF | Page load performance | 2 | 2 | 4 | Asset optimization | Dev Team |
| R-006 | REL | File upload failures | 1 | 3 | 3 | Error handling | Dev Team |

### Low-Priority Risks (Score 1-2)

| Risk ID | Category | Description | Probability | Impact | Score | Action |
| ------- | -------- | ----------- | ----------- | ------ | ----- | ------ |
| R-007 | MAINT | Code quality | 2 | 1 | 2 | CI quality gates |
| R-008 | REL | WebSocket disconnects | 1 | 2 | 2 | Auto-reconnect |

### Risk Category Legend

- **TECH**: Technical/Architecture (flaws, integration, scalability)
- **SEC**: Security (access controls, auth, data exposure)
- **PERF**: Performance (SLA violations, degradation, resource limits)
- **DATA**: Data Integrity (loss, corruption, inconsistency)
- **BUS**: Business Impact (UX harm, logic errors, revenue)
- **OPS**: Operations (deployment, config, monitoring)

---

## Test Coverage Plan

### P0 (Critical) - Run on every commit

**Criteria**: Blocks core journey + High risk (≥6) + No workaround

| Requirement | Test Level | Risk Link | Test Count | Owner | Notes |
| ----------- | ---------- | --------- | ---------- | ----- | ----- |
| Multi-tenant isolation | Integration | R-001 | 5 | QA | API and database level |
| WhatsApp pairing flow | E2E | R-003 | 3 | QA | QR scan, status updates |
| Car upload via WhatsApp | E2E | R-004 | 4 | QA | Photo + text parsing |
| API performance | Integration | R-002 | 3 | QA | Response time validation |

**Total P0**: 15 tests, 30 hours

### P1 (High) - Run on PR to main

**Criteria**: Important features + Medium risk (3-4) + Common workflows

| Requirement | Test Level | Risk Link | Test Count | Owner | Notes |
| ----------- | ---------- | --------- | ---------- | ----- | ----- |
| Sales team management | Integration | - | 6 | QA | CRUD operations |
| Conversation monitoring | Integration | - | 4 | QA | List and detail APIs |
| Admin UI responsiveness | E2E | R-005 | 3 | QA | Mobile and desktop |

**Total P1**: 13 tests, 13 hours

### P2 (Medium) - Run nightly/weekly

**Criteria**: Secondary features + Low risk (1-2) + Edge cases

| Requirement | Test Level | Risk Link | Test Count | Owner | Notes |
| ----------- | ---------- | --------- | ---------- | ----- | ----- |
| Error handling | Unit | R-006 | 8 | Dev | Validation and recovery |
| File upload security | Integration | R-006 | 4 | QA | Size and type validation |
| WebSocket pairing | Integration | R-008 | 3 | QA | Connection management |

**Total P2**: 15 tests, 7.5 hours

### P3 (Low) - Run on-demand

**Criteria**: Nice-to-have + Exploratory + Performance benchmarks

| Requirement | Test Level | Test Count | Owner | Notes |
| ----------- | ---------- | ---------- | ----- | ----- |
| Performance benchmarks | Performance | 5 | QA | k6 load testing |
| LLM parsing accuracy | Integration | 6 | QA | Various input formats |
| Code maintainability | Unit | 10 | Dev | Coverage and quality |

**Total P3**: 21 tests, 5.25 hours

---

## Execution Order

### Smoke Tests (<5 min)

**Purpose**: Fast feedback, catch build-breaking issues

- [ ] API health check
- [ ] Database connection
- [ ] WhatsApp service initialization
- [ ] Basic admin login

**Total**: 4 scenarios

### P0 Tests (<10 min)

**Purpose**: Critical path validation

- [ ] Tenant isolation verification
- [ ] WhatsApp pairing success
- [ ] Car upload completion
- [ ] API response times

**Total**: 15 scenarios

### P1 Tests (<30 min)

**Purpose**: Important feature coverage

- [ ] Sales CRUD operations
- [ ] Conversation APIs
- [ ] UI responsiveness

**Total**: 13 scenarios

### P2/P3 Tests (<60 min)

**Purpose**: Full regression coverage

- [ ] Error scenarios
- [ ] File upload edge cases
- [ ] Performance validation

**Total**: 36 scenarios

---

## Resource Estimates

### Test Development Effort

| Priority | Count | Hours/Test | Total Hours | Notes |
| -------- | ----- | ---------- | ----------- | ----- |
| P0 | 15 | 2.0 | 30 | Complex setup, external services |
| P1 | 13 | 1.0 | 13 | Standard coverage |
| P2 | 15 | 0.5 | 7.5 | Simple scenarios |
| P3 | 21 | 0.25 | 5.25 | Exploratory |
| **Total** | **64** | **-** | **55.75** | **~7 days** |

### Prerequisites

**Test Data:**
- User factory (faker-based, auto-cleanup)
- Car factory (with photos, auto-cleanup)
- Conversation factory (multi-tenant aware)

**Tooling:**
- Playwright for E2E testing
- k6 for performance testing
- TestContainers for database isolation
- Mock servers for external APIs

**Environment:**
- CI runners with Docker support
- Staging environment with real WhatsApp (rate limited)
- Performance testing infrastructure

---

## Quality Gate Criteria

### Pass/Fail Thresholds

- **P0 pass rate**: 100% (no exceptions)
- **P1 pass rate**: ≥95% (waivers required for failures)
- **P2/P3 pass rate**: ≥90% (informational)
- **High-risk mitigations**: 100% complete or approved waivers

### Coverage Targets

- **Critical paths**: ≥80%
- **Security scenarios**: 100%
- **API endpoints**: ≥90%
- **Error handling**: ≥85%

### Non-Negotiable Requirements

- [ ] All P0 tests pass
- [ ] No high-risk (≥6) items unmitigated
- [ ] Multi-tenant isolation verified
- [ ] WhatsApp integration functional
- [ ] Performance targets met

---

## Mitigation Plans

### R-001: Multi-tenant data leakage (Score: 9)

**Mitigation Strategy:** Implement comprehensive tenant isolation testing and middleware validation
**Owner:** Dev Team
**Timeline:** Sprint 1
**Status:** Planned
**Verification:** Automated tests verify no cross-tenant data access

### R-002: API performance degradation (Score: 6)

**Mitigation Strategy:** Database query optimization and response time monitoring
**Owner:** Dev Team
**Timeline:** Sprint 1
**Status:** Planned
**Verification:** k6 performance tests with <200ms target

### R-003: WhatsApp connection failures (Score: 6)

**Mitigation Strategy:** Implement retry logic and connection monitoring
**Owner:** Dev Team
**Timeline:** Sprint 1
**Status:** Planned
**Verification:** E2E tests verify pairing reliability

### R-004: LLM prompt injection (Score: 6)

**Mitigation Strategy:** Input sanitization and rate limiting for LLM calls
**Owner:** Dev Team
**Timeline:** Sprint 1
**Status:** Planned
**Verification:** Security tests verify input validation

---

## Assumptions and Dependencies

### Assumptions

1. External services (WhatsApp, LLM) will be available during testing
2. Test data can be safely created and cleaned up
3. Multi-tenant isolation is properly implemented in middleware

### Dependencies

1. Framework workflow completion - Test infrastructure setup
2. CI workflow completion - Pipeline configuration
3. External API access for integration testing

### Risks to Plan

- **Risk**: External service instability
  - **Impact**: Test flakiness and false failures
  - **Contingency**: Implement comprehensive mocking strategy

---

## Approval

**Test Design Approved By:**

- [ ] Product Manager: ___________________ Date: __________
- [ ] Tech Lead: ___________________ Date: __________
- [ ] QA Lead: ___________________ Date: __________

**Comments:**

---

---

## Appendix

### Knowledge Base References

- `nfr-criteria.md` - NFR validation approach
- `test-levels-framework.md` - Test levels strategy guidance
- `risk-governance.md` - Risk identification framework
- `test-quality.md` - Quality standards

### Related Documents

- PRD: docs/PRD-Admin-Tenant.md
- Architecture: docs/architecture.md
- Epics: docs/epics.md

---

**Generated by**: BMad TEA Agent - Test Architect Module
**Workflow**: `.bmad/bmm/testarch/test-design`
**Version**: 4.0 (BMad v6)