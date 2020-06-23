package streams

import "io"

var (
	// compile time checking of io.Reader compliance
	_ io.Reader = new(ReaderMapper)
)

// ReaderMapper takes a reader and a mapping function that
// mutates each byte read with the mapping function
type ReaderMapper struct {
	r          io.Reader
	mapperFunc func(byte) byte
}

// NewReaderMapper returns a new ReaderMapper
func NewReaderMapper(r io.Reader, mapper func(byte) byte) *ReaderMapper {
	return &ReaderMapper{
		r:          r,
		mapperFunc: mapper,
	}
}

// Read implements io.Reader interface; mapping each byte with the configured function
func (rm *ReaderMapper) Read(p []byte) (n int, err error) {
	if n, err = rm.r.Read(p); err != nil && err != io.EOF {
		return 0, err
	}

	for i := range p[:n] {
		p[i] = rm.mapperFunc(p[i])
	}

	return n, err
}
