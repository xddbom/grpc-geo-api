package app

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/xddbom/grpc-geo-api/services/edge-gateway/internal/config"
)

type Gateway struct {
	server  *http.Server
	address string
	logger  *zap.Logger
	opts    []grpc.DialOption
}

func NewGateway(cfg config.Config, logger *zap.Logger) (*Gateway, error) {

	mux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	server := &http.Server{
		Addr:    cfg.Gateway.Address(),
		Handler: mux,
	}

	gatewayLogger := logger.With(zap.String("component", "gateway"))

	return &Gateway{
		server:  server,
		address: cfg.Gateway.Address(),
		logger:  gatewayLogger,
		opts: 	 opts,
	}, nil
}

func (g *Gateway) Start() error {
	g.logger.Info("Starting HTTP gateway", zap.String("address", g.address))
	if err := g.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		g.logger.Error("Gateway stopped with error", zap.Error(err))
		return err
	}
	g.logger.Info("Gateway stopped gracefully")
	return nil
}

func (g *Gateway) Stop() {
	g.logger.Info("Stopping HTTP gateway")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := g.server.Shutdown(ctx); err != nil {
		g.logger.Error("Error during gateway shutdown", zap.Error(err))
	}
	g.logger.Info("HTTP gateway stopped")
}
