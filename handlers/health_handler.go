package grpc

import (
	"context"

	healthpb "grpc-geo-api/gen/go/health/v1"
)

type HealthHandler struct {
	healthpb.UnimplementedHealthServer
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(ctx context.Context, req *healthpb.HealthRequest) (*healthpb.HealthResponse, error) {
	return &healthpb.HealthResponse{
		Status: "OK",
	}, nil
}
