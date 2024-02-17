package queue

import "github.com/EwvwGeN/authService/internal/domain/models"

type messageWithErrCh struct {
	msg *models.Message
	errCh chan error
}