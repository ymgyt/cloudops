package main

import (
	"fmt"

	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"

	cli "github.com/jawher/mow.cli"

	"github.com/ymgyt/cloudops/core"
	"github.com/ymgyt/cloudops/usecase"
)

const (
	localRemoteHdler = "localRemote"
)

// CopyCommand represents copy operation.
type CopyCommand struct {
	ctx *core.Context

	fileOps usecase.FileOps

	dryrun     bool
	recursive  bool
	createDir  bool
	skipPrompt bool
	remove     bool
	regexp     string

	src  string
	dest string
}

// Run is a entrypoint.
func (cmd *CopyCommand) Run() {
	cmd.printStart()
	if err := cmd.run(); err != nil {
		cmd.ctx.Log.Error("copy", zap.Error(err))
		cli.Exit(2)
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
	Dryrun     bool
	Recursive  bool
	CreateDir  bool
	SkipPrompt bool
	Remove     bool
	Regexp     string
	Src        string `validate:"required"`
	Dest       string `validate:"required"`
}

func (cmd *CopyCommand) newRequest() (*copyRequest, error) {
	req := &copyRequest{
		Dryrun:     cmd.dryrun,
		Recursive:  cmd.recursive,
		CreateDir:  cmd.createDir,
		SkipPrompt: cmd.skipPrompt,
		Remove:     cmd.remove,
		Regexp:     cmd.regexp,
		Src:        cmd.src,
		Dest:       cmd.dest,
	}
	return req, req.validate(cmd.ctx.Validate)
}

// MEMO: should i return copy output ?
type copyHandler func(*copyRequest) error

func (cmd *CopyCommand) dispatch(req *copyRequest) (handler string, err error) {
	src, dest := core.InspectPath(req.Src), core.InspectPath(req.Dest)
	if src == core.InvalidResource || dest == core.InvalidResource {
		return "", core.NewError(core.InvalidParam, fmt.Sprintf("invalid arguments src:%s, dest:%s", src, dest))
	}
	switch {
	case src.IsLocal() && dest.IsRemote():
		handler = localRemoteHdler
	default:
		err = core.NewError(core.NotImplementedYet, fmt.Sprintf("copy %v to %v", src, dest))
	}
	return handler, err
}

// FIXME: more good name
func (cmd *CopyCommand) localRemote(req *copyRequest) error {
	out, err := cmd.fileOps.CopyLocalToRemote(&usecase.CopyLocalToRemoteInput{
		Dryrun:      req.Dryrun,
		Recursive:   req.Recursive,
		CreateDir:   req.CreateDir,
		Regexp:      req.Regexp,
		SkipConfirm: req.SkipPrompt,
		Remove:      req.Remove,
		Src:         req.Src,
		Dest:        req.Dest})
	if err != nil {
		return err
	}

	return cmd.localRemoteReport(out)
}

func (req *copyRequest) validate(validate *validator.Validate) error {
	if err := validate.Struct(req); err != nil {
		return core.WrapError(core.InvalidParam, "", err)
	}
	return nil
}

func (cmd *CopyCommand) localRemoteReport(out *usecase.CopyLocalToRemoteOutput) error {
	cmd.ctx.Log.Info("copy", zap.Int("copied_files", out.CopiedNum))
	return nil
}

func (cmd *CopyCommand) printStart() {
	cmd.ctx.Log.Info("copy",
		zap.String("src", cmd.src), zap.String("dest", cmd.dest),
		zap.Bool("dryrun", cmd.dryrun), zap.Bool("recursive", cmd.recursive))
}
