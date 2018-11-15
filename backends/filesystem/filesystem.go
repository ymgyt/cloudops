package filesystem

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dustin/go-humanize"

	"go.uber.org/zap"

	"github.com/ymgyt/cloudops/core"
)

// New -
func New(ctx *core.Context) (*FileSystem, error) {
	return newFileSystem(ctx)
}

func newFileSystem(ctx *core.Context) (*FileSystem, error) {
	return &FileSystem{ctx: ctx}, nil
}

// FileSystem -
type FileSystem struct {
	ctx *core.Context
}

// Put -
func (fs *FileSystem) Put(in *core.PutInput) (*core.PutOutput, error) {
	return nil, core.NotImplementedError("fileSystem.Put()")
}

// Fetch -
func (fs *FileSystem) Fetch(in *core.FetchInput) (*core.FetchOutput, error) {
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
func (fs *FileSystem) Remove(in *core.RemoveInput) (*core.RemoveOutput, error) {
	sideEffect, err := fs.removeSideEffect(in.Resources, in.Dryrun)
	if err != nil {
		return nil, err
	}
	if err := fs.doRemove(sideEffect); err != nil {
		return nil, err
	}
	removedNum := len(in.Resources)
	if in.Dryrun {
		removedNum = 0
	}
	return &core.RemoveOutput{
		RemoveNum: removedNum,
	}, nil
}

func (fs *FileSystem) removeSideEffect(resources core.Resources, dryrun bool) ([]*removeSideEffect, error) {
	var se = make([]*removeSideEffect, 0, len(resources))
	for _, r := range resources {
		fp, err := fs.trimScheme(r.URI())
		if err != nil {
			return nil, err
		}
		se = append(se, &removeSideEffect{
			Dryrun:   dryrun,
			Filepath: fp,
		})
	}
	return se, nil
}

func (fs *FileSystem) trimScheme(path string) (string, error) {
	const scheme = "file://"
	if !strings.HasPrefix(path, scheme) {
		return "", core.NewError(core.InvalidParam, fmt.Sprintf("invalid file path %s", path))
	}
	if len(path) <= len(scheme) {
		return "", core.NewError(core.InvalidParam, fmt.Sprintf("invalid file path %s", path))
	}
	return path[len(scheme):], nil
}

func (fs *FileSystem) doRemove(sideEffects []*removeSideEffect) error {
	log := fs.ctx.Log
	for _, se := range sideEffects {
		if se.Dryrun {
			log.Info("remove", zap.String("file", se.Filepath), zap.Bool("dryrun", se.Dryrun))
			continue
		}
		log.Info("remove", zap.String("file", se.Filepath))
		if err := os.Remove(se.Filepath); err != nil {
			return core.WrapError(core.Internal, "", err)
		}
	}
	return nil
}

func (fs *FileSystem) fetchFiles(path string, exp string) (core.Resources, error) {
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

func (fs *FileSystem) fetchFile(path string) (core.Resources, error) {
	return core.Resources{fs.resource(path)}, nil
}

func (fs *FileSystem) resource(path string) core.Resource {
	return &fileResource{path: path}
}

func (fs *FileSystem) includePredicator(exp string) (func(string) bool, error) {
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

// removeSideEffect represents remove operation for remove side effect to be thin.
type removeSideEffect struct {
	Dryrun   bool
	Filepath string
}

// DiskUsageInput -
type DiskUsageInput struct {
	Root string
}

// DiskUsageOutput -
type DiskUsageOutput struct {
	Root *Dir
}

// Dir -
type Dir struct {
	Path  string
	Files []os.FileInfo
	Dirs  []*Dir
}

// Size -
func (d *Dir) Size() uint64 {
	var s uint64
	for _, info := range d.Files {
		s += uint64(info.Size())
	}
	for _, dir := range d.Dirs {
		s += uint64(dir.Size())
	}
	return s
}

// Dump -
func (d *Dir) Dump(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%s %s\n", d.Path, humanize.Bytes(d.Size())); err != nil {
		return err
	}
	for _, dir := range d.Dirs {
		if err := dir.Dump(w); err != nil {
			return err
		}
	}
	return nil
}

func makeDirTree(ctx context.Context, root string) (*Dir, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	infos, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}
	d := &Dir{Path: root}

	for _, info := range infos {
		if info.IsDir() {
			dir, err := makeDirTree(ctx, filepath.Join(root, info.Name()))
			if err != nil {
				return nil, err
			}
			d.Dirs = append(d.Dirs, dir)
		}

		d.Files = append(d.Files, info)
	}

	return d, nil
}

// DiskUsage -
func (fs *FileSystem) DiskUsage(in *DiskUsageInput) (*DiskUsageOutput, error) {
	root, err := makeDirTree(fs.ctx.Ctx, in.Root)

	return &DiskUsageOutput{
		Root: root,
	}, core.WrapError(core.Internal, "", err)
}
