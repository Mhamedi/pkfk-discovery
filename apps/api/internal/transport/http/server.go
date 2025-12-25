package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	mw "github.com/pkfk-discovery/api/internal/middleware"
	"github.com/pkfk-discovery/api/internal/domain"
	"github.com/pkfk-discovery/api/internal/integrations/minio"
	"github.com/pkfk-discovery/api/internal/repositories/postgres"
	"github.com/pkfk-discovery/api/internal/services"
)

type Config struct {
	Addr          string
	DatabaseURL   string
	RedisURL      string
	MinIOEndpoint string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOUseSSL   bool
	JWTSecret     string
	EncryptionKey string
}

type Server struct {
	config            *Config
	router            *chi.Mux
	server            *http.Server
	logger            *logrus.Logger
	db                *pgxpool.Pool
	authService       *services.AuthService
	authMiddleware    *mw.AuthMiddleware
	connectionService *services.ConnectionService
	adapterService    *services.AdapterService
	userService       *services.UserService
	auditService      *services.AuditService
	settingsService   *services.SettingsService
	aiProviderService *services.AIProviderService
	scanService       *services.ScanService
	draftService      *services.AdapterDraftService
	bundleService     *services.AdapterBundleService
	validationService *services.ValidationService
}

func NewServer(cfg *Config) (*Server, error) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Initialize database
	db, err := initDB(cfg.DatabaseURL, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	connRepo := postgres.NewConnectionRepository(db)
	adapterRepo := postgres.NewAdapterRepository(db)
	auditRepo := postgres.NewAuditLogRepository(db)
	settingsRepo := postgres.NewSettingsRepository(db)
	aiProviderRepo := postgres.NewAIProviderRepository(db)
	scanRepo := postgres.NewScanRepository(db)
	draftRepo := postgres.NewAdapterDraftRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	connectionService, err := services.NewConnectionService(connRepo, cfg.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize connection service: %w", err)
	}
	adapterService := services.NewAdapterService(adapterRepo)
	userService := services.NewUserService(userRepo)
	auditService := services.NewAuditService(auditRepo)
	settingsService := services.NewSettingsService(settingsRepo)
	aiProviderService, err := services.NewAIProviderService(aiProviderRepo, cfg.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AI provider service: %w", err)
	}
	scanService := services.NewScanService(scanRepo)
	draftService := services.NewAdapterDraftService(draftRepo)

	// Initialize MinIO client for bundle service
	minioClient, err := minio.NewClient(
		cfg.MinIOEndpoint,
		cfg.MinIOAccessKey,
		cfg.MinIOSecretKey,
		cfg.MinIOUseSSL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	bundleService := services.NewAdapterBundleService(minioClient, cfg.EncryptionKey)
	validationService := services.NewValidationService()

	// Initialize middleware
	authMiddleware := mw.NewAuthMiddleware(cfg.JWTSecret)

	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	s := &Server{
		config:            cfg,
		router:            router,
		logger:            logger,
		db:                db,
		authService:       authService,
		authMiddleware:    authMiddleware,
		connectionService: connectionService,
		adapterService:    adapterService,
		userService:       userService,
		auditService:      auditService,
		settingsService:   settingsService,
		aiProviderService: aiProviderService,
		scanService:       scanService,
		draftService:      draftService,
		bundleService:     bundleService,
		validationService: validationService,
		server: &http.Server{
			Addr:         cfg.Addr,
			Handler:      router,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}

	s.setupRoutes()

	return s, nil
}

func (s *Server) setupRoutes() {
	// Health endpoints
	s.router.Get("/healthz", s.handleHealthz)
	s.router.Get("/readyz", s.handleReadyz)
	s.router.Get("/metrics", s.handleMetrics)
	
	// OpenAPI spec
	s.router.Get("/api/v1/openapi.json", s.handleOpenAPISpec)
	s.router.Get("/api/v1/openapi.yaml", s.handleOpenAPIYAML)

	// API routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// Auth routes
		r.Post("/auth/login", s.handleLogin)
		r.Post("/auth/refresh", s.handleRefresh)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(s.authMiddleware.Handler)

			// Adapters
			r.Route("/adapters", func(r chi.Router) {
				r.Get("/", s.handleListAdapters)
				r.Post("/", s.handleCreateAdapter)
				r.Get("/{id}", s.handleGetAdapter)
				r.Put("/{id}", s.handleUpdateAdapter)
				r.Delete("/{id}", s.handleDeleteAdapter)
				r.Post("/{id}/publish", s.handlePublishAdapter)
			})

			// Adapter drafts (Studio)
			r.Route("/studio", func(r chi.Router) {
				r.Post("/drafts", s.handleCreateDraft)
				r.Get("/drafts/{id}", s.handleGetDraft)
				r.Put("/drafts/{id}", s.handleUpdateDraft)
				r.Post("/drafts/{id}/probe", s.handleProbeDraft)
				r.Post("/drafts/{id}/validate", s.handleValidateDraft)
				r.Post("/drafts/{id}/optimize", s.handleOptimizeDraft)
				r.Post("/drafts/{id}/publish", s.handlePublishDraft)
			})

			// Connections
			r.Route("/connections", func(r chi.Router) {
				r.Get("/", s.handleListConnections)
				r.Post("/", s.handleCreateConnection)
				r.Get("/{id}", s.handleGetConnection)
				r.Put("/{id}", s.handleUpdateConnection)
				r.Delete("/{id}", s.handleDeleteConnection)
				r.Post("/{id}/test", s.handleTestConnection)
			})

			// Scans (Engine)
			r.Route("/scans", func(r chi.Router) {
				r.Get("/", s.handleListScans)
				r.Post("/", s.handleCreateScan)
				r.Get("/{id}", s.handleGetScan)
				r.Get("/{id}/results", s.handleGetScanResults)
				r.Delete("/{id}", s.handleCancelScan)
			})

			// AI Providers
			r.Route("/ai/providers", func(r chi.Router) {
				r.Get("/", s.handleListAIProviders)
				r.Post("/", s.handleCreateAIProvider)
				r.Get("/{id}", s.handleGetAIProvider)
				r.Put("/{id}", s.handleUpdateAIProvider)
				r.Delete("/{id}", s.handleDeleteAIProvider)
			})

			// Admin
			r.Route("/admin", func(r chi.Router) {
				r.Use(mw.NewRBACMiddleware(domain.RoleAdmin).Handler)
				r.Route("/users", func(r chi.Router) {
					r.Get("/", s.handleListUsers)
					r.Post("/", s.handleCreateUser)
					r.Get("/{id}", s.handleGetUser)
					r.Put("/{id}", s.handleUpdateUser)
					r.Delete("/{id}", s.handleDeleteUser)
				})

				r.Route("/audit", func(r chi.Router) {
					r.Get("/", s.handleListAuditLogs)
					r.Get("/{id}", s.handleGetAuditLog)
					r.Get("/export", s.handleExportAuditLogs)
				})
			})

			// Settings
			r.Route("/settings", func(r chi.Router) {
				r.Get("/", s.handleGetSettings)
				r.Put("/", s.handleUpdateSettings)
			})
		})
	})
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Health endpoints
func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) handleReadyz(w http.ResponseWriter, r *http.Request) {
	// TODO: Check database, Redis, MinIO connectivity
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Prometheus metrics
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("# Metrics endpoint\n"))
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req services.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	resp, err := s.authService.Login(req)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "Invalid credentials"})
			return
		}
		s.logger.WithError(err).Error("Login failed")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Internal server error"})
		return
	}

	render.JSON(w, r, resp)
}

