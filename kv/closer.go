package kv

import "io"

type noopCloser struct{}

func (noopCloser) Close() error { return nil }

func getCloser(k KeyValue) io.Closer {
	closer, ok := k.(io.Closer)
	if !ok {
		return noopCloser{}
	}

	return closer
}
