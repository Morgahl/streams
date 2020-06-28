package convert

import (
	"io"

	"github.com/curlymon/streams"
)

const (
	scratchSize = 1024
)

type Converter interface {
	// Convert converts src into dst, returning the actual number of bytes written to dst.
	Convert(dst, src []byte) (int, error)

	// ChunkSize returns an integer to use for scratch space.
	//
	// This should always return in the range: 1 <= return <= n.
	// A return > n will result in a panic.
	// A return <= 0 will result in an infinite loop or a panic.
	ChunkSize(n int) int
}

func ConvertingInterceptor(converter Converter) streams.Interceptor {
	return &convertingInterceptor{converter: converter}
}

type convertingInterceptor struct {
	converter Converter
	scratch   [scratchSize]byte
}

func (ci *convertingInterceptor) InterceptWrite(w io.Writer, p []byte) (n int, err error) {
	chunkSize := ci.converter.ChunkSize(len(ci.scratch))
	for len(p) > 0 && err == nil {
		if len(p) < chunkSize {
			chunkSize = len(p)
		}

		var converted, _ int
		if converted, err = ci.converter.Convert(ci.scratch[:], p[:chunkSize]); err != nil {
			continue
		}

		_, err = w.Write(ci.scratch[:converted])
		n += chunkSize
		p = p[chunkSize:]
	}
	return n, err
}

func (ci *convertingInterceptor) InterceptRead(r io.Reader, p []byte) (n int, err error) {
	chunkSize := ci.converter.ChunkSize(len(ci.scratch))
	for len(p) > 0 && err == nil {
		var read, converted int
		var rErr error
		read, err = r.Read(ci.scratch[:chunkSize])
		switch err {
		default:
			continue
		case io.EOF:
			rErr = err
		case nil:
			// fallthrough
		}

		converted, err = ci.converter.Convert(p, ci.scratch[:read])
		n += converted
		p = p[converted:]
		if rErr != nil {
			err = rErr
		}
	}
	return n, err
}
