package app

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/xddbom/grpc-geo-api/services/edge-gateway/internal/config"
	"go.uber.org/zap"
)

type App struct {
	grpcServer *Server
	gateway    *Gateway
	config     config.Config
	logger     *zap.Logger
}

func New(cfg config.Config, logger *zap.Logger) (*App, error) {
	grpcServer, err := NewServer(cfg.GRPC, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC server: %w", err)
	}
	grpcServer.RegisterServices()

	gateway, err := NewGateway(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create gateway: %w", err)
	}

	return &App{
		grpcServer: grpcServer,
		gateway:    gateway,
		config:     cfg,
		logger:     logger,
	}, nil
}

func (a *App) Run() error {
	errChan := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		a.logger.Info("Starting gRPC server", zap.String("address", a.config.GRPC.Address()))
		if err := a.grpcServer.Start(); err != nil {
			a.logger.Error("gRPC server failed", zap.Error(err))
			errChan <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		a.logger.Info("Starting HTTP gateway", zap.String("address", a.config.Gateway.Address()))
		if err := a.gateway.Start(); err != nil {
			a.logger.Error("Gateway failed", zap.Error(err))
			errChan <- fmt.Errorf("gateway error: %w", err)
		}
	}()

	return a.waitForShutdown(errChan, &wg)
}

func (a *App) waitForShutdown(errChan chan error, wg *sync.WaitGroup) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		a.logger.Error("Server error occurred", zap.Error(err))
		a.Stop()
		return err
	case sig := <-sigChan:
		a.logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
		a.Stop()
		wg.Wait()
		return nil
	}
}

func (a *App) Stop() {
	a.logger.Info("Initiating graceful shutdown")

	a.logger.Debug("Stopping gateway")
	a.gateway.Stop()

	a.logger.Debug("Stopping gRPC server")
	a.grpcServer.Stop()

	a.logger.Info("Application stopped gracefully")
}
