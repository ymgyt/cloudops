package main

import (
	"fmt"

	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/davecgh/go-spew/spew"
	cli "github.com/jawher/mow.cli"

	"github.com/ymgyt/cloudops/core"
)

const (
	localRemoteHdler = "localRemote"
)

// CopyCommand represents copy operation.
type CopyCommand struct {
	ctx *core.Context

	dryrun    bool
	recursive bool

	src  string
	dest string
}

// Run is a entrypoint.
func (cmd *CopyCommand) Run() {
	cmd.printStart()
	if err := cmd.run(); err != nil {
		cmd.ctx.Log.Error("copy", zap.Error(err))
		cli.Exit(1)
	}
}

func (cmd *CopyCommand) run() error {
	req, err := cmd.newRequest()
	if err != nil {
		return err
	}

	name, err := cmd.dispatch(req)
	if err != nil {
		return err
	}
	var handler copyHandler
	switch name {
	case localRemoteHdler:
		handler = cmd.localRemote
	}

	if handler == nil {
		return core.NewError(core.NotImplementedYet, fmt.Sprintf("handler not implemented yet %+v", req))
	}

	if err := handler(req); err != nil {
		return err
	}
	return nil
}

type copyRequest struct {
	Dryrun    bool
	Recursive bool
	Src       string `validate:"required"`
	Dest      string `validate:"required"`
}

func (cmd *CopyCommand) newRequest() (*copyRequest, error) {
	req := &copyRequest{
		Dryrun:    cmd.dryrun,
		Recursive: cmd.recursive,
		Src:       cmd.src,
		Dest:      cmd.dest,
	}
	return req, req.validate(cmd.ctx.Validate)
}

// MEMO: should i return copy output ?
type copyHandler func(*copyRequest) error

func (cmd *CopyCommand) dispatch(req *copyRequest) (handler string, err error) {
	src, dest, err := cmd.readTargets(req)
	spew.Dump(src, dest)
	if err != nil {
		return "", err
	}
	switch {
	case src.Type == core.LocalResource && dest.Type == core.RemoteResource:
		handler = localRemoteHdler
	default:
		err = core.NewError(core.NotImplementedYet, fmt.Sprintf("copy %v to %v", src.Type, dest.Type))
	}
	return handler, err
}

// FIXME: more good name
func (cmd *CopyCommand) localRemote(req *copyRequest) error {
	spew.Dump(req)
	return nil
}

func (req *copyRequest) validate(validate *validator.Validate) error {
	if err := validate.Struct(req); err != nil {
		return core.WrapError(core.InvalidParam, "", err)
	}
	return nil
}

func (cmd *CopyCommand) readTargets(req *copyRequest) (src *core.Resource, dest *core.Resource, err error) {
	src, err = core.NewResource(req.Src)
	if err != nil {
		return nil, nil, err
	}
	dest, err = core.NewResource(req.Dest)
	if err != nil {
		return nil, nil, err
	}
	return src, dest, nil
}

func (cmd *CopyCommand) printStart() {
	cmd.ctx.Log.Info("copy start",
		zap.String("src", cmd.src), zap.String("dest", cmd.dest),
		zap.Bool("dryrun", cmd.dryrun), zap.Bool("recursive", cmd.recursive))
}
