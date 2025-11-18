package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/riz/auto-lmk/internal/handler"
	"github.com/riz/auto-lmk/internal/llm"
	appMiddleware "github.com/riz/auto-lmk/internal/middleware"
	"github.com/riz/auto-lmk/internal/repository"
	"github.com/riz/auto-lmk/internal/service"
	"github.com/riz/auto-lmk/internal/whatsapp"
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
	fmt.Printf("Config loaded successfully: LLM Provider = %s\n", cfg.LLM.Provider)

	// Setup logger
	logger.Setup(cfg.Server.Env)
	slog.Info("starting application", "env", cfg.Server.Env, "port", cfg.Server.Port)

	// Connect to database
	fmt.Printf("Connecting to database...\n")
	db, err := database.Connect(cfg.DatabaseURL())
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	fmt.Printf("Database connected successfully\n")

	// Initialize LLM provider if configured
	var llmProvider llm.Provider
	fmt.Printf("LLM Provider: %s, API Key: %s\n", cfg.LLM.Provider, cfg.LLM.APIKey)
	if cfg.LLM.Provider != "" && cfg.LLM.APIKey != "" {
		fmt.Printf("Initializing LLM provider...\n")
		llmCfg := llm.Config{
			Provider:    cfg.LLM.Provider,
			APIKey:      cfg.LLM.APIKey,
			Model:       cfg.LLM.Model,
			ZAIEndpoint: cfg.LLM.ZAIEndpoint,
		}

		llmProvider, err = llm.NewProvider(llmCfg)
		if err != nil {
			slog.Warn("failed to initialize LLM provider", "error", err, "provider", cfg.LLM.Provider)
			fmt.Printf("LLM provider initialization failed: %v\n", err)
		} else {
			slog.Info("LLM provider initialized", "provider", cfg.LLM.Provider, "model", cfg.LLM.Model)
			fmt.Printf("LLM provider initialized successfully\n")
		}
	} else {
		slog.Warn("LLM provider not configured, bot functionality will be disabled")
		fmt.Printf("LLM provider not configured\n")
	}

	// Setup router
	r := setupRouter(cfg, db, llmProvider)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second, // Increased for AI generation (can take 30-60s)
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

