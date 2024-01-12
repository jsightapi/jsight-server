package sync

import (
	"bytes"
	"sync"
)

// BufferPool a wrapper under sync.Pool which holds buffers.
type BufferPool struct {
	pool sync.Pool
}

// NewBufferPool creates new instance of BufferPool.
func NewBufferPool(size int) *BufferPool {
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, size))
			},
		},
	}
}

// Get returns new buffer from pool.
func (p *BufferPool) Get() *bytes.Buffer {
	return p.pool.Get().(*bytes.Buffer)
}

// Put returns buffer to pool.
func (p *BufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	p.pool.Put(b)
}
