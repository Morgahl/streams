package streams

import "io"

var (
	// compile time checking of io.Reader compliance
	_ io.Reader = new(Reader)
)

// Reader takes an io.Reader and an Interceptor that observes and intercepts the bytes read. It can
// optionally modify the intercepted bytes in the slice reporting back the new read length and error
// state.
type Reader struct {
	r  io.Reader
	in Interceptor
}

// NewReader returns a new Reader utilizaing the passed Interceptor.
func NewReader(r io.Reader, in Interceptor) *Reader {
	return &Reader{
		r:  r,
		in: in,
	}
}

// Read implements io.Reader interface; intercepting bytes read with the Interceptor
func (rm *Reader) Read(p []byte) (n int, err error) {
	return rm.in.InterceptRead(rm.r, p)
}
