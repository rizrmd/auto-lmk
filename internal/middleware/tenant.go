package middleware

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strings"

	"github.com/riz/auto-lmk/internal/model"
)

// TenantExtractor extracts tenant ID from domain and adds to context
func TenantExtractor(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract domain from request host
			host := r.Host

			// Remove port if present
			if colonIdx := strings.Index(host, ":"); colonIdx != -1 {
				host = host[:colonIdx]
			}

			// Special handling for root admin domain and localhost (development)
			if host == "admin.localhost" || host == "admin.platform.com" || host == "localhost" {
				// Root admin or development - try to get default tenant
				var tenantID int
				err := db.QueryRow("SELECT id FROM tenants WHERE status = 'active' ORDER BY id LIMIT 1").Scan(&tenantID)
				if err == nil {
					// Found a tenant, add to context
					ctx := model.WithTenantID(r.Context(), tenantID)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
				// No tenant found, continue without tenant context
				next.ServeHTTP(w, r)
				return
			}

			// Look up tenant by domain
			var tenantID int
			err := db.QueryRow("SELECT id FROM tenants WHERE domain = $1 AND status = 'active'", host).Scan(&tenantID)
			if err != nil {
				if err == sql.ErrNoRows {
					slog.Warn("tenant not found", "domain", host)
					http.Error(w, "Tenant not found", http.StatusNotFound)
					return
				}
				slog.Error("failed to lookup tenant", "error", err, "domain", host)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			// Add tenant ID to context
			ctx := model.WithTenantID(r.Context(), tenantID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
