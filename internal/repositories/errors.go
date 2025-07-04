package repositories

import "errors"

var (
	ErrEntityNotFound  = errors.New("entity not found")
	ErrSessionNotFound = errors.New("session not found")
)
