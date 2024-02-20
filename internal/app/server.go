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
	cfg config.GRPCConfig
}

func NewSerever(
	log *slog.Logger,
	validator config.Validator,
	cfg config.GRPCConfig,
	auth grpcHandler.Auth,
	msgSender grpcHandler.MessageSender,
) *Server {
	grpcS := grpc.NewServer()

	grpcHandler.Register(grpcS, validator, cfg, log, auth, msgSender)

	return &Server{
		log: log,
		grpcServer: grpcS,
		cfg: cfg,
	}
}

func (s *Server) GRPCRun() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port))
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