package tap

import (
	"io"

	filec "github.com/curlymon/bufioc/file"
)

// File returns a new
func File(r io.Reader, path string) (io.Reader, error) {
	f, err := filec.NewWriterCreate(path)
	if err != nil {
		return nil, err
	}

	return &file{
		r: io.TeeReader(r, f),
		f: f,
	}, nil
}

// file wraps an io.TeeReader, it will close the underlying file automatically when io.EOF is
// recieved
type file struct {
	r io.Reader
	f io.WriteCloser
}

// Read implements the io.Reader interface, it will close the underlying file
// automatically when io.EOF is recieved
func (tfc *file) Read(p []byte) (n int, err error) {
	if n, err = tfc.r.Read(p); err == io.EOF {
		if cErr := tfc.f.Close(); cErr != nil {
			return n, cErr
		}
	}
	return n, err
}
