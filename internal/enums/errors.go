package enums

import "errors"

var (
	ErrAlreadyRegsitered   = errors.New("user already registered")
	ErrUnknownWorkflowType = errors.New("unknown workflow type")
	ErrDummyEmailOnProd    = errors.New("forbidden to use dummy email on prod")
)
