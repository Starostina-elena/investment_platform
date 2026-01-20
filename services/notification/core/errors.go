package core

import "errors"

var (
	ErrInvalidEmailRequest = errors.New("invalid email request")
	ErrEmailSendFailed     = errors.New("failed to send email")
	ErrUnknownNotifType    = errors.New("unknown notification type")
)
