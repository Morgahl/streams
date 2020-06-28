package streams

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
