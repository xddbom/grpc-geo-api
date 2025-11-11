package main

import (
	"log"

	"github.com/xddbom/grpc-geo-api/internal/platform/logging"
	"github.com/xddbom/grpc-geo-api/services/edge-gateway/internal/app"
	"github.com/xddbom/grpc-geo-api/services/edge-gateway/internal/config"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	logger, err := logging.New(logging.Config{
		Level:       cfg.Log.Level,
		Environment: cfg.Log.Environment,
	})
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	logger.Info("Service started successfully")
	defer logger.Sync()

	application, err := app.New(*cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize application", zap.Error(err))
	}

	if err := application.Run(); err != nil {
		logger.Fatal("Application runtime error", zap.Error(err))
	}

	logger.Info("Application exited successfully")
}
