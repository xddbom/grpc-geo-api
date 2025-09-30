package app

import (
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-geo-api/internal/config"

	handlers "grpc-geo-api/handlers"
	healthpb "grpc-geo-api/gen/go/health/v1"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
	address    string
	logger     *zap.Logger
}

func NewServer(cfg config.GRPCConfig, logger *zap.Logger) (*Server, error) {
	lis, err := net.Listen("tcp", cfg.Address())
	if err != nil {
		return nil, err 
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	grpcLogger := logger.With(zap.String("component", "grpc"))

	return &Server{
		grpcServer: grpcServer,
		listener:   lis,
		address:    cfg.Address(),
		logger:     grpcLogger,
	}, nil
}

func (s *Server) Start() error {
	s.logger.Info("Starting gRPC server", zap.String("address", s.address))
	err := s.grpcServer.Serve(s.listener)
	if err != nil {
		s.logger.Error("gRPC server stopped with error", zap.Error(err))
		return err
	}
	s.logger.Info("gRPC server stopped")
	return nil
}

func (s *Server) Stop() {
	s.logger.Info("Stopping gRPC server")
	s.grpcServer.GracefulStop()
	s.logger.Info("gRPC server stopped gracefully")
}

func (s *Server) RegisterServices() {
	healthHandler := handlers.NewHealthHandler()
	healthpb.RegisterHealthServer(s.grpcServer, healthHandler)
	s.logger.Info("gRPC services registered")
}
