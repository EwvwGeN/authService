package storage

import "errors"

var (
	ErrDbNotExist = errors.New("database does not exist")
	ErrCollNotExist = errors.New("some of collections does not exist")

	ErrUserNotFound = errors.New("user not found")
	ErrUserExist    = errors.New("user already exist")
	ErrUserConfirm    = errors.New("user already confirmed")
	ErrAppNotFound  = errors.New("app not found")
	ErrValidation = errors.New("validation error")
)