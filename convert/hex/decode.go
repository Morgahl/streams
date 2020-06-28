package hex

import (
	"encoding/hex"
	"io"

	"github.com/curlymon/streams"
	"github.com/curlymon/streams/convert"
)

var (
	Decoding = Decoder{}
)

func NewHexDecodingWriter(w io.Writer) *streams.Writer {
	return streams.NewWriter(w, convert.ConvertingInterceptor(Decoding))
}

func NewHexDecodingReader(r io.Reader) *streams.Reader {
	return streams.NewReader(r, convert.ConvertingInterceptor(Decoding))
}

var (
	_ convert.Converter = Decoder{}
)

type Decoder struct{}

func (Decoder) Convert(dst, src []byte) (n int, err error) {
	return hex.Decode(dst, src)
}

func (Decoder) ChunkSize(n int) int {
	return n
}
