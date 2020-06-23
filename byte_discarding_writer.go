package streams

import "io"

var (
	// compile time checking of io.WritCloser compliance
	_ io.Writer = new(ByteDiscardingWriter)
)

// ByteDiscardingWriter takes a Writer and a filter function that
// tests each byte keeping only those that return as true from the filter function
type ByteDiscardingWriter struct {
	r          io.Writer
	filterFunc func(byte) bool
}

// NewByteDiscardingWriter returns a new ByteDiscardingWriter
func NewByteDiscardingWriter(r io.Writer, filterFunc func(byte) bool) *ByteDiscardingWriter {
	return &ByteDiscardingWriter{
		r:          r,
		filterFunc: filterFunc,
	}
}

// Writ implements io.Writer interface; filtering each byte with the configured function
func (bdr *ByteDiscardingWriter) Write(p []byte) (n int, err error) {
	if n, err = bdr.r.Write(p); err != nil && err != io.EOF {
		return 0, err
	}

	newP := p[:0] // this way we reuse the underlying buffer in an io.Writer smart way
	for _, byt := range p {
		if bdr.filterFunc(byt) {
			newP = append(newP, byt)
		}
	}

	// report lenght of newP as this is the actual fully filtered buffer
	return len(newP), err
}
