package filesystem

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/ymgyt/cloudops/core"
)

func TestFileSystem_trimScheme(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"file:///path/to/file.txt", "/path/to/file.txt"},
		{"file://relative/file.txt", "relative/file.txt"},
		{"file://./dot/file.txt", "./dot/file.txt"},
	}

	for _, tt := range tests {
		fs := &fileSystem{}
		got, err := fs.trimScheme(tt.path)

		if err != nil {
			t.Fatalf("trim %s,want no error, got %s", tt.path, err)
		}
		if diff := cmp.Diff(got, tt.want); diff != "" {
			t.Errorf("(-got +want)\n%s", diff)
		}
	}
}

func TestFileSystem_removeSideEffect(t *testing.T) {
	tests := []struct {
		desc      string
		resources core.Resources
		dryrun    bool
		want      []removeSideEffect
	}{
		{
			desc: "ok dryrun false",
			resources: core.Resources{
				&fileResource{path: "dammy1.txt"},
				&fileResource{path: "/root/path/dammy2.txt"},
				&fileResource{path: "./relative/dammy3.txt"},
			},
			dryrun: false,
			want: []removeSideEffect{
				{Dryrun: false, Filepath: "dammy1.txt"},
				{Dryrun: false, Filepath: "/root/path/dammy2.txt"},
				{Dryrun: false, Filepath: "./relative/dammy3.txt"},
			},
		},
		{
			desc: "ok dryrun true",
			resources: core.Resources{
				&fileResource{path: "dammy1.txt"},
				&fileResource{path: "/root/path/dammy2.txt"},
				&fileResource{path: "./relative/dammy3.txt"},
			},
			dryrun: true,
			want: []removeSideEffect{
				{Dryrun: true, Filepath: "dammy1.txt"},
				{Dryrun: true, Filepath: "/root/path/dammy2.txt"},
				{Dryrun: true, Filepath: "./relative/dammy3.txt"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			fs := &fileSystem{}
			got, err := fs.removeSideEffect(tt.resources, tt.dryrun)

			if err != nil {
				t.Fatalf("want no error, got %s", err)
			}

			for i := range got {
				if diff := cmp.Diff(got[i], &tt.want[i]); diff != "" {
					t.Errorf("(-got +want)\n%s", diff)
				}
			}
		})
	}
}
