package autoio

import "io"

type FlushCloser interface {
	Flush() error
	io.Closer
}

type AutoFlusher struct {
	f FlushCloser
}

func NewAutoFlushCloser(f FlushCloser) io.Closer {
	return &AutoFlusher{f: f}
}

func (f *AutoFlusher) Close() error {
	flushErr := f.f.Flush()
	closeErr := f.f.Close()
	if flushErr != nil {
		return flushErr
	}
	return closeErr
}

type FlushReadCloser interface {
	Flush() error
	io.ReadCloser
}

type AutoFlushReadCloser struct {
	f FlushReadCloser
}

func NewAutoFlushReadCloser(f FlushReadCloser) io.ReadCloser {
	return &AutoFlushReadCloser{f: f}
}

func (f *AutoFlushReadCloser) Close() error {
	flushErr := f.f.Flush()
	closeErr := f.f.Close()
	if flushErr != nil {
		return flushErr
	}
	return closeErr
}

func (f *AutoFlushReadCloser) Read(b []byte) (int, error) {
	return f.f.Read(b)
}

func NewAutoFlushWriteCloser(f FlushWriteCloser) io.WriteCloser {
	return &AutoFlushWriteCloser{f: f}
}

func (f *AutoFlushWriteCloser) Close() error {
	flushErr := f.f.Flush()
	closeErr := f.f.Close()
	if flushErr != nil {
		return flushErr
	}
	return closeErr
}

func (f *AutoFlushWriteCloser) Write(b []byte) (int, error) {
	return f.f.Write(b)
}

type FlushReadWriteCloser interface {
	Flush() error
	io.ReadWriteCloser
}

type AutoFlushReadWriteCloser struct {
	f FlushReadWriteCloser
}

func NewAutoFlushReadWriteCloser(f FlushReadWriteCloser) io.ReadWriteCloser {
	return &AutoFlushReadWriteCloser{f: f}
}

func (f *AutoFlushReadWriteCloser) Close() error {
	flushErr := f.f.Flush()
	closeErr := f.f.Close()
	if flushErr != nil {
		return flushErr
	}
	return closeErr
}

func (f *AutoFlushReadWriteCloser) Read(b []byte) (int, error) {
	return f.f.Read(b)
}

func (f *AutoFlushReadWriteCloser) Write(b []byte) (int, error) {
	return f.f.Write(b)
}

type FlushWriteCloser interface {
	Flush() error
	io.WriteCloser
}

type AutoFlushWriteCloser struct {
	f FlushWriteCloser
}
