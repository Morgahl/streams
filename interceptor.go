package streams

import "io"

type Interceptor interface {
	InterceptRead(io.Reader, []byte) (int, error)
	InterceptWrite(io.Writer, []byte) (int, error)
}

type InterceptCloser interface {
	Interceptor
	InterceptClose(io.Closer) error
}

type BytesCounter struct {
	bytes uint64
}

func (bc *BytesCounter) InterceptWrite(w io.Writer, p []byte) (n int, err error) {
	n, err = w.Write(p)
	bc.bytes += uint64(n)
	return
}

func (bc *BytesCounter) InterceptRead(w io.Reader, p []byte) (n int, err error) {
	n, err = w.Read(p)
	bc.bytes += uint64(n)
	return
}

func (bc *BytesCounter) Count() uint64 {
	return bc.bytes
}

// // ByteMapperInterceptor remaps all intercepted bytes based on the passed ByteMapperFunc. It should
// // be safe to use either a statefull or idempotent function in this. However you should avoid reuse
// // of a stateful ByteMapperFunc as correct behvior is difficult and error prone to implement.
// func ByteMapperInterceptor(byteMapperFunc ByteMapperFunc) InterceptorFunc {
// 	return func(rwFunc ReaderWriterFunc, isReader bool, p []byte) (n int, err error) {
// 		if !isReader {
// 			for i, byt := range p {
// 				p[i] = byteMapperFunc(byt)
// 			}

// 			return rwFunc(p)
// 		}

// 		if n, err = rwFunc(p); err != nil {
// 			return // n, err
// 		}
// 		for i, byt := range p[:n] {
// 			p[i] = byteMapperFunc(byt)
// 		}

// 		return // n, err
// 	}
// }

// // CompiledByteMapperInterceptor prebuilds an array based on the output of the passed
// // ByteMapperFunc, using this array as a fast lookup to remap all intercepted bytes. However as this
// // assumes that the passed ByteMapperFunc is idempotent or stateless in nature, it has undefined
// // behvaior if this is not the case.
// func CompiledByteMapperInterceptor(byteMapperFunc ByteMapperFunc) InterceptorFunc {
// 	return ByteMapperInterceptor(CompileByteMapperFunc(byteMapperFunc))
// }

// // ByteFilterInterceptor filters all intercepted bytes based on the passed ByteFilterFunc. It should
// // be safe to use either a statefull or idempotent function in this. However you should avoid reuse
// // of a stateful ByteFilterFunc as correct behvior is difficult and error prone to implement.
// func ByteFilterInterceptor(bytefilterFunc ByteFilterFunc) InterceptorFunc {
// 	return func(rwFunc ReaderWriterFunc, isReader bool, p []byte) (n int, err error) {
// 		if !isReader {
// 			newP := p[:0] // this way we reuse the underlying buffer safely
// 			for _, byt := range p {
// 				if bytefilterFunc(byt) {
// 					newP = append(newP, byt)
// 				}
// 			}

// 			return rwFunc(newP)
// 		}

// 		if n, err = rwFunc(p); err != nil {
// 			return // n, err
// 		}

// 		newP := p[:0] // this way we reuse the underlying buffer safely
// 		for _, byt := range p[:n] {
// 			if bytefilterFunc(byt) {
// 				newP = append(newP, byt)
// 			}
// 		}

// 		return len(newP), err
// 	}
// }

// // CompiledByteFilterInterceptor prebuilds an array based on the output of the passed
// // ByteFilterFunc, using this array as a fast lookup to filter all intercepted bytes. However as
// // this assumes that the passed ByteFilterFunc is idempotent or stateless in nature, it has
// // undefined behvaior if this is not the case.
// func CompiledByteFilterInterceptor(bytefilterFunc ByteFilterFunc) InterceptorFunc {
// 	return ByteFilterInterceptor(CompileByteFilterFunc(bytefilterFunc))
// }
