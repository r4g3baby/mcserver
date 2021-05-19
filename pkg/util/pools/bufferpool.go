package pools

import (
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
	"sync"
)

type (
	BufferPool interface {
		Get(buf []byte) *bytes.Buffer
		Put(buffer *bytes.Buffer)
	}

	bufferPool struct {
		buffers sync.Pool
	}
)

var Buffer = NewBufferPool()

func (pool *bufferPool) Get(buf []byte) *bytes.Buffer {
	if buffer := pool.buffers.Get(); buffer != nil {
		buffer := buffer.(*bytes.Buffer)
		if buf != nil {
			_, _ = buffer.Write(buf)
		}
		return buffer
	}
	return bytes.NewBuffer(buf)
}

func (pool *bufferPool) Put(buffer *bytes.Buffer) {
	buffer.Reset()
	pool.buffers.Put(buffer)
}

func NewBufferPool() BufferPool {
	return &bufferPool{}
}
