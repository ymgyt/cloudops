package core

import (
	"fmt"

	validator "gopkg.in/go-playground/validator.v9"
)

// ErrCode -
type ErrCode int

//go:generate stringer -type ErrCode errors.go
const (
	OK ErrCode = iota
	InvalidParam
	Conflict
	Timeout
	Internal
	External
	NotFound
	Unauthorized
	Unauthenticated
	RateLimit
	NotImplementedYet
	Canceled
	Undefined
)

// Error -
type Error interface {
	Code() ErrCode
	error
}

// NewError -
func NewError(code ErrCode, msg string) Error {
	return newError(code, msg)
}

// WrapError -
func WrapError(code ErrCode, msg string, err error) Error {
	if err == nil {
		return nil
	}
	e := newError(code, msg)
	e.err = err
	return e
}

// NotImplementedError -
func NotImplementedError(msg string) Error {
	return newError(NotImplementedYet, msg)
}

// Validate -
func Validate(input interface{}, v *validator.Validate) Error {
	return WrapError(InvalidParam, fmt.Sprintf("%v", input), v.Struct(input))
}

// ErrorCode -
func ErrorCode(err error) ErrCode {
	if err == nil {
		return OK
	}
	if coreErr, ok := err.(Error); ok {
		return coreErr.Code()
	}
	return Undefined
}

func newError(code ErrCode, msg string) *coreError {
	return &coreError{
		code: code,
		msg:  msg,
	}
}

type coreError struct {
	code ErrCode
	msg  string
	err  error
}

func (e *coreError) Code() ErrCode {
	return e.code
}

func (e *coreError) Error() string {
	if e.err == nil {
		return fmt.Sprintf("%s %s", e.code, e.msg)
	}
	return fmt.Sprintf("%s %s:%s", e.code, e.msg, e.err.Error())
}
