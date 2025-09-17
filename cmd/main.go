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
	"github.com/vishalanarase/order-system/internal/kafka"
	"github.com/vishalanarase/order-system/internal/services"
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
	log, err := logger.InitLogger(cfg.Environment)
	if err != nil {
		log.Sugar().Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	log.Info("Starting Order Processing System")

	log.Sugar().Infof("config: %+v", cfg)

	// Create Kafka producer
	producer, err := kafka.NewProducer(cfg.Kafka.Brokers, "orders")
	if err != nil {
		log.Sugar().Fatal("Failed to create Kafka producer", err)
	}
	defer producer.Close()

	// Initialize services
	orderService := services.NewOrderService(log, producer)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start background processors
	orderService.StartOrderProcessor(ctx, cfg.Kafka.Brokers)

	// Setup HTTP server
	router := gin.Default()
	api.RegisterRoutes(log, router, orderService)

	server := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Sugar().Infof("HTTP server starting on port: %s", cfg.HTTP.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Sugar().Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Sugar().Errorf("Server forced to shutdown: ", err)
	}

	cancel()                    // Cancel background processors context
	time.Sleep(1 * time.Second) // Give time for graceful shutdown

	log.Info("Server exited properly")
}
