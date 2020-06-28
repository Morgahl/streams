package streams

import "io"

var (
	// compile time checking of io.Writer compliance
	_ io.Writer = new(Writer)
)

// Writer takes an io.Writer and an Interceptor that observes and intercepts the bytes read. It can
// optionally modify the intercepted bytes in the slice reporting back the new read length and error
// state.
type Writer struct {
	w  io.Writer
	in Interceptor
}

// NewWriter returns a new Writer utilizaing the passed Interceptor.
func NewWriter(w io.Writer, in Interceptor) *Writer {
	return &Writer{
		w:  w,
		in: in,
	}
}

// Write implements io.Writer interface; intercepting bytes read with the Interceptor
func (rm *Writer) Write(p []byte) (n int, err error) {
	return rm.in.InterceptWrite(rm.w, p)
}
