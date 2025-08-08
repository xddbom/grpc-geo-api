package main

import (
	"context"
	"log"
	"net"

	_ "grpc-geo-api/docs"

	healthpb "grpc-geo-api/gen/go/health/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type HealthServer struct {
	healthpb.UnimplementedHealthServer
}

func (s *HealthServer) Check(ctx context.Context, req *healthpb.HealthRequest) (*healthpb.HealthResponse, error) {
	return &healthpb.HealthResponse{Status: "OK"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	healthpb.RegisterHealthServer(s, &HealthServer{})

	reflection.Register(s)

	log.Println("Server running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
