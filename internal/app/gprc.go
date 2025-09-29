package app

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-geo-api/internal/config"
	
	handlers "grpc-geo-api/handlers"
	healthpb "grpc-geo-api/gen/go/health/v1"
)

type Server struct {
	grpcServer *grpc.Server
	listener   	net.Listener
	address		string
}

func NewServer(cfg config.GRPCConfig) (*Server, error) {
	lis, err := net.Listen("tcp", cfg.Address())
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", cfg.Address(), err)
	}

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	return &Server{
		grpcServer: grpcServer,
		listener: 	lis,
		address:  	cfg.Address(),
	}, nil
}


func (s *Server) Start() error {
	return s.grpcServer.Serve(s.listener)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}

func (s *Server) RegisterServices() {
	healthHandler := handlers.NewHealthHandler()
	
	healthpb.RegisterHealthServer(s.grpcServer, healthHandler)
}