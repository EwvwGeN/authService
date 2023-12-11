package grpcHandler

import (
	"context"
	"fmt"
	"log/slog"

	c "github.com/EwvwGeN/authService/internal/config"
	"github.com/EwvwGeN/authService/internal/validator"
	authProto "github.com/EwvwGeN/authService/proto/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	log *slog.Logger
	validator c.Validator
	authProto.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server, validator c.Validator, log *slog.Logger) {
	authProto.RegisterAuthServer(gRPC, &Server{
		log: log,
		validator: validator,
	})
}

func (s *Server) Login(
	ctx context.Context,
	req *authProto.LoginRequest,
	) (*authProto.LoginResponse, error) {
		if !validator.ValideteByRegex(req.GetEmail(), s.validator.EmailValidate) {
			return nil, status.Error(codes.InvalidArgument, "incorrect email")
		}
		if !validator.ValideteByRegex(req.GetPassword(), s.validator.PasswordValidate) {
			return nil, status.Error(codes.InvalidArgument, "incorrect password")
		}
		if !validator.ValideteByRegex(fmt.Sprintf("%d",req.GetAppId()), s.validator.AppIDValidate) {
			return nil, status.Error(codes.InvalidArgument, "incorrect app id")
		}
	return &authProto.LoginResponse{
		Token: "123",
	}, nil
}

func (s *Server) Register(
	ctx context.Context,
	req *authProto.RegisterRequest,
	) (*authProto.RegisterResponse, error) {
		if !validator.ValideteByRegex(req.GetEmail(), s.validator.EmailValidate) {
			return nil, status.Error(codes.InvalidArgument, "incorrect email")
		}
		if !validator.ValideteByRegex(req.GetPassword(), s.validator.PasswordValidate) {
			return nil, status.Error(codes.InvalidArgument, "incorrect password")
		}
	return &authProto.RegisterResponse{
		UserId: 1,
	}, nil
}

func (s *Server) IsAdmin(
	ctx context.Context,
	req *authProto.IsAdminRequest,
	) (*authProto.IsAdminResponse, error) {
		if !validator.ValideteByRegex(fmt.Sprintf("%d",req.GetUserId()), s.validator.EmailValidate) {
			return nil, status.Error(codes.InvalidArgument, "incorrect user id")
		}
	return &authProto.IsAdminResponse{
		IsAdmin: true,
	}, nil
}