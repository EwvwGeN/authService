package storage

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExist    = errors.New("user already exist")
	ErrAppNotFound  = errors.New("app not found")
)