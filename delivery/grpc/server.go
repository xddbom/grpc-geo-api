package grpc

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	healthpb "grpc-geo-api/gen/go/health/v1"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

func NewServer(addr string) (*Server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	return &Server{
		grpcServer: grpcServer,
		listener:   lis,
	}, nil
}

func (s *Server) RegisterService(handler *HealthHandler) { // TODO: it's not a health service (custom interface needed)
	healthpb.RegisterHealthServer(s.grpcServer, handler)
}

func (s *Server) Start() error {
	return s.grpcServer.Serve(s.listener)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}
