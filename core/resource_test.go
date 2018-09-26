package core_test

import (
	"testing"

	"github.com/ymgyt/cloudops/core"
	"github.com/ymgyt/cloudops/testutil"
)

func TestNewResource(t *testing.T) {
	tests := []struct {
		desc      string
		path      string
		want      core.Resource
		wantError bool
		err       error
	}{
		{
			desc: "relative path",
			path: "./path/to/src",
			want: core.Resource{
				Type: core.LocalResource,
			},
		},
		{
			desc: "absolute path",
			path: "/path/to/src",
			want: core.Resource{
				Type: core.LocalResource,
			},
		},
		{
			desc: "empty is invalid",
			path: "",
			want: core.Resource{
				Type: core.InvalidResource,
			},
		},
		{
			desc: "s3 is remote",
			path: "s3://",
			want: core.Resource{
				Type: core.RemoteResource,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := core.NewResource(tt.path)

			if !testutil.AssertError(t, tt.wantError, err, tt.err) {
				return
			}

			testutil.Diff(t, *got, tt.want)
		})
	}
}
