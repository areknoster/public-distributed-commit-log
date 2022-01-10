package storage

import "errors"

var (
	ErrNotFound = errors.New("message was not found")
	ErrTimeout  = errors.New("timeout when accessing message")
	ErrInternal = errors.New("internal error")
)
