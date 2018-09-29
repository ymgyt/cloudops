package testutil

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/ymgyt/cloudops/core"
)

// AssertError -
func AssertError(t *testing.T, wantError bool, got, want error) bool {
	t.Helper()

	if wantError && got == nil {
		t.Errorf("want %v, got nil", want)
		return false
	}
	if !wantError && got != nil {
		t.Errorf("want no error, got %s", got)
		return false
	}
	if core.ErrorCode(got) != core.ErrorCode(want) {
		t.Errorf("error code does not match. got %s, want %s", got, want)
		return false
	}
	return !wantError
}

// Diff -
func Diff(t *testing.T, got, want interface{}) {
	t.Helper()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("(-got +want)\n%s", diff)
	}
}

// DiffUnex diff unexported struct.
func DiffUnex(t *testing.T, got, want interface{}) {
	t.Helper()
	if diff := cmp.Diff(got, want, cmp.AllowUnexported(got)); diff != "" {
		t.Errorf("(-got +want)\n%s", diff)
	}
}

// DammyContext -
func DammyContext() *core.Context {
	return &core.Context{
		Log:      zap.NewNop(),
		Validate: validator.New(),
		Ctx:      context.Background(),
	}
}

// Resource
type Resource struct {
	FakeType    core.ResourceType
	FakeURI     string
	FakeContent string
}

// Type -
func (r *Resource) Type() core.ResourceType {
	return r.FakeType
}

// URI -
func (r *Resource) URI() string {
	return r.FakeURI
}

// Open -
func (r *Resource) Open() (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewReader([]byte(r.FakeContent))), nil
}
