package app

import (
    "context"
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    healthpb "grpc-geo-api/gen/go/health/v1"
    config "grpc-geo-api/internal/config"
)

type Gateway struct {
	server *http.Server
	address string
}

func NewGateway(cfg config.Config) (*Gateway, error) {
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

    return &Gateway{
        server:  server,
        address: cfg.Gateway.Address(),
    }, nil
}
	
func (g *Gateway) Start() error {
    log.Printf("Starting HTTP gateway on %s", g.address)
    return g.server.ListenAndServe()
}

func (g *Gateway) Stop() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := g.server.Shutdown(ctx); err != nil {
        log.Printf("Gateway shutdown error: %v", err)
    }
}