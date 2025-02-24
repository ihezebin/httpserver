package httpserver

import (
	"net/http"

	"github.com/pkg/errors"
)

type Code int

const (
	CodeOK Code = iota

	CodeValidateRuleFailed
	CodeInternalServerError
	CodeBadRequest
	CodeUnauthorized
	CodeNotFound
	CodeForbidden
	CodeTimeout

	CodeCreated
	CodeAccepted
	CodeNoContent
	CodeResetContent
	CodeAuthorizationFailed
)

var code2MessageM = map[Code]string{
	CodeOK:                  "OK",
	CodeInternalServerError: "Internal Server Error",
	CodeBadRequest:          "Bad Request",
	CodeUnauthorized:        "Unauthorized",
	CodeNotFound:            "Not Found",
	CodeForbidden:           "Forbidden",
	CodeTimeout:             "Timeout",
	CodeCreated:             "Created",
	CodeAccepted:            "Accepted",
	CodeNoContent:           "No Content",
	CodeResetContent:        "Reset Content",
	CodeValidateRuleFailed:  "Validate Rule Failed",
}

type Err struct {
	Status int
	Code   Code
	Err    error
}

var _ error = &Err{}

func (e *Err) Error() string {
	return e.Err.Error()
}

func (e *Err) WithStatus(status int) *Err {
	e.Status = status
	return e
}

func ErrorWithCode(code Code) *Err {
	return &Err{
		Code: code,
		Err:  errors.New(code2MessageM[code]),
	}
}

func NewError(code Code, msg string) *Err {
	return &Err{
		Code: code,
		Err:  errors.New(msg),
	}
}

func ErrorWithBadRequest() *Err {
	return &Err{
		Status: http.StatusBadRequest,
		Code:   CodeBadRequest,
		Err:    errors.New(code2MessageM[CodeBadRequest]),
	}
}
func ErrorWithInternalServer() *Err {
	return &Err{
		Code: CodeInternalServerError,
		Err:  errors.New(code2MessageM[CodeInternalServerError]),
	}
}

func ErrWithUnAuthorized() *Err {
	return &Err{
		Status: http.StatusUnauthorized,
		Code:   CodeUnauthorized,
		Err:    errors.New(code2MessageM[CodeUnauthorized]),
	}
}

func ErrorWithAuthorizationFailed(reason string) *Err {
	return &Err{
		Status: http.StatusUnauthorized,
		Code:   CodeAuthorizationFailed,
		Err:    errors.New(reason),
	}
}
