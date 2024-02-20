package verification

import "errors"

var (
	ErrInternal = errors.New("internal error")
	ErrKeySession = errors.New("key session not started")
	ErrCodeSize = errors.New("code too short")
)