func (s *Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (s *Server) handleListAdapters(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0
	// TODO: Parse limit/offset from query params

	adapters, err := s.adapterService.List(limit, offset)
	if err != nil {
		s.logger.WithError(err).Error("Failed to list adapters")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to list adapters"})
		return
	}

	render.JSON(w, r, adapters)
}

func (s *Server) handleCreateAdapter(w http.ResponseWriter, r *http.Request) {
	var adapter domain.Adapter
	if err := json.NewDecoder(r.Body).Decode(&adapter); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	if err := s.adapterService.Create(&adapter); err != nil {
		s.logger.WithError(err).Error("Failed to create adapter")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to create adapter"})
		return
	}

	userID, _ := mw.GetUserID(r.Context())
	s.auditService.Log(&userID, "adapter.create", "adapter", &adapter.ID, adapter, r.RemoteAddr)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, adapter)
}

func (s *Server) handleGetAdapter(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid adapter ID"})
		return
	}

	adapter, err := s.adapterService.GetByID(id)
	if err != nil {
		if err == services.ErrAdapterNotFound {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "Adapter not found"})
			return
		}
		s.logger.WithError(err).Error("Failed to get adapter")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to get adapter"})
		return
	}

	render.JSON(w, r, adapter)
}

func (s *Server) handleUpdateAdapter(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid adapter ID"})
		return
	}

	var adapter domain.Adapter
	if err := json.NewDecoder(r.Body).Decode(&adapter); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	adapter.ID = id
	if err := s.adapterService.Update(&adapter); err != nil {
		s.logger.WithError(err).Error("Failed to update adapter")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to update adapter"})
		return
	}

	userID, _ := mw.GetUserID(r.Context())
	s.auditService.Log(&userID, "adapter.update", "adapter", &id, adapter, r.RemoteAddr)

	render.JSON(w, r, adapter)
}

