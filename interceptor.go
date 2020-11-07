package streams

import "io"

type Interceptor interface {
	InterceptRead(io.Reader, []byte) (int, error)
	InterceptWrite(io.Writer, []byte) (int, error)
}

type InterceptCloser interface {
	Interceptor
	InterceptClose(io.Closer) error
}

type BytesCounter struct {
	bytes uint64
}

func (bc *BytesCounter) InterceptWrite(w io.Writer, p []byte) (n int, err error) {
	n, err = w.Write(p)
	bc.bytes += uint64(n)
	return
}

func (bc *BytesCounter) InterceptRead(w io.Reader, p []byte) (n int, err error) {
	n, err = w.Read(p)
	bc.bytes += uint64(n)
	return
}

func (bc *BytesCounter) Count() uint64 {
	return bc.bytes
}
