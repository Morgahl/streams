package interceptor

import "io"

// ByteMapperFunc is used to map a byte to another byte. This can be a stateful or stateless
// function. However you should avoid mutating or reuse of a stateful ByteMapperFunc as
// correct behvior is difficult and error prone to perfect.
type ByteMapperFunc func(byte) byte

// CompileByteMapperFunc prebuilds an array based on the passed ByteMapperFunc. Using this
// array to perform a fast lookup remap of bytes. However as this assumes that the passed
// ByteMapperFunc is idempotent or stateless in nature, it has undefined behvaior if this
// is not the case.
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

// ByteFilterFunc is used to filter bytes returning true if the byte should be retained. This
// can be a stateful or stateless function. However you should avoid mutating or reuse of a
// stateful ByteFilterFunc as correct behvior is difficult and error prone to perfect.
type ByteFilterFunc func(byte) bool

// CompileByteFilterFunc prebuilds an array based on the passed ByteFilterFunc. Using this
// array to perform a fast lookup filter of bytes. However as this assumes that the passed
// ByteFilterFunc is idempotent or stateless in nature, it has undefined behvaior if this
// is not the case.
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

// InterceptorFunc is an interceptor function that is used to intercept the bytes read/written;
// recieving a reference to the byte buffer and the results of the Read or Write load that
// occured from the underlying io.Reader or io.Writer. It can optionally modify the
// intercepted bytes in the slice reporting back the new read/write length and error state.
type InterceptorFunc func([]byte, int, error) (int, error)

// ByteMapperInterceptor remaps all intercepted bytes based on the passed ByteMapperFunc.
// It should be safe to use either a statefull or idempotent function in this. However
// you should avoid reuse of a stateful ByteMapperFunc as correct behvior is difficult
// and error prone.
func ByteMapperInterceptor(byteMapperFunc ByteMapperFunc) InterceptorFunc {
	return func(p []byte, n int, err error) (int, error) {
		if err != nil && err != io.EOF {
			return n, err
		}

		for i, byt := range p[:n] {
			p[i] = byteMapperFunc(byt)
		}

		return n, err
	}
}

// CompiledByteMapperInterceptor prebuilds a fast lookup array based on the passed
// ByteMapperFunc. Using this array to remap all intercepted bytes. However as this assumes
// that the passed ByteMapperFunc is idempotent or stateless in nature, it has undefined
// behvaior if this is not the case.
func CompiledByteMapperInterceptor(byteMapperFunc ByteMapperFunc) InterceptorFunc {
	return ByteMapperInterceptor(CompileByteMapperFunc(byteMapperFunc))
}

// ByteFilterInterceptor filters all intercepted bytes based on the passed ByteFilterFunc.
// It should be safe to use either a statefull or idempotent function in this. However
// you should avoid reuse of a stateful ByteFilterFunc as correct behvior is difficult
// and error prone.
func ByteFilterInterceptor(bytefilterFunc ByteFilterFunc) InterceptorFunc {
	return func(p []byte, n int, err error) (int, error) {
		if err != nil && err != io.EOF {
			return n, err
		}

		newP := p[:0] // this way we reuse the underlying buffer safely
		for _, byt := range p {
			if bytefilterFunc(byt) {
				newP = append(newP, byt)
			}
		}

		return len(newP), err
	}
}

// CompiledByteFilterInterceptor prebuilds a fast lookup array based on the passed
// ByteFilterFunc. Using this array to filter all intercepted bytes. However as this assumes
// that the passed ByteFilterFunc is idempotent or stateless in nature, it has undefined
// behvaior if this is not the case.
func CompiledByteFilterInterceptor(bytefilterFunc ByteFilterFunc) InterceptorFunc {
	return ByteFilterInterceptor(CompileByteFilterFunc(bytefilterFunc))
}
