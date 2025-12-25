package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkfk-discovery/worker/internal/jobs"
)

func main() {
	// Load configuration from environment
	cfg := loadConfig()

	// Initialize worker
	worker, err := jobs.NewWorker(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize worker: %v", err)
	}

	// Start worker in goroutine
	go func() {
		log.Println("Starting worker...")
		if err := worker.Start(); err != nil {
			log.Fatalf("Worker failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := worker.Shutdown(ctx); err != nil {
		log.Fatalf("Worker forced to shutdown: %v", err)
	}

	log.Println("Worker exited")
}

func loadConfig() *jobs.Config {
	return &jobs.Config{
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		RedisURL:       getEnv("REDIS_URL", ""),
		MinIOEndpoint:  getEnv("MINIO_ENDPOINT", ""),
		MinIOAccessKey: getEnv("MINIO_ACCESS_KEY", ""),
		MinIOSecretKey: getEnv("MINIO_SECRET_KEY", ""),
		MinIOUseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
		EncryptionKey:  getEnv("ENCRYPTION_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