func setupRouter(cfg *config.Config, db *database.DB, llmProvider llm.Provider) *chi.Mux {
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

	// Get working directory
	workDir, _ := os.Getwd()

	// Static files with proper MIME types
	filesDir := http.Dir(workDir + "/static")
	fileServer(r, "/static", filesDir)

	// Serve uploaded files
	uploadsDir := http.Dir(workDir + "/uploads")
	r.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(uploadsDir)))

	// Initialize repositories
	tenantRepo := repository.NewTenantRepository(db.DB)
	carRepo := repository.NewCarRepository(db.DB)
	salesRepo := repository.NewSalesRepository(db.DB)
	conversationRepo := repository.NewConversationRepository(db.DB)
	analyticsRepo := repository.NewAnalyticsRepository(db.DB)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsRepo)
	blogRepo := repository.NewBlogRepository(db.DB)

	// Initialize WhatsApp settings repository
	whatsappSettingsRepo := repository.NewWhatsAppSettingsRepository(db.DB)

	// Initialize branding repository
	brandingRepo := repository.NewBrandingRepository(db.DB)

	// Initialize showroom repository
	showroomRepo := repository.NewShowroomRepository(db.DB)

	// Initialize services
	carService := service.NewCarService(carRepo)

	// Initialize WhatsApp client if LLM is configured
	var waClient *whatsapp.Client
	var waService *service.WhatsAppService
	if llmProvider != nil {
		var err error
		waClient, err = whatsapp.NewClient(salesRepo, cfg.DatabaseURL())
		if err != nil {
			slog.Error("failed to initialize WhatsApp client", "error", err)
		} else {
			slog.Info("WhatsApp client initialized")

			// Initialize bot with adapters
			convAdapter := llm.NewConversationRepoAdapter(conversationRepo)
			carAdapter := llm.NewCarRepoAdapter(carRepo)
			bot := llm.NewBot(llmProvider, convAdapter, carAdapter)

			// Initialize WhatsApp service
			waService = service.NewWhatsAppService(waClient, bot, salesRepo, conversationRepo, carService)

			// Set message handler
			waClient.SetMessageHandler(waService.ProcessIncomingMessage)

			slog.Info("WhatsApp bot service initialized")
		}
	}

	// Initialize handlers
	tenantHandler := handler.NewTenantHandler(tenantRepo)
	var carHandler *handler.CarHandler
	if llmProvider != nil {
		carHandler = handler.NewCarHandlerWithAnalyticsAndLLM(carRepo, analyticsRepo, llmProvider)
	} else {
		carHandler = handler.NewCarHandlerWithAnalytics(carRepo, analyticsRepo)
	}
	salesHandler := handler.NewSalesHandler(salesRepo)
	conversationHandler := handler.NewConversationHandler(conversationRepo)
	pageHandler := handler.NewPageHandler(carRepo, salesRepo, tenantRepo, conversationRepo, blogRepo, brandingRepo, showroomRepo)

	// Blog handler (with LLM support if available)
	var blogHandler *handler.BlogHandler
	if llmProvider != nil {
		blogHandler = handler.NewBlogHandlerWithLLM(blogRepo, llmProvider)
	} else {
		blogHandler = handler.NewBlogHandler(blogRepo)
	}

	// Branding handler
	brandingHandler := handler.NewBrandingHandler(brandingRepo)

	// Showroom handler
	showroomHandler := handler.NewShowroomHandler(showroomRepo)

	// WhatsApp handler (if WhatsApp client is initialized)
	var whatsappHandler *handler.WhatsAppHandler
	if waClient != nil {
		whatsappHandler = handler.NewWhatsAppHandlerWithSettings(waClient, tenantRepo, whatsappSettingsRepo)
	}

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
				r.Post("/ai-generate", carHandler.AIGenerate)
				r.Get("/", carHandler.List)
				r.Get("/search", carHandler.Search)
				r.Get("/{id}", carHandler.Get)
				r.Put("/{id}", carHandler.Update)
				r.Delete("/{id}", carHandler.Delete)
				r.Post("/{id}/photos", carHandler.UploadPhotos)
				r.Delete("/photos/{photoId}", carHandler.DeletePhoto)
			})

			// Sales management
			r.Route("/sales", func(r chi.Router) {
				r.Post("/", salesHandler.Create)
				r.Get("/", salesHandler.List)
				r.Get("/stats", salesHandler.Stats)
				r.Delete("/{id}", salesHandler.Delete)
			})

			// Conversations
			r.Route("/conversations", func(r chi.Router) {
				r.Get("/", conversationHandler.List)
				r.Get("/stats", conversationHandler.Stats)
				r.Get("/{id}", conversationHandler.Get)
			})

			// Analytics admin routes (tenant-scoped)
			r.Route("/admin/analytics", func(r chi.Router) {
				r.Get("/search-keywords", analyticsHandler.GetTopKeywords)
				r.Get("/car-views", analyticsHandler.GetTopCars)
				r.Get("/trends", analyticsHandler.GetTrends)
				r.Get("/export", analyticsHandler.ExportCSV)
			})

			// WhatsApp admin routes (tenant-scoped)
			if whatsappHandler != nil {
				r.Route("/admin/whatsapp", func(r chi.Router) {
					r.Get("/status", whatsappHandler.GetStatus)
					r.Post("/pair", whatsappHandler.InitiatePairing)
					r.Post("/disconnect", whatsappHandler.Disconnect)
					r.Post("/test", whatsappHandler.SendTestMessage)
					r.Get("/qr/{tenant_id}", whatsappHandler.GetQRCodeImage)
					r.Get("/settings", whatsappHandler.GetSettings)
					r.Put("/settings", whatsappHandler.UpdateSettings)
					r.Get("/effective-number", whatsappHandler.GetEffectiveNumber)
				})
			}

			// Blog admin routes (tenant-scoped)
			r.Route("/admin/blog", func(r chi.Router) {
				r.Get("/", blogHandler.List)
				r.Post("/", blogHandler.Create)
				r.Get("/{id}", blogHandler.Get)
				r.Put("/{id}", blogHandler.Update)
				r.Delete("/{id}", blogHandler.Delete)
				r.Post("/generate-ai", blogHandler.GenerateAI)
			})

			// Branding admin routes (tenant-scoped)
			r.Route("/admin/branding", func(r chi.Router) {
				r.Get("/", brandingHandler.GetSettings)
				r.Put("/", brandingHandler.UpdateSettings)
				r.Post("/upload-logo", brandingHandler.UploadLogo)
				r.Post("/upload-favicon", brandingHandler.UploadFavicon)
			})

			// Showroom admin routes (tenant-scoped)
			r.Route("/admin/showroom", func(r chi.Router) {
				r.Get("/", showroomHandler.GetAdminSettings)
				r.Put("/", showroomHandler.UpdateSettings)
			})

			// Public showroom route
			r.Get("/showroom", showroomHandler.GetSettings)
		})
	})

	// Public frontend routes (with tenant middleware for multi-tenant support)
	r.Group(func(r chi.Router) {
		r.Use(appMiddleware.TenantExtractor(db.DB))

		r.Get("/", pageHandler.Home)
		r.Get("/mobil", pageHandler.Cars)
		r.Get("/mobil/{id}", pageHandler.CarDetail)
		r.Get("/kontak", pageHandler.Contact)
		r.Get("/blog", pageHandler.BlogList)
		r.Get("/blog/{slug}", pageHandler.BlogDetail)
	})

	// Logout route
	r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Admin frontend routes (with tenant middleware)
	r.Route("/admin", func(r chi.Router) {
		r.Use(appMiddleware.TenantExtractor(db.DB))

		r.Get("/", pageHandler.AdminDashboard)
		r.Get("/dashboard", pageHandler.AdminDashboard)
		r.Get("/cars", pageHandler.AdminCars)
		r.Get("/analytics", pageHandler.AdminAnalytics)
		r.Get("/cars/new", pageHandler.AdminCarsNew)
		r.Get("/cars/{id}/edit", pageHandler.AdminCarsEdit)
		r.Get("/cars/{id}/edit", pageHandler.AdminCarsEdit)
		r.Get("/sales", pageHandler.AdminSales)
		r.Get("/sales/new", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/admin/sales", http.StatusSeeOther)
		})
		r.Get("/sales/table", pageHandler.AdminSalesTable)
		r.Get("/whatsapp", pageHandler.AdminWhatsApp)
		r.Get("/conversations", pageHandler.AdminConversations)
		r.Get("/conversations/table", pageHandler.AdminConversationsTable)
		r.Get("/conversations/{id}", pageHandler.AdminConversationDetail)
		r.Get("/blog", pageHandler.AdminBlog)
		r.Get("/blog/new", pageHandler.AdminBlogNew)
		r.Get("/blog/{id}/edit", pageHandler.AdminBlogEdit)
		r.Get("/settings", pageHandler.AdminSettings)
		r.Get("/branding", pageHandler.AdminBranding)
		r.Get("/showroom", pageHandler.AdminShowroom)
	})

	// Suppress unused variable warnings (will be used when handlers are implemented)
	_ = waService

	return r
}

// fileServer sets up a http.FileServer handler to serve static files with proper MIME types
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
