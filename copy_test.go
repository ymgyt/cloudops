package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/ymgyt/cloudops/core"
	"github.com/ymgyt/cloudops/testutil"
)

func TestCopyCommand_newRequest(t *testing.T) {
	tests := []struct {
		desc      string
		copyCmd   CopyCommand
		want      copyRequest
		wantError bool
		err       error
	}{
		{
			desc: "ok",
			copyCmd: CopyCommand{
				src:  "src.txt",
				dest: "dest.txt",
			},
			want: copyRequest{
				Src:       "src.txt",
				Dest:      "dest.txt",
				Dryrun:    false,
				Recursive: false,
			},
		},
		{
			desc: "src is required",
			copyCmd: CopyCommand{
				src:  "",
				dest: "dest.txt",
			},
			wantError: true,
			err:       core.NewError(core.InvalidParam, ""),
		},
		{
			desc: "dest is required",
			copyCmd: CopyCommand{
				src:  "src.txt",
				dest: "",
			},
			wantError: true,
			err:       core.NewError(core.InvalidParam, ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.copyCmd.ctx = testutil.DammyContext()
			got, err := tt.copyCmd.newRequest()

			if !testutil.AssertError(t, tt.wantError, err, tt.err) {
				return
			}

			if diff := cmp.Diff(*got, tt.want, cmp.AllowUnexported(*got)); diff != "" {
				errorf(t, diff)
			}
		})
	}
}

func TestCopyCommand_dispatch(t *testing.T) {

	copy := &CopyCommand{ctx: testutil.DammyContext()}

	tests := []struct {
		desc      string
		req       copyRequest
		want      string
		wantError bool
		err       error
	}{
		{
			desc: "localToRemote",
			req: copyRequest{
				Src:  "./src.txt",
				Dest: "s3://bucket/object",
			},
			want: "localRemote",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := copy.dispatch(&tt.req)

			if !testutil.AssertError(t, tt.wantError, err, tt.err) {
				return
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("dispatched handler does not match. got %q, want %q", got, tt.want)
			}
		})
	}
}
func errorf(t *testing.T, diff string) {
	t.Errorf("(-got +want)\n%s", diff)
}