func (s *Server) handleDeleteAdapter(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid adapter ID"})
		return
	}

	if err := s.adapterService.Delete(id); err != nil {
		s.logger.WithError(err).Error("Failed to delete adapter")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to delete adapter"})
		return
	}

	userID, _ := mw.GetUserID(r.Context())
	s.auditService.Log(&userID, "adapter.delete", "adapter", &id, nil, r.RemoteAddr)

	render.Status(r, http.StatusNoContent)
}

func (s *Server) handlePublishAdapter(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement adapter publishing logic
	render.Status(r, http.StatusNotImplemented)
	render.JSON(w, r, map[string]string{"error": "Not implemented"})
}

func (s *Server) handleCreateDraft(w http.ResponseWriter, r *http.Request) {
	var draft domain.AdapterDraft
	if err := json.NewDecoder(r.Body).Decode(&draft); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	userID, _ := mw.GetUserID(r.Context())
	draft.CreatedBy = userID

	if err := s.draftService.Create(&draft); err != nil {
		s.logger.WithError(err).Error("Failed to create draft")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to create draft"})
		return
	}

	s.auditService.Log(&userID, "draft.create", "adapter_draft", &draft.ID, draft, r.RemoteAddr)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, draft)
}

func (s *Server) handleGetDraft(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid draft ID"})
		return
	}

	draft, err := s.draftService.GetByID(id)
	if err != nil {
		if err == services.ErrAdapterDraftNotFound {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "Draft not found"})
			return
		}
		s.logger.WithError(err).Error("Failed to get draft")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to get draft"})
		return
	}

	render.JSON(w, r, draft)
}

func (s *Server) handleUpdateDraft(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid draft ID"})
		return
	}

	var draft domain.AdapterDraft
	if err := json.NewDecoder(r.Body).Decode(&draft); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	draft.ID = id
	if err := s.draftService.Update(&draft); err != nil {
		s.logger.WithError(err).Error("Failed to update draft")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to update draft"})
		return
	}

	userID, _ := mw.GetUserID(r.Context())
	s.auditService.Log(&userID, "draft.update", "adapter_draft", &id, draft, r.RemoteAddr)

	render.JSON(w, r, draft)
}

func (s *Server) handleProbeDraft(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	_, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid draft ID"})
		return
	}

	// TODO: Enqueue probe job to Redis queue
	// For now, return a placeholder response
	render.JSON(w, r, map[string]string{
		"status": "queued",
		"job_id": uuid.New().String(),
		"message": "Probe job queued",
	})
}

func (s *Server) handleValidateDraft(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid draft ID"})
		return
	}

	draft, err := s.draftService.GetByID(id)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "Draft not found"})
		return
	}

	// Parse target maturity level from query or body
	targetLevel := domain.MaturityL2 // Default
	if levelStr := r.URL.Query().Get("level"); levelStr != "" {
		targetLevel = domain.MaturityLevel(levelStr)
	}

	// TODO: Enqueue validation job to Redis queue
	// For now, run validation synchronously
	result, err := s.validationService.Validate(draft.ID.String(), targetLevel)
	if err != nil {
		s.logger.WithError(err).Error("Validation failed")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Validation failed"})
		return
	}

	render.JSON(w, r, result)
}

