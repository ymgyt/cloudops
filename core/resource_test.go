package core_test

import (
	"testing"

	"github.com/ymgyt/cloudops/core"
	"github.com/ymgyt/cloudops/testutil"
)

func TestInspectPath(t *testing.T) {
	tests := []struct {
		desc string
		path string
		want core.ResourceType
	}{
		{"dot relative path", "./path/to/src/txt", core.LocalFileResource},
		{"relative path", "path/to/src.txt", core.LocalFileResource},
		{"absolute path", "/path/to/src.txt", core.LocalFileResource},
		{"s3", "s3://bucket/prefix/object", core.S3Resource},
		{"empty", "", core.InvalidResource},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, want := core.InspectPath(tt.path), tt.want
			testutil.Diff(t, got, want)
		})
	}
}

func TestResourceType_IsLocal(t *testing.T) {
	tests := []struct {
		desc string
		rsc  core.ResourceType
		want bool
	}{
		{"local file", core.LocalFileResource, true},
		{"invalid never local", core.InvalidResource, false},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, want := tt.rsc.IsLocal(), tt.want
			if got != want {
				t.Errorf("%v IsLocal does not match. want %v, got %v", tt.rsc, want, got)
			}
		})
	}
}
