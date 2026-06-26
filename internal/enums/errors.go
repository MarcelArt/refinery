package enums

import "errors"

var (
	ErrAlreadyRegsitered   = errors.New("user already registered")
	ErrUnknownWorkflowType = errors.New("unknown workflow type")
)
