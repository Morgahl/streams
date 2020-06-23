package interceptor

import "io"

var (
	// compile time checking of io.Writer compliance
	_ io.Writer = new(Writer)
)

// Writer takes an io.Writer and a mapping function that intercepts the bytes
// writen, recieving a reference to the byte buffer and the results of the Write load that
// occured from the underlying io.Writer. It can optionally modify the intercepted bytes
// in the slice reporting back the new written length and error state.
type Writer struct {
	r               io.Writer
	interceptorFunc InterceptorFunc
}

// NewWriter returns a new Writer
func NewWriter(r io.Writer, interceptorFunc func([]byte, int, error) (int, error)) *Writer {
	return &Writer{
		r:               r,
		interceptorFunc: interceptorFunc,
	}
}

// Write implements io.Writer interface; intercepting bytes writen with the configured function
func (rm *Writer) Write(p []byte) (n int, err error) {
	n, err = rm.r.Write(p)
	return rm.interceptorFunc(p, n, err)
}
