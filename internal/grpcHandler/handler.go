package grpcHandler

import (
	"context"
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
	auth Auth
	authProto.UnimplementedAuthServer
}

type Auth interface {
	Login(ctx context.Context,
		email string,
		password string,
		appId string,
	) (token string, err error)
	RegisterNewUser(ctx context.Context,
		email string,
		password string,
	) (userId string, err error)
	IsAdmin(ctx context.Context,
	userId string,
	) (bool, error)
}

func Register(gRPC *grpc.Server, validator c.Validator, log *slog.Logger, auth Auth) {
	authProto.RegisterAuthServer(gRPC, &Server{
		log: log,
		validator: validator,
		auth: auth,
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
		if !validator.ValideteByRegex(req.GetAppId(), s.validator.AppIDValidate) {
			return nil, status.Error(codes.InvalidArgument, "incorrect app id")
		}
		token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
		if err != nil {
			return nil, status.Error(codes.Internal, "cant login user")
		}
		return &authProto.LoginResponse{
			Token: token,
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
		uId, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
		if err != nil {
			return nil, status.Error(codes.Internal, "cant register user")
		}
		return &authProto.RegisterResponse{
			UserId: uId,
		}, nil
}

func (s *Server) IsAdmin(
	ctx context.Context,
	req *authProto.IsAdminRequest,
	) (*authProto.IsAdminResponse, error) {
		if !validator.ValideteByRegex(req.GetUserId(), s.validator.UserIDValidate) {
			return nil, status.Error(codes.InvalidArgument, "incorrect user id")
		}
		isAdm, err := s.auth.IsAdmin(ctx, req.GetUserId())
		if err != nil {
			return nil, status.Error(codes.Internal, "cant check user for admin rights")
		}
		return &authProto.IsAdminResponse{
			IsAdmin: isAdm,
		}, nil
}