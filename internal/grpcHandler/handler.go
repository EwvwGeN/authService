package grpcHandler

import (
	"context"

	authProto "github.com/EwvwGeN/authService/proto/gen/go"
	"google.golang.org/grpc"
)

type Server struct {
	authProto.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	authProto.RegisterAuthServer(gRPC, &Server{})
}

func (s *Server) Login(
	ctx context.Context,
	req *authProto.LoginRequest,
	) (*authProto.LoginResponse, error) {
	return nil, nil
}

func (s *Server) Register(
	ctx context.Context,
	req *authProto.RegisterRequest,
	) (*authProto.RegisterResponse, error) {
	return nil, nil
}

func (s *Server) IsAdmin(
	ctx context.Context,
	req *authProto.IsAdminRequest,
	) (*authProto.IsAdminResponse, error) {
	return nil, nil
}