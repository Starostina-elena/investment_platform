package core

import "errors"

var (
	ErrProjectNotFound     = errors.New("project not found")
	ErrNotAuthorized       = errors.New("not authorized")
	ErrInvalidInput        = errors.New("invalid input")
	ErrPaybackStarted      = errors.New("cannot change is_completed when payback has started")
	ErrPaybackNotSupported = errors.New("payback is not supported for charity and custom monetization types")
	ErrNotEnoughFunds      = errors.New("not enough funds to complete payback")
)
