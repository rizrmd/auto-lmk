package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// ErrorResponse represents a standardized API error response
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Code    string                 `json:"code,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// WriteError writes a standardized JSON error response
func WriteError(w http.ResponseWriter, status int, message, code string, details map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   message,
		Code:    code,
		Details: details,
	})
}

// BadRequest returns a 400 Bad Request error
func BadRequest(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusBadRequest, message, "BAD_REQUEST", nil)
}

// BadRequestWithDetails returns a 400 Bad Request error with field details
func BadRequestWithDetails(w http.ResponseWriter, message string, details map[string]interface{}) {
	WriteError(w, http.StatusBadRequest, message, "BAD_REQUEST", details)
}

// ValidationError returns a 400 error with field-level validation details
func ValidationError(w http.ResponseWriter, fields map[string]string) {
	details := map[string]interface{}{
		"fields": fields,
	}
	WriteError(w, http.StatusBadRequest, "Validation failed", "VALIDATION_ERROR", details)
}

// NotFound returns a 404 Not Found error
func NotFound(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusNotFound, message, "NOT_FOUND", nil)
}

// Forbidden returns a 403 Forbidden error
func Forbidden(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusForbidden, message, "FORBIDDEN", nil)
}

// Conflict returns a 409 Conflict error
func Conflict(w http.ResponseWriter, message string, details map[string]interface{}) {
	WriteError(w, http.StatusConflict, message, "CONFLICT", details)
}

// InternalServerError returns a 500 Internal Server Error
func InternalServerError(w http.ResponseWriter, message string) {
	// Never expose sensitive information in production errors
	WriteError(w, http.StatusInternalServerError, message, "INTERNAL_ERROR", nil)
}

// Unauthorized returns a 401 Unauthorized error
func Unauthorized(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusUnauthorized, message, "UNAUTHORIZED", nil)
}

// RecoveryMiddleware recovers from panics and logs them
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("Panic recovered",
					"error", err,
					"stack", debug.Stack(),
					"path", r.URL.Path,
					"method", r.Method)

				InternalServerError(w, "Terjadi kesalahan internal server")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request
		slog.Info("HTTP Request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(wrapped, r)

		// Log response
		slog.Info("HTTP Response",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status codes
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}
