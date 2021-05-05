package pools

import (
	"github.com/klauspost/compress/zlib"
	"io"
	"sync"
)

type (
	ZlibPool interface {
		GetReader(src io.Reader) (io.ReadCloser, error)
		PutReader(reader io.ReadCloser)
		GetWriter(dst io.Writer, level int) (*zlib.Writer, error)
		PutWriter(writer *zlib.Writer, level int)
	}

	zlibPool struct {
		readers sync.Pool
		writers map[int]sync.Pool
	}
)

var Zlib = NewZlibPool()

func (pool *zlibPool) GetReader(src io.Reader) (io.ReadCloser, error) {
	if reader := pool.readers.Get(); reader != nil {
		reader := reader.(io.ReadCloser)
		err := reader.(zlib.Resetter).Reset(src, nil)
		return reader, err
	}

	reader, err := zlib.NewReader(src)
	return reader, err
}

func (pool *zlibPool) PutReader(reader io.ReadCloser) {
	_ = reader.Close()
	pool.readers.Put(reader)
}

func (pool *zlibPool) GetWriter(dst io.Writer, level int) (*zlib.Writer, error) {
	if levelPool, ok := pool.writers[level]; ok {
		if writer := levelPool.Get(); writer != nil {
			writer := writer.(*zlib.Writer)
			writer.Reset(dst)
			return writer, nil
		}
	}

	writer, err := zlib.NewWriterLevel(dst, level)
	return writer, err
}

func (pool *zlibPool) PutWriter(writer *zlib.Writer, level int) {
	_ = writer.Close()
	if levelPool, ok := pool.writers[level]; ok {
		levelPool.Put(writer)
	} else {
		levelPool := sync.Pool{}
		levelPool.Put(writer)
		pool.writers[level] = levelPool
	}
}

func NewZlibPool() ZlibPool {
	return &zlibPool{
		readers: sync.Pool{},
		writers: make(map[int]sync.Pool),
	}
}
