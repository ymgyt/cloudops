package usecase

import (
	"fmt"

	"github.com/ymgyt/cloudops/core"
)

//go:generate mockgen -destination=../testutil/file_operation_mock.go -package=testutil github.com/ymgyt/cloudops/usecase FileOps

// FileOps -
type FileOps interface {
	CopyLocalToRemote(*CopyLocalToRemoteInput) (*CopyLocalToRemoteOutput, error)
}

// CopyLocalToRemoteInput -
type CopyLocalToRemoteInput struct {
	Dryrun      bool
	Recursive   bool
	CreateDir   bool
	SkipConfirm bool
	Regexp      string
	Remove      bool
	Src         string `validate:"required"`
	Dest        string `validate:"required"`
}

// CopyLocalToRemoteOutput -
type CopyLocalToRemoteOutput struct {
	CopiedNum  int
	RemovedNum int
}

// RemoveLocalInput -
type RemoveLocalInput struct{}

// RemoveLocalOutput -
type RemoveLocalOutput struct{}

// NewFileOps -
func NewFileOps(ctx *core.Context, fileSystem core.Backend, s3 core.Backend, confirmer core.Confirmer) (FileOps, error) {
	return &fileOps{
		ctx:        ctx,
		fileSystem: fileSystem,
		s3:         s3,
		confirmer:  confirmer,
	}, nil
}

type fileOps struct {
	ctx        *core.Context
	fileSystem core.Backend
	s3         core.Backend
	confirmer  core.Confirmer
}

// CopyLocalToRemote -
func (fo *fileOps) CopyLocalToRemote(in *CopyLocalToRemoteInput) (*CopyLocalToRemoteOutput, error) {
	if err := core.Validate(in, fo.ctx.Validate); err != nil {
		return nil, err
	}

	out, err := fo.copyLocalToRemote(in)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (fo *fileOps) copyLocalToRemote(in *CopyLocalToRemoteInput) (*CopyLocalToRemoteOutput, error) {
	fetchOut, err := fo.fileSystem.Fetch(&core.FetchInput{
		Recursive: in.Recursive,
		Regexp:    in.Regexp,
		Src:       in.Src})
	if err != nil {
		return nil, err
	}

	resources := fetchOut.Resources
	if !in.SkipConfirm {
		var msg string = fmt.Sprintf("copy files to %s", in.Dest)
		if in.Dryrun {
			msg = "[Dryrun] " + msg
		}
		ok, err := fo.confirmer.Confirm(msg, resources)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, core.NewError(core.Canceled, "")
		}
	}

	putOut, err := fo.s3.Put(&core.PutInput{
		Dryrun:    in.Dryrun,
		Recursive: in.Recursive,
		CreateDir: in.CreateDir,
		Dest:      in.Dest,
		Resources: resources})
	if err != nil {
		return nil, err
	}

	var removedNum int
	if in.Remove {
		if !in.SkipConfirm {
			var msg string = fmt.Sprintf("delete above file(s)")
			if in.Dryrun {
				msg = "[Dryrun] " + msg
			}
			ok, err := fo.confirmer.Confirm(msg, resources)
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, core.NewError(core.Canceled, "")
			}
		}
		removeOut, err := fo.fileSystem.Remove(&core.RemoveInput{
			Dryrun:    in.Dryrun,
			Resources: resources,
		})
		if err != nil {
			return nil, err
		}
		removedNum = removeOut.RemoveNum
	}

	return &CopyLocalToRemoteOutput{
		CopiedNum:  putOut.PutNum,
		RemovedNum: removedNum,
	}, nil
}