func (s *Server) handleOptimizeDraft(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid draft ID"})
		return
	}

	var req struct {
		SQLTemplate string            `json:"sql_template"`
		Context     map[string]string `json:"context"`
		ProviderID  string            `json:"provider_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	userID, _ := mw.GetUserID(r.Context())

	aiProviderRepo := postgres.NewAIProviderRepository(s.db)
	aiInteractionRepo := postgres.NewAIInteractionRepository(s.db)
	optimizationService := services.NewAIOptimizationService(aiProviderRepo, aiInteractionRepo)
	optReq := &services.OptimizationRequest{
		DraftID:     id.String(),
		SQLTemplate: req.SQLTemplate,
		Context:     req.Context,
		ProviderID:  req.ProviderID,
	}

	result, err := optimizationService.Optimize(r.Context(), optReq, userID.String())
	if err != nil {
		s.logger.WithError(err).Error("Optimization failed")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Optimization failed"})
		return
	}

	s.auditService.Log(&userID, "draft.optimize", "adapter_draft", &id, result, r.RemoteAddr)

	render.JSON(w, r, result)
}

func (s *Server) handlePublishDraft(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid draft ID"})
		return
	}

	draft, err := s.draftService.GetByID(id)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "Draft not found"})
		return
	}

	// Parse target maturity level
	var req struct {
		MaturityLevel domain.MaturityLevel `json:"maturity_level"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.MaturityLevel = domain.MaturityL2 // Default
	}

	// Validate maturity level gates
	// TODO: Get adapter if draft.AdapterID is set
	// For now, skip gate check for new adapters

	// Package bundle
	// TODO: Collect files from draft
	files := []services.BundleFile{
		{Path: "manifest.yaml", Data: []byte("placeholder")},
	}
	bundleData, signature, err := s.bundleService.PackageBundle(files)
	if err != nil {
		s.logger.WithError(err).Error("Failed to package bundle")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to package bundle"})
		return
	}

	// Upload to MinIO
	objectName := fmt.Sprintf("adapters/%s/%s-%s.tar.gz", draft.ID.String(), draft.Name, uuid.New().String())
	ctx := r.Context()
	if err := s.bundleService.UploadBundle(ctx, objectName, bundleData); err != nil {
		s.logger.WithError(err).Error("Failed to upload bundle")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to upload bundle"})
		return
	}

	// Create adapter record
	adapter := &domain.Adapter{
		ID:            uuid.New(),
		Name:          draft.Name,
		Vendor:        "pkfk-discovery", // TODO: Get from draft
		DBFamily:      "postgresql",      // TODO: Get from draft
		Version:       "1.0.0",           // TODO: Get from draft
		MaturityLevel: req.MaturityLevel,
		BundlePath:    objectName,
		Signature:     signature,
	}

	if err := s.adapterService.Create(adapter); err != nil {
		s.logger.WithError(err).Error("Failed to create adapter")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to create adapter"})
		return
	}

	// Update draft
	draft.AdapterID = &adapter.ID
	draft.Status = "published"
	if err := s.draftService.Update(draft); err != nil {
		s.logger.WithError(err).Error("Failed to update draft")
	}

	userID, _ := mw.GetUserID(r.Context())
	s.auditService.Log(&userID, "draft.publish", "adapter_draft", &id, adapter, r.RemoteAddr)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, adapter)
}

func (s *Server) handleListConnections(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0
	// TODO: Parse limit/offset from query params

	conns, err := s.connectionService.List(limit, offset)
	if err != nil {
		s.logger.WithError(err).Error("Failed to list connections")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to list connections"})
		return
	}

	render.JSON(w, r, conns)
}

func (s *Server) handleCreateConnection(w http.ResponseWriter, r *http.Request) {
	var conn domain.Connection
	if err := json.NewDecoder(r.Body).Decode(&conn); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	userID, _ := mw.GetUserID(r.Context())
	conn.ID = uuid.New()
	conn.CreatedBy = userID

	if err := s.connectionService.Create(&conn); err != nil {
		s.logger.WithError(err).Error("Failed to create connection")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to create connection"})
		return
	}

	// Audit log
	s.auditService.Log(&userID, "connection.create", "connection", &conn.ID, conn, r.RemoteAddr)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, conn)
}

func (s *Server) handleGetConnection(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid connection ID"})
		return
	}

	conn, err := s.connectionService.GetByID(id)
	if err != nil {
		if err == services.ErrConnectionNotFound {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "Connection not found"})
			return
		}
		s.logger.WithError(err).Error("Failed to get connection")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to get connection"})
		return
	}

	render.JSON(w, r, conn)
}

func (s *Server) handleUpdateConnection(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid connection ID"})
		return
	}

	var conn domain.Connection
	if err := json.NewDecoder(r.Body).Decode(&conn); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	conn.ID = id
	if err := s.connectionService.Update(&conn); err != nil {
		s.logger.WithError(err).Error("Failed to update connection")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to update connection"})
		return
	}

	userID, _ := mw.GetUserID(r.Context())
	s.auditService.Log(&userID, "connection.update", "connection", &id, conn, r.RemoteAddr)

	render.JSON(w, r, conn)
}

