package streams

import "io"

type closeFunc func() error

// closeIfCloser will call Close() if the passed interface can be cast to io.Closer
func closeIfCloser(v interface{}) error {
	if vc, ok := v.(io.Closer); ok {
		return vc.Close()
	}
	return nil
}
