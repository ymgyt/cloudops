package testutil

import (
	"context"
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