func (s *Server) handleDeleteConnection(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid connection ID"})
		return
	}

	if err := s.connectionService.Delete(id); err != nil {
		s.logger.WithError(err).Error("Failed to delete connection")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to delete connection"})
		return
	}

	userID, _ := mw.GetUserID(r.Context())
	s.auditService.Log(&userID, "connection.delete", "connection", &id, nil, r.RemoteAddr)

	render.Status(r, http.StatusNoContent)
}

func (s *Server) handleTestConnection(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid connection ID"})
		return
	}

	if err := s.connectionService.TestConnection(id); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.JSON(w, r, map[string]string{"status": "success"})
}

func (s *Server) handleListScans(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0
	// TODO: Parse limit/offset from query params

	scans, err := s.scanService.List(limit, offset)
	if err != nil {
		s.logger.WithError(err).Error("Failed to list scans")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to list scans"})
		return
	}

	render.JSON(w, r, scans)
}

func (s *Server) handleCreateScan(w http.ResponseWriter, r *http.Request) {
	var scan domain.Scan
	if err := json.NewDecoder(r.Body).Decode(&scan); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	userID, _ := mw.GetUserID(r.Context())
	scan.CreatedBy = userID

	if err := s.scanService.Create(&scan); err != nil {
		s.logger.WithError(err).Error("Failed to create scan")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to create scan"})
		return
	}

	// TODO: Enqueue scan job to Redis queue

	s.auditService.Log(&userID, "scan.create", "scan", &scan.ID, scan, r.RemoteAddr)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, scan)
}

func (s *Server) handleGetScan(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid scan ID"})
		return
	}

	scan, err := s.scanService.GetByID(id)
	if err != nil {
		if err == services.ErrScanNotFound {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "Scan not found"})
			return
		}
		s.logger.WithError(err).Error("Failed to get scan")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to get scan"})
		return
	}

	render.JSON(w, r, scan)
}

func (s *Server) handleGetScanResults(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid scan ID"})
		return
	}

	scan, err := s.scanService.GetByID(id)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get scan")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to get scan"})
		return
	}

	if scan.Results == nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "Scan results not available"})
		return
	}

	var results interface{}
	if err := json.Unmarshal(scan.Results, &results); err != nil {
		s.logger.WithError(err).Error("Failed to unmarshal scan results")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to parse scan results"})
		return
	}

	render.JSON(w, r, results)
}

func (s *Server) handleCancelScan(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid scan ID"})
		return
	}

	if err := s.scanService.UpdateStatus(id, domain.ScanStatusCancelled); err != nil {
		s.logger.WithError(err).Error("Failed to cancel scan")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to cancel scan"})
		return
	}

	// TODO: Cancel job in Redis queue

	userID, _ := mw.GetUserID(r.Context())
	s.auditService.Log(&userID, "scan.cancel", "scan", &id, nil, r.RemoteAddr)

	render.Status(r, http.StatusNoContent)
}

func (s *Server) handleListAIProviders(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0
	// TODO: Parse limit/offset from query params

	providers, err := s.aiProviderService.List(limit, offset)
	if err != nil {
		s.logger.WithError(err).Error("Failed to list AI providers")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to list AI providers"})
		return
	}

	render.JSON(w, r, providers)
}

func (s *Server) handleCreateAIProvider(w http.ResponseWriter, r *http.Request) {
	var provider domain.AIProvider
	if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	if err := s.aiProviderService.Create(&provider); err != nil {
		s.logger.WithError(err).Error("Failed to create AI provider")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to create AI provider"})
		return
	}

	provider.EncryptedAPIKey = "[REDACTED]"
	userID, _ := mw.GetUserID(r.Context())
	s.auditService.Log(&userID, "ai_provider.create", "ai_provider", &provider.ID, provider, r.RemoteAddr)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, provider)
}

func (s *Server) handleGetAIProvider(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid AI provider ID"})
		return
	}

	provider, err := s.aiProviderService.GetByID(id)
	if err != nil {
		if err == services.ErrAIProviderNotFound {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "AI provider not found"})
			return
		}
		s.logger.WithError(err).Error("Failed to get AI provider")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to get AI provider"})
		return
	}

	render.JSON(w, r, provider)
}

