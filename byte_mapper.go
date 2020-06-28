package streams

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
