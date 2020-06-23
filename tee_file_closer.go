package streams

import (
	"io"

	"github.com/curlymon/bufioc/file"
)

var (
	// compile time checking of io.Reader compliance
	_ io.Reader = new(TeeFileReader)
)

// TeeFileReader wraps an io.TeeReader and automatically handles the close of the
// underlying reader and writer if they happen to be an io.Closer
type TeeFileReader struct {
	r      io.Reader
	fClose closeFunc
}

// NewTeeFileReader returns a new TeeFileReader
func NewTeeFileReader(r io.Reader, path string) (io.Reader, error) {
	f, err := file.NewWriterCreate(path)
	if err != nil {
		return nil, err
	}

	return &TeeFileReader{
		r:      io.TeeReader(r, f),
		fClose: f.Close,
	}, nil
}

// Read implements the io.Reader interface, it will close the underlying file
// automatically if io.EOF is recieved
func (tfc *TeeFileReader) Read(p []byte) (n int, err error) {
	if n, err = tfc.r.Read(p); err == io.EOF {
		if cErr := tfc.fClose(); cErr != nil {
			return n, cErr
		}
	}
	return n, err
}