func (s *Server) handleUpdateAIProvider(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid AI provider ID"})
		return
	}

	var provider domain.AIProvider
	if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	provider.ID = id
	if err := s.aiProviderService.Update(&provider); err != nil {
		s.logger.WithError(err).Error("Failed to update AI provider")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to update AI provider"})
		return
	}

	provider.EncryptedAPIKey = "[REDACTED]"
	userID, _ := mw.GetUserID(r.Context())
	s.auditService.Log(&userID, "ai_provider.update", "ai_provider", &id, provider, r.RemoteAddr)

	render.JSON(w, r, provider)
}

func (s *Server) handleDeleteAIProvider(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid AI provider ID"})
		return
	}

	if err := s.aiProviderService.Delete(id); err != nil {
		s.logger.WithError(err).Error("Failed to delete AI provider")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to delete AI provider"})
		return
	}

	userID, _ := mw.GetUserID(r.Context())
	s.auditService.Log(&userID, "ai_provider.delete", "ai_provider", &id, nil, r.RemoteAddr)

	render.Status(r, http.StatusNoContent)
}

func (s *Server) handleListUsers(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0
	// TODO: Parse limit/offset from query params

	users, err := s.userService.List(limit, offset)
	if err != nil {
		s.logger.WithError(err).Error("Failed to list users")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to list users"})
		return
	}

	render.JSON(w, r, users)
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	if err := s.userService.Create(&user); err != nil {
		s.logger.WithError(err).Error("Failed to create user")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to create user"})
		return
	}

	user.PasswordHash = "" // Don't return password hash
	userID, _ := mw.GetUserID(r.Context())
	s.auditService.Log(&userID, "user.create", "user", &user.ID, user, r.RemoteAddr)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid user ID"})
		return
	}

	user, err := s.userService.GetByID(id)
	if err != nil {
		if err == services.ErrUserNotFound {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "User not found"})
			return
		}
		s.logger.WithError(err).Error("Failed to get user")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to get user"})
		return
	}

	user.PasswordHash = "" // Don't return password hash
	render.JSON(w, r, user)
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid user ID"})
		return
	}

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	user.ID = id
	if err := s.userService.Update(&user); err != nil {
		s.logger.WithError(err).Error("Failed to update user")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to update user"})
		return
	}

	user.PasswordHash = ""
	userID, _ := mw.GetUserID(r.Context())
	s.auditService.Log(&userID, "user.update", "user", &id, user, r.RemoteAddr)

	render.JSON(w, r, user)
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid user ID"})
		return
	}

	if err := s.userService.Delete(id); err != nil {
		s.logger.WithError(err).Error("Failed to delete user")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to delete user"})
		return
	}

	userID, _ := mw.GetUserID(r.Context())
	s.auditService.Log(&userID, "user.delete", "user", &id, nil, r.RemoteAddr)

	render.Status(r, http.StatusNoContent)
}

func (s *Server) handleListAuditLogs(w http.ResponseWriter, r *http.Request) {
	filters := domain.AuditLogFilters{}
	// TODO: Parse filters from query params

	limit := 50
	offset := 0
	// TODO: Parse limit/offset from query params

	logs, err := s.auditService.List(filters, limit, offset)
	if err != nil {
		s.logger.WithError(err).Error("Failed to list audit logs")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to list audit logs"})
		return
	}

	render.JSON(w, r, logs)
}

func (s *Server) handleGetAuditLog(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid audit log ID"})
		return
	}

	log, err := s.auditService.GetByID(id)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get audit log")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to get audit log"})
		return
	}

	render.JSON(w, r, log)
}

func (s *Server) handleExportAuditLogs(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement CSV/JSON export
	render.Status(r, http.StatusNotImplemented)
	render.JSON(w, r, map[string]string{"error": "Not implemented"})
}

func (s *Server) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := s.settingsService.GetAll()
	if err != nil {
		s.logger.WithError(err).Error("Failed to get settings")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to get settings"})
		return
	}

	render.JSON(w, r, settings)
}

func (s *Server) handleUpdateSettings(w http.ResponseWriter, r *http.Request) {
	var updates map[string]map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	for key, value := range updates {
		if err := s.settingsService.Set(key, value); err != nil {
			s.logger.WithError(err).Errorf("Failed to update setting: %s", key)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to update settings"})
			return
		}
	}

	userID, _ := mw.GetUserID(r.Context())
	s.auditService.Log(&userID, "settings.update", "settings", nil, updates, r.RemoteAddr)

	render.JSON(w, r, map[string]string{"status": "success"})
}

