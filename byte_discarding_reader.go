package streams

import "io"

var (
	// compile time checking of io.ReadCloser compliance
	_ io.Reader = new(ByteDiscardingReader)
)

// ByteDiscardingReader takes a reader and a filter function that
// tests each byte keeping only those that return as true from the filter function
type ByteDiscardingReader struct {
	r          io.Reader
	filterFunc func(byte) bool
}

// NewByteDiscardingReader returns a new ByteDiscardingReader
func NewByteDiscardingReader(r io.Reader, filterFunc func(byte) bool) *ByteDiscardingReader {
	return &ByteDiscardingReader{
		r:          r,
		filterFunc: filterFunc,
	}
}

// Read implements io.Reader interface; filtering each byte with the configured function
func (bdr *ByteDiscardingReader) Read(p []byte) (n int, err error) {
	if n, err = bdr.r.Read(p); err != nil && err != io.EOF {
		return 0, err
	}

	newP := p[:0] // this way we reuse the underlying buffer in an io.Reader smart way
	for _, byt := range p {
		if bdr.filterFunc(byt) {
			newP = append(newP, byt)
		}
	}

	// report lenght of newP as this is the actual fully filtered buffer
	return len(newP), err
}
