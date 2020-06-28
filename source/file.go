package source

import (
	"io"

	filec "github.com/curlymon/bufioc/file"
)

func File(path string) (io.Reader, error) {
	f, err := filec.NewReader(path)
	if err != nil {
		return nil, err
	}

	return &file{f: f}, nil
}

type file struct {
	f io.ReadCloser
}

func (f *file) Read(p []byte) (n int, err error) {
	if n, err = f.f.Read(p); err != nil && err == io.EOF {
		if fErr := f.f.Close(); fErr != nil {
			return n, fErr
		}
	}

	return
}
