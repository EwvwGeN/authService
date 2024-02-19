package grpcHandler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	c "github.com/EwvwGeN/authService/internal/config"
	"github.com/EwvwGeN/authService/internal/domain/models"
	"github.com/EwvwGeN/authService/internal/services/auth"
	"github.com/EwvwGeN/authService/internal/storage"
	tmpl "github.com/EwvwGeN/authService/internal/template"
	"github.com/EwvwGeN/authService/internal/validator"
	"github.com/EwvwGeN/authService/internal/verification"
	authProto "github.com/EwvwGeN/authService/proto/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	log *slog.Logger
	validator c.Validator
	template c.Template
	auth Auth
	msgSender MessageSender
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

type MessageSender interface {
	SendMsg(ctx context.Context,
		msg *models.Message,
	) error
}

func Register(gRPC *grpc.Server, validator c.Validator, template c.Template, log *slog.Logger, auth Auth, msgSender MessageSender) {
	authProto.RegisterAuthServer(gRPC, &Server{
		log: log,
		validator: validator,
		template: template,
		auth: auth,
		msgSender: msgSender,
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
			if errors.Is(err, auth.ErrInvalidCredentials) {
				return nil, status.Error(codes.InvalidArgument, "invalid email or password")
			}
	
			return nil, status.Error(codes.Internal, "failed to login")
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
			if errors.Is(err, storage.ErrUserExist) {
				return nil, status.Error(codes.AlreadyExists, "user already exists")
			}
	
			return nil, status.Error(codes.Internal, "failed to register user")
		}
		go func (uId string){
			vCode, err := verification.GenerateVerificationCode(uId)
			if err != nil {
				s.log.Error("error while creating verification code", slog.String("userId", uId), slog.String("error", err.Error()))
				return
			}
			s.log.Debug("create verification code", slog.String("userId", uId), slog.String("code", vCode))
			link := fmt.Sprintf("%s?verificationCode=%s", s.template.RgistrationLink, verification.ConvertForURL(vCode))
			msgBody, err := tmpl.Register(link)
			s.log.Debug("created message", slog.String("userId", uId), slog.String("msg", string(msgBody)))
			if err != nil {
				s.log.Error("error while parsing template", slog.String("userId", uId), slog.String("error", err.Error()))
				return
			}
			err = s.msgSender.SendMsg(ctx, &models.Message{
				Subject: "Register msg",
				EmailTo: req.GetEmail(),
				Body: msgBody,
			})
			if err != nil {
				s.log.Error("error while sending registration msg", slog.String("userId", uId), slog.String("error", err.Error()))
			} else {
				s.log.Info("sended mail", slog.String("userId", uId))
			}
		}(uId)
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
			if errors.Is(err, storage.ErrUserNotFound) {
				return nil, status.Error(codes.NotFound, "user not found")
			}
			return nil, status.Error(codes.Internal, "failed to check admin status")
		}
		return &authProto.IsAdminResponse{
			IsAdmin: isAdm,
		}, nil
}