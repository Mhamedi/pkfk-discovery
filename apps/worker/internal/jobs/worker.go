package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/pkfk-discovery/worker/internal/domain"
	"github.com/pkfk-discovery/worker/internal/integrations/minio"
	"github.com/pkfk-discovery/worker/internal/repositories/postgres"
	"github.com/pkfk-discovery/worker/internal/services"
)
type Config struct {
	DatabaseURL   string
	RedisURL      string
	MinIOEndpoint string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOUseSSL   bool
	EncryptionKey string
}
type Worker struct {
	config        *Config
	client        *asynq.Client
	server        *asynq.Server
	mux           *asynq.ServeMux
	logger        *logrus.Logger
	db            *pgxpool.Pool
	minioClient   *minio.Client
	scanPipeline  *ScanPipeline
	scanRepo      domain.ScanRepository
	adapterRepo   domain.AdapterRepository
	connectionRepo domain.ConnectionRepository
	sqlSafety     *services.SQLSafety
}
func NewWorker(cfg *Config) (*Worker, error) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	// Initialize database connection
	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	// Initialize MinIO client
	minioClient, err := minio.NewClient(
		cfg.MinIOEndpoint,
		cfg.MinIOAccessKey,
		cfg.MinIOSecretKey,
		cfg.MinIOUseSSL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}
	// Initialize SQL safety
	sqlSafety := services.NewSQLSafety()
	// Initialize repositories
	scanRepo := postgres.NewScanRepository(db)
	adapterRepo := postgres.NewAdapterRepository(db)
	connectionRepo := postgres.NewConnectionRepository(db)
	// Initialize scan pipeline
	scanPipeline := NewScanPipeline(db, minioClient, logger, sqlSafety)
	redisOpt := asynq.RedisClientOpt{
		Addr: cfg.RedisURL,
	}
	client := asynq.NewClient(redisOpt)
	server := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":       1,
		},
	})
	mux := asynq.NewServeMux()
	w := &Worker{
		config:        cfg,
		client:        client,
		server:        server,
		mux:           mux,
		logger:        logger,
		db:            db,
		minioClient:   minioClient,
		scanPipeline:  scanPipeline,
		scanRepo:      scanRepo,
		adapterRepo:   adapterRepo,
		connectionRepo: connectionRepo,
		sqlSafety:     sqlSafety,
	}
	w.registerHandlers()
	return w, nil
}
func (w *Worker) registerHandlers() {
	w.mux.HandleFunc("adapter_probe", w.handleAdapterProbe)
	w.mux.HandleFunc("adapter_validate", w.handleAdapterValidate)
	w.mux.HandleFunc("scan_run", w.handleScanRun)
}
func (w *Worker) Start() error {
	return w.server.Run(w.mux)
}
func (w *Worker) Shutdown(ctx context.Context) error {
	w.server.Shutdown()
	w.client.Close()
	if w.db != nil {
		w.db.Close()
	}
	return nil
}
func (w *Worker) handleAdapterProbe(ctx context.Context, t *asynq.Task) error {
	// TODO: Implement adapter probe logic
	w.logger.Info("Processing adapter_probe job")
	return nil
}
func (w *Worker) handleAdapterValidate(ctx context.Context, t *asynq.Task) error {
	// TODO: Implement adapter validation logic
	w.logger.Info("Processing adapter_validate job")
	return nil
}
func (w *Worker) handleScanRun(ctx context.Context, t *asynq.Task) error {
	var payload ScanJobPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal scan job payload: %w", err)
	}
	scanID := payload.ScanID
	w.logger.WithField("scan_id", scanID).Info("Processing scan_run job")
	// Load scan from database
	scan, err := w.scanRepo.GetByID(scanID)
	if err != nil {
		return fmt.Errorf("failed to load scan: %w", err)
	}
	// Update status to running
	scan.Status = domain.ScanStatusRunning
	if err := w.scanRepo.Update(scan); err != nil {
		w.logger.WithError(err).Error("Failed to update scan status")
	}
	// Load adapter
	adapter, err := w.adapterRepo.GetByID(scan.AdapterID)
	if err != nil {
		scan.Status = domain.ScanStatusFailed
		w.scanRepo.Update(scan)
		return fmt.Errorf("failed to load adapter: %w", err)
	}
	// Load connection
	connection, err := w.connectionRepo.GetByID(scan.ConnectionID)
	if err != nil {
		scan.Status = domain.ScanStatusFailed
		w.scanRepo.Update(scan)
		return fmt.Errorf("failed to load connection: %w", err)
	}
	// Parse policy
	var policy domain.ScanPolicy
	if err := json.Unmarshal(scan.Policy, &policy); err != nil {
		policy = domain.ScanPolicy{
			SampleMode:  true,
			DeepMode:    false,
			Timeout:     300,
			MaxRows:     10000,
			Concurrency: 5,
		}
	}
	// Run scan pipeline
	results, err := w.scanPipeline.Run(ctx, scanID, adapter, connection, policy)
	if err != nil {
		scan.Status = domain.ScanStatusFailed
		w.scanRepo.Update(scan)
		return fmt.Errorf("scan pipeline failed: %w", err)
	}
	// Save results
	resultsJSON, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}
	scan.Results = resultsJSON
	scan.Status = domain.ScanStatusCompleted
	if err := w.scanRepo.Update(scan); err != nil {
		return fmt.Errorf("failed to save results: %w", err)
	}
	w.logger.WithField("scan_id", scanID).Info("Scan completed successfully")
	return nil
}
