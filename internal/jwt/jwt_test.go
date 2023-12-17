package jwt

import (
	"testing"
	"time"

	"github.com/EwvwGeN/authService/internal/domain/models"
)

func TestNewToken(t *testing.T) {
	type args struct {
		user     models.User
		app      models.App
		duration time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewToken(tt.args.user, tt.args.app, tt.args.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
