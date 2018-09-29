package s3

import (
	"testing"

	"github.com/ymgyt/cloudops/core"
	"github.com/ymgyt/cloudops/testutil"
)

func TestS3Client_split(t *testing.T) {
	tests := []struct {
		desc      string
		path      string
		bucket    string
		key       string
		wantError bool
		err       error
	}{
		{"ok", "s3://bucket/object.ext", "bucket", "object.ext", false, nil},
		{"ok prefix", "s3://bucket/prefix/object.ext", "bucket", "prefix/object.ext", false, nil},
		{"invalid no colon", "s3//bucket/object.ext", "", "", true, core.NewError(core.InvalidParam, "")},
		{"invalid no slash", "s3:bucket/object.ext", "", "", true, core.NewError(core.InvalidParam, "")},
		{"invalid no key", "s3://bucket/", "", "", true, core.NewError(core.InvalidParam, "")},
		{"invalid slash slash", "s3:///", "", "", true, core.NewError(core.InvalidParam, "")},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			client := s3Client{ctx: testutil.DammyContext()}

			bucket, key, err := client.split(tt.path)

			if !testutil.AssertError(t, tt.wantError, err, tt.err) {
				return
			}

			testutil.Diff(t, bucket, tt.bucket)
			testutil.Diff(t, key, tt.key)
		})
	}
}
