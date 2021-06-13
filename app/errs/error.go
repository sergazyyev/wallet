package errs

import (
	"net/http"
)

//ErrorCode type wrapper int
type ErrorCode int

//CustomError custom struct implementing error
type CustomError struct {
	code ErrorCode
}

//Const error codes
const (
	InternalServerError ErrorCode = 100
	NotEnoughBalance    ErrorCode = 101
	EmptyReport         ErrorCode = 102
)

var statusCodes = map[ErrorCode]int{
	InternalServerError: http.StatusInternalServerError,
	NotEnoughBalance:    http.StatusBadRequest,
	EmptyReport:         http.StatusNotFound,
}

var messages = map[ErrorCode]string{
	InternalServerError: "Internal server error, try later",
	NotEnoughBalance:    "Insufficient funds for debit",
	EmptyReport:         "Report is empty by given parameters",
}

//Error implement error interface
func (e *CustomError) Error() string {
	return messages[e.code]
}

//Status return http status for error
func (e *CustomError) Status() int {
	return statusCodes[e.code]
}

//New creates new custom error
func New(code ErrorCode) error {
	return &CustomError{code: code}
}

//IsCustomErr checks is custom error
func IsCustomErr(err error) bool {
	switch err.(type) {
	default:
		return false
	case *CustomError:
		return true
	}
}
