package streams

import "io"

// ByteMapper remaps all intercepted bytes based on the passed ByteMapperFunc. It should be safe
// to use either a statefull or idempotent function in this. However you should avoid reuse of a
// stateful ByteMapperFunc as correct behavior is difficult and error prone to implement.
type ByteMapper struct {
	bmfn ByteMapperFunc
}

func NewByteMapper(bmfn ByteMapperFunc) Interceptor {
	return &ByteMapper{
		bmfn: bmfn,
	}
}

func (bm *ByteMapper) InterceptWrite(w io.Writer, p []byte) (n int, err error) {
	for i, b := range p {
		p[i] = bm.bmfn(b)
	}

	return w.Write(p)
}

func (bm *ByteMapper) InterceptRead(r io.Reader, p []byte) (n int, err error) {
	if n, err = r.Read(p); err != nil && err != io.EOF {
		return
	}

	for i, b := range p[:n] {
		p[i] = bm.bmfn(b)
	}

	return
}

// ByteMapperFunc is used to map a byte to another byte. This can be a stateful or stateless
// function. However you should avoid mutating or reuse of a stateful ByteMapperFunc as correct
// behvior is difficult and error prone to perfect.
type ByteMapperFunc func(byte) byte

// CompileByteMapperFunc prebuilds an array based on the passed ByteMapperFunc. Using this array to
// perform a fast lookup remap of bytes. However as this assumes that the passed ByteMapperFunc is
// idempotent or stateless in nature, it has undefined behvaior if this is not the case.
func CompileByteMapperFunc(byteMapperFunc ByteMapperFunc) ByteMapperFunc {
	// This compiles this into a a very fast array lookup
	var byteMap [256]byte
	for i := range byteMap {
		byteMap[i] = byteMapperFunc(byte(i))
	}

	return func(b byte) byte {
		return byteMap[b]
	}
}

// CompiledByteMapper prebuilds an array based on the output of the passed ByteMapperFunc, using this
// array as a fast lookup to remap all intercepted bytes. However as this assumes that the passed
// ByteMapperFunc is idempotent or stateless in nature, it has undefined behavr if this is not the
// case.
func CompiledByteMapper(byteMapperFunc ByteMapperFunc) Interceptor {
	return NewByteMapper(CompileByteMapperFunc(byteMapperFunc))
}
