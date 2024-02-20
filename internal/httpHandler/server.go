package httpHandler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/EwvwGeN/authService/internal/config"
	"github.com/EwvwGeN/authService/internal/domain/models"
	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	cfg config.HttpConfig
	log *slog.Logger
	confirmator UserСonfirmator
	usrProvider UserProvider
	closer func(ctx context.Context) error
}

type UserСonfirmator interface {
	ConfirmUser(ctx context.Context,
		userId string,
	) error
}

type UserProvider interface {
	GetUserById(ctx context.Context,
		userId string,
	) (models.User, error)
}

func NewServer(cfg config.HttpConfig, log *slog.Logger, confirmator UserСonfirmator, usrProvider UserProvider) *Server {
	return &Server{
		router: mux.NewRouter(),
		cfg: cfg,
		log: log,
		confirmator: confirmator,
		usrProvider: usrProvider,
	}
}

func (s *Server) Start() error {
	s.configureRouter()
	srv := &http.Server{
		Handler: s.router,
		Addr:    fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	s.closer = func(ctx context.Context) error {
		return srv.Shutdown(ctx)
	}
	s.log.Info("starting listening", slog.String("addres", srv.Addr))
	return srv.ListenAndServe()
}

func (s *Server) configureRouter() {
	s.log.Info("configuring router")
	s.router.HandleFunc("/user/verification/{verificationCode}", s.confirmUser())
}

func (s *Server) Close(ctx context.Context) error {
	return s.closer(ctx)
}