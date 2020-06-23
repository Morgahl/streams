package streams

import (
	"io"

	"github.com/curlymon/bufioc/file"
)

var (
	// compile time checking of io.Reader compliance
	_ io.Reader = new(TeeReaderToFile)
)

// TeeReaderToFile wraps an io.TeeReader, it will close the underlying file
// automatically when io.EOF is recieved
type TeeReaderToFile struct {
	r      io.Reader
	fClose CloseFunc
}

// NewTeeReaderToFile returns a new TeeReaderToFile
func NewTeeReaderToFile(r io.Reader, path string) (io.Reader, error) {
	f, err := file.NewWriterCreate(path)
	if err != nil {
		return nil, err
	}

	return &TeeReaderToFile{
		r:      io.TeeReader(r, f),
		fClose: f.Close,
	}, nil
}

// Read implements the io.Reader interface, it will close the underlying file
// automatically when io.EOF is recieved
func (tfc *TeeReaderToFile) Read(p []byte) (n int, err error) {
	if n, err = tfc.r.Read(p); err == io.EOF {
		if cErr := tfc.fClose(); cErr != nil {
			return n, cErr
		}
	}
	return n, err
}
