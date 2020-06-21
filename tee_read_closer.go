package streams

import (
	"io"
)

var (
	// compile time checking of io.ReadCloser compliance
	_ io.ReadCloser = new(TeeReadCloser)
)

// TeeReadCloser wraps an io.TeeReader and automatically handles the close of the
// underlying reader and writer if they happen to be an io.Closer
type TeeReadCloser struct {
	io.Reader
	close func() error
}

// NewTeeReadCloser returns a new TeeReadCloser
func NewTeeReadCloser(r io.Reader, w io.Writer) *TeeReadCloser {
	return &TeeReadCloser{
		Reader: io.TeeReader(r, w),
		close: func() error {
			if err := closeIfCloser(r); err != nil {
				return err
			}

			return closeIfCloser(w)
		},
	}
}

// Close implements io.Closer interface
func (trc *TeeReadCloser) Close() error {
	return trc.close()
}
