package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httptransport "github.com/pkfk-discovery/api/internal/transport/http"
)

func main() {
	// Load configuration from environment
	cfg := loadConfig()

	// Initialize server
	server, err := httptransport.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting API server on %s", cfg.Addr)
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func loadConfig() *httptransport.Config {
	return &httptransport.Config{
		Addr:         getEnv("ADDR", ":8080"),
		DatabaseURL:   getEnv("DATABASE_URL", ""),
		RedisURL:     getEnv("REDIS_URL", ""),
		MinIOEndpoint: getEnv("MINIO_ENDPOINT", ""),
		MinIOAccessKey: getEnv("MINIO_ACCESS_KEY", ""),
		MinIOSecretKey: getEnv("MINIO_SECRET_KEY", ""),
		MinIOUseSSL:   getEnv("MINIO_USE_SSL", "false") == "true",
		JWTSecret:    getEnv("JWT_SECRET", ""),
		EncryptionKey: getEnv("ENCRYPTION_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

