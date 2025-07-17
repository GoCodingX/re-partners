package repository

import "fmt"

// NewAlreadyExistsError creates a new AlreadyExistsError instance.
func NewAlreadyExistsError(field, value string, err error) *AlreadyExistsError {
	return &AlreadyExistsError{
		Msg: fmt.Sprintf("record with %s=%s already exists", field, value),
		Err: err,
	}
}

type AlreadyExistsError struct {
	Msg string
	Err error
}

func (e *AlreadyExistsError) Error() string {
	return e.Msg + ": " + e.Err.Error()
}

func (e *AlreadyExistsError) Unwrap() error {
	return e.Err
}
