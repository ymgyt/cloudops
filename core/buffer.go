package core

import (
	"bytes"
	"sync"
)

var pool = &bufferPool{p: &sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}}

type bufferPool struct {
	p *sync.Pool
}

// Buffer -
func Buffer() *bytes.Buffer {
	b := pool.p.Get().(*bytes.Buffer)
	return b
}

// PutBuffer -
func PutBuffer(b *bytes.Buffer) {
	b.Reset()
	pool.p.Put(b)
}
