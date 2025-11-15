package app

import (
	"net"

	"github.com/xddbom/grpc-geo-api/services/edge-gateway/internal/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/reflection"
  
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	grpcServer *grpc.Server
	listener    net.Listener
	address     string
	logger      *zap.Logger
	healthSrv	*health.Server
}

func NewServer(cfg config.GRPCConfig, logger *zap.Logger) (*Server, error) {
	lis, err := net.Listen("tcp", cfg.Address())
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()

	grpcLogger := logger.With(zap.String("component", "grpc"))

	return &Server{
		grpcServer: grpcServer,
		listener:   lis,
		address:    cfg.Address(),
		logger:     grpcLogger,
	}, nil
}

func (s *Server) RegisterServices() {
	s.healthSrv = health.NewServer()
    s.healthSrv.SetServingStatus("grpc.health.v1.Health", healthpb.HealthCheckResponse_SERVING)		// ?
	s.healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)							// ?
	healthpb.RegisterHealthServer(s.grpcServer, s.healthSrv)										// ?

	reflection.Register(s.grpcServer)
	s.logger.Info("gRPC services registered")
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
	s.healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	s.grpcServer.GracefulStop()
	s.logger.Info("gRPC server stopped gracefully")
}
