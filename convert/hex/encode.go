package hex

import (
	"encoding/hex"
	"io"

	"github.com/curlymon/streams"
	"github.com/curlymon/streams/convert"
)

var (
	Encoding = Encoder{}
)

func NewHexEncodingWriter(w io.Writer) *streams.Writer {
	return streams.NewWriter(w, convert.ConvertingInterceptor(Encoding))
}

func NewHexEncodingReader(r io.Reader) *streams.Reader {
	return streams.NewReader(r, convert.ConvertingInterceptor(Encoding))
}

var (
	_ convert.Converter = Encoder{}
)

type Encoder struct{}

func (Encoder) Convert(dst, src []byte) (n int, err error) {
	return hex.Encode(dst, src), nil
}

func (Encoder) ChunkSize(n int) int {
	return n >> 1
}
