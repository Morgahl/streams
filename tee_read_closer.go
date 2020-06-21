package streams

import (
	"io"
)

// TeeReadCloser wraps an io.TeeReader and handles automatically close of the
// underlying reader and writer when io.EOF is returned from the source reader
type TeeReadCloser struct {
	tee   io.Reader
	close func() error
}

// NewTeeReadCloser returns a new TeeReadCloser
func NewTeeReadCloser(r io.ReadCloser, w io.WriteCloser) (io.ReadCloser, error) {
	return &TeeReadCloser{
		tee: io.TeeReader(r, w),
		close: func() error {
			if err := r.Close(); err != nil {
				return err
			}
			return w.Close()
		},
	}, nil
}

// Read implements io.Reader interface
func (trc *TeeReadCloser) Read(p []byte) (n int, err error) {
	if n, err = trc.tee.Read(p); err == io.EOF {
		err = trc.close()
	}
	return
}

// Close implements io.Closer interface
func (trc *TeeReadCloser) Close() error {
	return trc.close()
}
