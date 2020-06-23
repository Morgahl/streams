package interceptor

import "io"

var (
	// compile time checking of io.Reader compliance
	_ io.Reader = new(Reader)
)

// Reader takes an io.Reader and a mapping function that intercepts the bytes
// read, recieving a reference to the byte buffer and the results of the Read load that
// occured from the underlying io.Reader. It can optionally modify the intercepted bytes
// in the slice reporting back the new read length and error state.
type Reader struct {
	r               io.Reader
	interceptorFunc InterceptorFunc
}

// NewReader returns a new Reader
func NewReader(r io.Reader, interceptorFunc func([]byte, int, error) (int, error)) *Reader {
	return &Reader{
		r:               r,
		interceptorFunc: interceptorFunc,
	}
}

// Read implements io.Reader interface; intercepting bytes read with the configured function
func (rm *Reader) Read(p []byte) (n int, err error) {
	n, err = rm.r.Read(p)
	return rm.interceptorFunc(p, n, err)
}
