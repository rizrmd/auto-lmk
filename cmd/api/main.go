package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/riz/auto-lmk/internal/handler"
	appMiddleware "github.com/riz/auto-lmk/internal/middleware"
	"github.com/riz/auto-lmk/internal/repository"
	"github.com/riz/auto-lmk/pkg/config"
	"github.com/riz/auto-lmk/pkg/database"
	"github.com/riz/auto-lmk/pkg/logger"
)

func main() {
	fmt.Println("Auto LMK - Multi-Tenant Car Sales Platform")
	fmt.Println("===========================================")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Setup logger
	logger.Setup(cfg.Server.Env)
	slog.Info("starting application", "env", cfg.Server.Env, "port", cfg.Server.Port)

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL())
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Setup router
	r := setupRouter(db)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		slog.Info("server starting", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server stopped")
}

func setupRouter(db *database.DB) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Health(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Database unhealthy"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Initialize repositories
	tenantRepo := repository.NewTenantRepository(db.DB)
	carRepo := repository.NewCarRepository(db.DB)
	salesRepo := repository.NewSalesRepository(db.DB)
	conversationRepo := repository.NewConversationRepository(db.DB)
	leadRepo := repository.NewLeadRepository(db.DB)

	// Initialize handlers
	tenantHandler := handler.NewTenantHandler(tenantRepo)
	carHandler := handler.NewCarHandler(carRepo)

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"message":"Auto LMK API","version":"1.0.0"}`))
		})

		// Root admin routes (no tenant middleware)
		r.Route("/admin", func(r chi.Router) {
			// Tenant management
			r.Route("/tenants", func(r chi.Router) {
				r.Post("/", tenantHandler.Create)
				r.Get("/", tenantHandler.List)
				r.Get("/{id}", tenantHandler.Get)
			})
		})

		// Tenant-scoped routes (with tenant middleware)
		r.Group(func(r chi.Router) {
			r.Use(appMiddleware.TenantExtractor(db.DB))

			// Car management
			r.Route("/cars", func(r chi.Router) {
				r.Post("/", carHandler.Create)
				r.Get("/", carHandler.List)
				r.Get("/search", carHandler.Search)
				r.Get("/{id}", carHandler.Get)
				r.Put("/{id}", carHandler.Update)
				r.Delete("/{id}", carHandler.Delete)
			})

			// Sales management (TODO: implement handler)
			r.Route("/sales", func(r chi.Router) {
				// r.Post("/", salesHandler.Create)
				// r.Get("/", salesHandler.List)
				// r.Delete("/{id}", salesHandler.Delete)
			})

			// Leads (TODO: implement handler)
			r.Route("/leads", func(r chi.Router) {
				// r.Get("/", leadHandler.List)
				// r.Get("/{id}", leadHandler.Get)
				// r.Put("/{id}/status", leadHandler.UpdateStatus)
			})

			// Conversations (TODO: implement handler)
			r.Route("/conversations", func(r chi.Router) {
				// r.Get("/", conversationHandler.List)
				// r.Get("/{id}", conversationHandler.Get)
			})
		})
	})

	// Suppress unused variable warnings
	_ = conversationRepo
	_ = leadRepo
	_ = salesRepo

	return r
}
