package main

import (
    "log"

    "go.uber.org/zap"
    "grpc-geo-api/internal/app"
    "grpc-geo-api/internal/config"
)

func main() {
    cfg := config.Load()

    logger, err := app.NewLogger(&cfg.Log)
    if err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }
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