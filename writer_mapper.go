package streams

import "io"

var (
	// compile time checking of io.WriteCloser compliance
	_ io.Writer = new(WriterMapper)
)

// WriterMapper takes a Writer and a mapping function that
// mutates each byte Write with the mapping function
type WriterMapper struct {
	r          io.Writer
	mapperFunc func(byte) byte
}

// NewWriterMapper returns a new WriterMapper
func NewWriterMapper(r io.Writer, mapper func(byte) byte) *WriterMapper {
	return &WriterMapper{
		r:          r,
		mapperFunc: mapper,
	}
}

// Write implements io.Writer interface; mapping each byte with the configured function
func (rm *WriterMapper) Write(p []byte) (n int, err error) {
	if n, err = rm.r.Write(p); err != nil && err != io.EOF {
		return 0, err
	}

	for i := range p[:n] {
		p[i] = rm.mapperFunc(p[i])
	}

	return n, err
}
