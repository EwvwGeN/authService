package jwt

import (
	"errors"
	"time"

	"github.com/EwvwGeN/authService/internal/domain/models"
	"github.com/golang-jwt/jwt"
)

var (
	ErrEmptyValue = errors.New("empty value")
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	if user.Id == "" {
		return "", ErrEmptyValue
	}
	if user.Email == "" {
		return "", ErrEmptyValue
	}
	if app.Id == "" {
		return "", ErrEmptyValue
	}
	if app.Secret == "" {
		return "", ErrEmptyValue
	}
	claims["user_id"] = user.Id
	claims["email"] = user.Email
	claims["app_id"] = app.Id
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}