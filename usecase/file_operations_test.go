package usecase_test

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/ymgyt/cloudops/core"
	"github.com/ymgyt/cloudops/testutil"
	"github.com/ymgyt/cloudops/usecase"
)

func TestFileOps_CopyLocalToRemote_S3(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctr := gomock.NewController(t)
		defer ctr.Finish()

		fs := testutil.NewMockBackend(ctr)
		s3 := testutil.NewMockBackend(ctr)
		confirmer := testutil.NewMockConfirmer(ctr)

		fs.EXPECT().Fetch(gomock.Any()).Return(&core.FetchOutput{}, nil).Times(1)
		s3.EXPECT().Put(gomock.Any()).Return(&core.PutOutput{}, nil).Times(1)
		confirmer.EXPECT().Confirm(gomock.Any(), gomock.Any()).Return(true, nil).Times(1)

		fileOps, err := usecase.NewFileOps(testutil.DammyContext(), fs, s3, confirmer)
		if fileOps == nil || err != nil {
			t.Fatalf("failed to NewFileOps %s", err)
		}

		got, err := fileOps.CopyLocalToRemote(&usecase.CopyLocalToRemoteInput{Src: "path/to/src.txt", Dest: "s3://bucket/object"})
		if err != nil {
			t.Fatalf("failed to FileOps.CopyLocalToRemote %s", err)
		}

		if got == nil {
			t.Fatal("FileOps.CopyLocalToRemote return nil")
		}

		want := usecase.CopyLocalToRemoteOutput{}
		testutil.Diff(t, *got, want)
	})
}
