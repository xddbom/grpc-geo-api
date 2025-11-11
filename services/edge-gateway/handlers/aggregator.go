package handlers

import (
	"context"
	aggregatorv1 "github.com/xddbom/grpc-geo-api/api/gen/go/aggregate/v1"
)

type AggregatorHandler struct {
	aggregatorv1.UnimplementedAggregatorServer
}

func NewAggregatorHandler() *AggregatorHandler {
	return &AggregatorHandler{}
}

func (h *AggregatorHandler) Aggregate(ctx context.Context, req *aggregatorv1.AggregateRequest) (*aggregatorv1.AggregateResponse, error) {
	return &aggregatorv1.AggregateResponse{
		Geo:     "Kyiv, Ukraine",
		Weather: "Sunny 20°C",
	}, nil
}
