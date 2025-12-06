package grpc

import (
	"context"
	"net"

	"go-templ/internal/resource"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	config *resource.Config
	server *grpc.Server
	logger *zap.Logger
	errors chan error
}

func NewServer(cfg *resource.Config, log *zap.Logger) *Server {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor(log)),
	)

	return &Server{
		config: cfg,
		server: grpcServer,
		logger: log,
	}
}

func (s *Server) Notify() <-chan error {
	return s.errors
}

func (s *Server) Start() {
	go func() {
		lis, err := net.Listen("tcp", net.JoinHostPort("", s.config.ApiPort))
		if err != nil {
			s.errors <- err
		}

		s.errors <- s.server.Serve(lis)
		close(s.errors)
	}()
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}

func loggingInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log.Info("gRPC request", zap.String("method", info.FullMethod))
		return handler(ctx, req)
	}
}
