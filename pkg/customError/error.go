package customError

import "fmt"

type Cause string

const (
	InvalidInput  Cause = "INVALID_INPUT"
	InternalError Cause = "INTERNAL_ERROR"
	NotFound      Cause = "NOT_FOUND"
	AlreadyExists Cause = "ALREADY_EXISTS"
)

type CustomError struct {
	Cause Cause
	Msg   string
}

func New(cause Cause, msg string) *CustomError {
	return &CustomError{
		Cause: cause,
		Msg:   msg,
	}
}

func (ce *CustomError) Error() string {
	return fmt.Sprintf("cause: %s, msg: %s", ce.Cause, ce.Msg)
}
