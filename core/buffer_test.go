package core_test

import (
	"testing"

	"github.com/ymgyt/cloudops/core"
)

func TestBufferPool(t *testing.T) {
	b := core.Buffer()
	b.WriteString("gopher")
	core.PutBuffer(b)
	b = core.Buffer()
	if s := b.String(); s != "" {
		t.Errorf("new buffer return %s, want clean buffer", s)
	}
}
