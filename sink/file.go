package sink

import (
	"io"

	"github.com/curlymon/bufioc/file"
)

func File(r io.Reader, path string) (n int, err error) {
	f, err := file.NewWriterCreate(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	var bN int64
	bN, err = io.Copy(f, r)
	return int(bN), err
}
