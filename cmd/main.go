package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	_ "grpc-geo-api/docs"

	runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpcDelivery "grpc-geo-api/delivery/grpc"
	healthpb "grpc-geo-api/gen/go/health/v1"
)

func startGrpcServer() {
	grpcSrv, err := grpcDelivery.NewServer(":50051")
	if err != nil {
		log.Fatalf("failed to create gRPC server: %v", err)
	}

	healthHandler := grpcDelivery.NewHealthHandler()
	grpcSrv.RegisterService(healthHandler)

	log.Println("Starting gRPC server on :50051")
	if err := grpcSrv.Start(); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}

func startGateway() {
	ctx := context.Background()

	gatewayMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := healthpb.RegisterHealthHandlerFromEndpoint(ctx, gatewayMux, ":50051", opts); err != nil {
		log.Fatalf("failed to register health handler: %v", err)
	}

	r := gin.Default()
	api := r.Group("/api")
	api.Any("/*path", gin.WrapH(gatewayMux))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Ok"})
	})

	log.Println("Starting HTTP gateway on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("HTTP gateway error: %v", err)
	}
}

func main() {
	go startGrpcServer()
	startGateway()
}
