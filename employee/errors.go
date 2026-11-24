package employee

import "errors"

var (
	ErrEmployeeNotFound = errors.New("employee not found")
	ErrCmdRepository    = errors.New("unable to command repository")
	ErrQueryRepository  = errors.New("unable to query repository")
)
