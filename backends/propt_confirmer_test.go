package backends_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ymgyt/cloudops/backends"
	"github.com/ymgyt/cloudops/core"
	"github.com/ymgyt/cloudops/testutil"
)

func TestPromptConfirmer(t *testing.T) {
	var gotMsg bytes.Buffer

	p, err := backends.NewPromptConfirmer(&gotMsg, strings.NewReader("yes\n"), "[yes/no]", []string{"yes"})
	if err != nil {
		t.Fatalf("NewPromptConfirmer() failed %s", err)
	}

	rs := core.Resources{
		&testutil.Resource{FakeType: core.LocalFileResource, FakeURI: "file:///aaa.txt"},
		&testutil.Resource{FakeType: core.LocalFileResource, FakeURI: "file:///bbb.txt"},
	}

	ok, err := p.Confirm("delete", rs)
	if err != nil {
		t.Fatalf("want no error, but got %s", err)
	}

	if !ok {
		t.Error("PromptConfirmer should return true but return false")
	}

	wantMsg := `file:///aaa.txt
file:///bbb.txt
delete [yes/no] `
	testutil.Diff(t, gotMsg.String(), wantMsg)
}
