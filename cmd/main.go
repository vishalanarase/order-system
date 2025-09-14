package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vishalanarase/order-system/internal/api"
	"github.com/vishalanarase/order-system/pkg/config"
	"github.com/vishalanarase/order-system/pkg/logger"
)

func main() {
	fmt.Println("Hello!")
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize logger
	if err := logger.InitLogger(cfg.Environment); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	logger.Log.Info("Starting Order Processing System")

	logger.Log.Sugar().Infof("config: %+v", cfg)

	// Create context for graceful shutdown
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup HTTP server
	router := gin.Default()
	api.RegisterRoutes(router)

	server := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Log.Sugar().Infof("HTTP server starting on port: %s", cfg.HTTP.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Sugar().Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Sugar().Errorf("Server forced to shutdown: ", err)
	}

	cancel()                    // Cancel background processors context
	time.Sleep(1 * time.Second) // Give time for graceful shutdown

	logger.Log.Info("Server exited properly")
}
