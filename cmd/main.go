package main

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/go-pricing-engine/internal/api"
	"github.com/go-pricing-engine/internal/config"
	"github.com/go-pricing-engine/internal/logger"
)

func main() {
	logger.InitLogger()
	defer logger.CloseLogger()

	// Load .env file
	if err := godotenv.Load(); err != nil {
		logger.Logger.Fatal("Error loading .env file", zap.Error(err))
	}

	// Load config
	cfg := config.NewConfig()
	if err := cfg.ParseFlags(); err != nil {
		logger.Logger.Fatal("Failed to parse command-line flags", zap.Error(err))
	}

	// Start server
	srv := api.NewServer(logger.Logger, cfg)
	err := srv.Run()
	if err != nil {
		logger.Logger.Fatal("Failed to start server", zap.Error(err))
	}
}
