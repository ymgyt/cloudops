package filesystem

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/ymgyt/cloudops/core"
)

// NewFileSystem -
func New(ctx *core.Context) (core.Backend, error) {
	return newFileSystem(ctx)
}

func newFileSystem(ctx *core.Context) (*fileSystem, error) {
	return &fileSystem{ctx: ctx}, nil
}

type fileSystem struct {
	ctx *core.Context
}

// Put -
func (fs *fileSystem) Put(in *core.PutInput) (*core.PutOutput, error) {
	return nil, core.NotImplementedError("fileSystem.Put()")
}

// Fetch -
func (fs *fileSystem) Fetch(in *core.FetchInput) (*core.FetchOutput, error) {
	src := in.Src
	stat, err := os.Stat(src)
	if os.IsNotExist(err) {
		return nil, core.WrapError(core.NotFound, "", err)
	}
	var resoureces core.Resources
	if stat.IsDir() {
		if !in.Recursive {
			return nil, core.NewError(core.InvalidParam, fmt.Sprintf("%s is directory", src))
		}
		resoureces, err = fs.fetchFiles(src, in.Regexp)
	} else {
		resoureces, err = fs.fetchFile(src)
	}
	if len(resoureces) == 0 {
		return nil, core.NewError(core.NotFound, fmt.Sprintf("files not found in %s", src))
	}
	return &core.FetchOutput{
		Resources: resoureces,
	}, err
}

// Remove -
func (fs *fileSystem) Remove(in *core.RemoveInput) (*core.RemoveOutput, error) {
	return nil, core.NotImplementedError("fileSystem.Remove()")
}

func (fs *fileSystem) fetchFiles(path string, exp string) (core.Resources, error) {
	includer, err := fs.includePredicator(exp)
	if err != nil {
		return nil, err
	}
	var rs core.Resources
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if includer(path) {
			rs = append(rs, fs.resource(path))
		}
		return nil
	})
	return rs, err
}

func (fs *fileSystem) fetchFile(path string) (core.Resources, error) {
	return core.Resources{fs.resource(path)}, nil
}

func (fs *fileSystem) resource(path string) core.Resource {
	return &fileResource{path: path}
}

func (fs *fileSystem) includePredicator(exp string) (func(string) bool, error) {
	if exp == "" {
		return includeAll, nil
	}
	r, err := regexp.Compile(exp)
	if err != nil {
		return nil, core.WrapError(core.InvalidParam, "", err)
	}
	return func(path string) bool {
		return r.MatchString(path)
	}, nil
}

func includeAll(_ string) bool {
	return true
}

type fileResource struct {
	path string
}

// Type -
func (r *fileResource) Type() core.ResourceType {
	return core.LocalFileResource
}

// URI -
func (r *fileResource) URI() string {
	const scheme = "file://"
	return scheme + r.path
}

// Open -
func (r *fileResource) Open() (io.ReadCloser, error) {
	f, err := os.Open(r.path)
	if os.IsNotExist(err) {
		return nil, core.WrapError(core.NotFound, fmt.Sprintf("%s not found", r.path), err)
	} else if os.IsPermission(err) {
		return nil, core.WrapError(core.Unauthenticated, fmt.Sprintf("open %s unauthorized", r.path), err)
	} else if err != nil {
		return nil, core.WrapError(core.Internal, "", err)
	}
	return f, nil
}
