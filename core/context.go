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

	cancel func()
}

// NewContext -
func NewContext(background context.Context, log *zap.Logger, v *validator.Validate) *Context {
	ctx, cancel := context.WithCancel(background)
	return &Context{Log: log, Validate: v, Ctx: ctx, cancel: cancel}
}

// Cancel -
func (c *Context) Cancel() {
	c.cancel()
}
