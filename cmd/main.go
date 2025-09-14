package main

import (
	"fmt"
	"log"

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
}
