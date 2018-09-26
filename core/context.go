package core

import (
	"context"

	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"
)

// Context -
type Context struct {
	Log      *zap.Logger
	Validate *validator.Validate
	Ctx      context.Context
}

// NewContext -
func NewContext(log *zap.Logger, v *validator.Validate) *Context {
	return &Context{Log: log, Validate: v}
}
