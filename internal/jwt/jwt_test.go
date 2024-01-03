package jwt

import (
	"testing"
	"time"

	"github.com/EwvwGeN/authService/internal/domain/models"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FullOkay(t *testing.T) {
	user := models.User{
		Id: "user_id_1",
		Email: "test@user.first",
		PassHash: []byte("test_password_1"),
	}
	app := models.App{
		Id: "app_id_1",
		Name: "test app",
		Secret: "123123123",
	}

	duration := time.Duration(10) * time.Second

	got, err := NewToken(user, app, duration)
	creatingTime := time.Now()
	require.NoError(t, err)

	parsedToken, err := jwt.Parse(got, func(t *jwt.Token) (interface{}, error) {
		return []byte(app.Secret), nil
	})

	require.NoError(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	require.True(t, ok)


	assert.Equal(t, user.Id, claims["user_id"].(string))
	assert.Equal(t, user.Email, claims["email"].(string))
	assert.Equal(t, app.Id, claims["app_id"].(string))

	assert.InDelta(t, creatingTime.Add(duration).Unix(), claims["exp"].(float64), 1)
}

func Test_UserWithoutValue(t *testing.T) {
	user := models.User{
		Id: "user_id_2",
		PassHash: []byte("test_password_2"),
	}
	app := models.App{
		Id: "app_id_2",
		Name: "test app",
		Secret: "123123123",
	}

	_, err := NewToken(user, app, time.Second)
	require.ErrorIs(t, err, ErrEmptyValue)
}

func Test_AppWithoutValue(t *testing.T) {
	user := models.User{
		Id: "user_id_3",
		Email: "test@user.third",
		PassHash: []byte("test_password_3"),
	}
	app := models.App{
		Id: "app_id_3",
		Name: "test app",
	}

	_, err := NewToken(user, app, time.Second)
	require.ErrorIs(t, err, ErrEmptyValue)
}