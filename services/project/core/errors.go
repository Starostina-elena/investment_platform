package core

import "errors"

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrNotAuthorized   = errors.New("not authorized")
	ErrInvalidInput    = errors.New("invalid input")
)
