package app

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	healthpb "grpc-geo-api/gen/go/health/v1"
	"grpc-geo-api/internal/config"
)

type Gateway struct {
	server  *http.Server
	address string
	logger  *zap.Logger
}

func NewGateway(cfg config.Config, logger *zap.Logger) (*Gateway, error) {
	ctx := context.Background()
	gatewayMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := healthpb.RegisterHealthHandlerFromEndpoint(ctx, gatewayMux, cfg.GRPC.Address(), opts); err != nil {
		return nil, err
	}

	r := gin.Default()
	api := r.Group("/api")
	api.Any("/*path", gin.WrapH(gatewayMux))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Ok"})
	})

	server := &http.Server{
		Addr:    cfg.Gateway.Address(),
		Handler: r,
	}

	gatewayLogger := logger.With(zap.String("component", "gateway"))

	return &Gateway{
		server:  server,
		address: cfg.Gateway.Address(),
		logger:  gatewayLogger,
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
