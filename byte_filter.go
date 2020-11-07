package streams

import "io"

// ByteFilterr filters all intercepted bytes based on the passed ByteFilterFunc. It should
// be safe to use either a statefull or idempotent function in this. However you should avoid reuse
// of a stateful ByteFilterFunc as correct behavior is difficult and error prone to implement.
type ByteFilter struct {
	bffn ByteFilterFunc
}

func NewByteFilter(bffn ByteFilterFunc) Interceptor {
	return &ByteFilter{
		bffn: bffn,
	}
}

func (bf *ByteFilter) InterceptWrite(w io.Writer, p []byte) (n int, err error) {
	newP := p[:0] // this way we reuse the underlying buffer safely
	for _, b := range p {
		if bf.bffn(b) {
			newP = append(newP, b)
		}
	}

	return w.Write(newP)
}

func (bf *ByteFilter) InterceptRead(r io.Reader, p []byte) (n int, err error) {
	if n, err = r.Read(p); err != nil && err != io.EOF {
		return
	}

	newP := p[:0] // this way we reuse the underlying buffer safely
	for _, b := range p[:n] {
		if bf.bffn(b) {
			newP = append(newP, b)
		}
	}

	return len(newP), err
}

// CompiledByteFilter prebuilds an array based on the output of the passed
// ByteFilterFunc, using this array as a fast lookup to filter all intercepted bytes. However as
// this assumes that the passed ByteFilterFunc is idempotent or stateless in nature, it has
// undefined behvaior if this is not the case.
func CompiledByteFilter(bytefilterFunc ByteFilterFunc) Interceptor {
	return NewByteFilter(CompileByteFilterFunc(bytefilterFunc))
}

// ByteFilterFunc is used to filter bytes returning true if the byte should be retained. This can be
// a stateful or stateless function. However you should avoid mutating or reuse of a stateful
// ByteFilterFunc as correct behvior is difficult and error prone to perfect.
type ByteFilterFunc func(byte) bool

// CompileByteFilterFunc prebuilds an array based on the passed ByteFilterFunc. Using this array to
// perform a fast lookup filter of bytes. However as this assumes that the passed ByteFilterFunc is
// idempotent or stateless in nature, it has undefined behvaior if this is not the case.
func CompileByteFilterFunc(byteFilterFunc ByteFilterFunc) ByteFilterFunc {
	// This compiles this into a a very fast array lookup
	var byteMap [256]bool
	for i := range byteMap {
		byteMap[i] = byteFilterFunc(byte(i))
	}

	return func(b byte) bool {
		return byteMap[b]
	}
}
