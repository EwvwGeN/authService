package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/EwvwGeN/authService/internal/domain/models"
	"github.com/EwvwGeN/authService/internal/jwt"
	"github.com/EwvwGeN/authService/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log *slog.Logger
	usrSaver UserSaver
	usrProvider UserProvider
	appProdiver AppProvider
	tokenTTL time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context,
	email string,
	passHash []byte,
	) (uid string, err error)
}

type UserProvider interface {
	GetUser(ctx context.Context,
		email string,
	) (models.User, error)
	IsAdmin(ctx context.Context,
		userId string,
	) (bool, error)
}

type AppProvider interface {
	GetApp(ctx context.Context,
		appId string,
	) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func NewAuthService(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
	) *Auth {
	return &Auth{
		log: log,
		usrSaver: userSaver,
		usrProvider: userProvider,
		appProdiver: appProvider,
		tokenTTL: tokenTTL,
	}
}

func (a *Auth) Login( ctx context.Context,
	email string,
	password string,
	appId string,
) (string, error) {
	a.log.Info("attempting to login user")

	user, err := a.usrProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", slog.String("error", err.Error()))
			return "", fmt.Errorf("can't login user: %w", ErrInvalidCredentials)
		}
		a.log.Error("failed to get user", slog.String("error", err.Error()))
		return "", fmt.Errorf("can't login user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Warn("invalid credentials", slog.String("error", err.Error()))
		return "", fmt.Errorf("can't login user: %w", ErrInvalidCredentials)
	}

	app, err := a.appProdiver.GetApp(ctx, appId)
	if err != nil {
		return "", fmt.Errorf("can't login user: %w", err)
	}

	a.log.Info("user logged successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", slog.String("error", err.Error()))
		return "", fmt.Errorf("can't login user: %w", err)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser( ctx context.Context,
	email string,
	password string,
) (string, error) {
	a.log.Info("registering user", slog.String("email", email))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error("failed to generate password hash", slog.String("error", err.Error()))
		return "", fmt.Errorf("can't register user: %w", err)
	}
	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExist) {
			a.log.Warn("failed to save user", slog.String("error", err.Error()))
			return "", fmt.Errorf("can't register user: %w", storage.ErrUserExist)
		}
		a.log.Error("failed to save user", slog.String("error", err.Error()))
		return "", fmt.Errorf("can't register user: %w", err)
	}
	a.log.Info("user registered", slog.String("UserId", id))
	return id, nil
}

func (a *Auth) IsAdmin( ctx context.Context,
	uId string,
) (bool, error) {
	a.log.Info("checking admin")

	isAdm, err := a.usrProvider.IsAdmin(ctx, uId)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", slog.String("error", err.Error()))
			return false, fmt.Errorf("can't check user: %w", storage.ErrUserNotFound)
		}
		a.log.Error("failed to get user", slog.String("error", err.Error()))
		return false, fmt.Errorf("can't check user: %w", err)
	}

	a.log.Info("checked user if user is admin", slog.Bool("isAdmin", isAdm))

	return isAdm, nil
}