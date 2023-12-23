package app

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/EwvwGeN/authService/internal/config"
	"github.com/EwvwGeN/authService/internal/grpcHandler"
	"google.golang.org/grpc"
)

type Server struct {
	log *slog.Logger
	grpcServer *grpc.Server
	port int
}

func NewSerever(
	log *slog.Logger,
	validator config.Validator,
	auth grpcHandler.Auth,
	port int,
) *Server {
	grpcS := grpc.NewServer()

	grpcHandler.Register(grpcS, validator, log, auth)

	return &Server{
		log: log,
		grpcServer: grpcS,
		port: port,
	}
}

func (s *Server) GRPCRun() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("can't create listener: %w", err)
	}
	s.log.Info("grpc server is running", slog.String("addr", l.Addr().String()))
	if err := s.grpcServer.Serve(l); err != nil {
		return fmt.Errorf("can't start grpc server: %w", err)
	}
	return nil
}

func (s *Server) GRPCStop() {
	s.log.Info("stopping grpc server")

	s.grpcServer.GracefulStop()
}