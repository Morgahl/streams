package streams

import (
	"io"

	"github.com/curlymon/bufioc/file"
)

// TeeFileCloser wraps an io.TeeReader and automatically handles the close of the
// underlying reader and writer if they happen to be an io.Closer
type TeeFileCloser struct {
	r      io.Reader
	fClose closeFunc
	close  closeFunc
}

// NewTeeFileCloser returns a new TeeFileCloser
func NewTeeFileCloser(r io.Reader, path string) (io.ReadCloser, error) {
	f, err := file.NewWriterCreate(path)
	if err != nil {
		return nil, err
	}

	return &TeeFileCloser{
		r:      io.TeeReader(r, f),
		fClose: f.Close,
		close: func() error {
			if err := closeIfCloser(r); err != nil {
				return err
			}

			return f.Close()
		},
	}, nil
}

// Read implements the io.Reader interface, it will close the underlying file
// automatically if io.EOF is recieved
func (tfc *TeeFileCloser) Read(p []byte) (n int, err error) {
	if n, err = tfc.r.Read(p); err == io.EOF {
		if cErr := tfc.fClose(); cErr != nil {
			return n, cErr
		}
	}
	return n, err
}

// Close implements io.Closer interface
func (tfc *TeeFileCloser) Close() error {
	return tfc.close()
}
