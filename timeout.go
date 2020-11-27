package autoio

import (
	"context"
	"io"
	sync "sync"
	"time"
)

type TimeoutCloser struct {
	openFunc       CloserOpenFunc
	closeAfterTime time.Duration
	ctx            context.Context
	cancel         context.CancelFunc
	err            error
	errMtx         sync.RWMutex
	resetTimer     func()
	c              io.Closer
}

type CloserOpenFunc func() (io.Closer, error)

func NewTimeoutCloser(openFunc CloserOpenFunc, closeAfterTime time.Duration) (io.Closer, error) {
	t := TimeoutCloser{openFunc: openFunc, closeAfterTime: closeAfterTime}
	if err := t.init(); err != nil {
		return nil, err
	}
	return &t, nil
}

func (t *TimeoutCloser) init() error {
	var err error
	t.c, err = t.openFunc()
	if err != nil {
		return err
	}
	t.ctx, t.cancel = context.WithCancel(context.Background())
	timer := time.AfterFunc(t.closeAfterTime, func() {
		if err := t.Close(); err != nil {
			t.errMtx.Lock()
			defer t.errMtx.Unlock()
			t.err = err
		}
	})
	t.resetTimer = func() { timer.Reset(t.closeAfterTime) }
	return nil
}

func (t *TimeoutCloser) Close() error {
	t.errMtx.Lock()
	defer t.errMtx.Unlock()
	t.cancel()
	if t.err != nil {
		err := t.err
		t.err = nil
		return err
	}
	return t.c.Close()
}

type TimeoutReadCloser struct {
	openFunc       ReadCloserOpenFunc
	closeAfterTime time.Duration
	ctx            context.Context
	cancel         context.CancelFunc
	err            error
	errMtx         sync.RWMutex
	resetTimer     func()
	rc             io.ReadCloser
}

type ReadCloserOpenFunc func() (io.ReadCloser, error)

func NewTimeoutReadCloser(openFunc ReadCloserOpenFunc, closeAfterTime time.Duration) (io.ReadCloser, error) {
	t := TimeoutReadCloser{openFunc: openFunc, closeAfterTime: closeAfterTime}
	if err := t.init(); err != nil {
		return nil, err
	}
	return &t, nil
}

func (t *TimeoutReadCloser) init() error {
	var err error
	t.rc, err = t.openFunc()
	if err != nil {
		return err
	}
	t.ctx, t.cancel = context.WithCancel(context.Background())
	timer := time.AfterFunc(t.closeAfterTime, func() {
		if err := t.Close(); err != nil {
			t.errMtx.Lock()
			defer t.errMtx.Unlock()
			t.err = err
		}
	})
	t.resetTimer = func() { timer.Reset(t.closeAfterTime) }
	return nil
}

func (t *TimeoutReadCloser) Close() error {
	t.errMtx.Lock()
	defer t.errMtx.Unlock()
	t.cancel()
	if t.err != nil {
		err := t.err
		t.err = nil
		return err
	}
	return t.rc.Close()
}

func (t *TimeoutReadCloser) Read(p []byte) (n int, err error) {
	select {
	case <-t.ctx.Done():
		t.init()
	default:
	}
	t.resetTimer()
	return t.rc.Read(p)
}

type TimeoutWriteCloser struct {
	openFunc       WriteCloserOpenFunc
	closeAfterTime time.Duration
	ctx            context.Context
	cancel         context.CancelFunc
	err            error
	errMtx         sync.RWMutex
	resetTimer     func()
	wc             io.WriteCloser
}

type WriteCloserOpenFunc func() (io.WriteCloser, error)

func NewTimeoutWriteCloser(openFunc WriteCloserOpenFunc, closeAfterTime time.Duration) (io.WriteCloser, error) {
	t := TimeoutWriteCloser{openFunc: openFunc, closeAfterTime: closeAfterTime}
	if err := t.init(); err != nil {
		return nil, err
	}
	return &t, nil
}

func (t *TimeoutWriteCloser) init() error {
	var err error
	t.wc, err = t.openFunc()
	if err != nil {
		return err
	}
	t.ctx, t.cancel = context.WithCancel(context.Background())
	timer := time.AfterFunc(t.closeAfterTime, func() {
		if err := t.Close(); err != nil {
			t.errMtx.Lock()
			defer t.errMtx.Unlock()
			t.err = err
		}
	})
	t.resetTimer = func() { timer.Reset(t.closeAfterTime) }
	return nil
}

func (t *TimeoutWriteCloser) Close() error {
	t.errMtx.Lock()
	defer t.errMtx.Unlock()
	t.cancel()
	if t.err != nil {
		err := t.err
		t.err = nil
		return err
	}
	return t.wc.Close()
}

func (t *TimeoutWriteCloser) Write(p []byte) (n int, err error) {
	select {
	case <-t.ctx.Done():
		t.init()
	default:
	}
	t.resetTimer()
	return t.wc.Write(p)
}

type TimeoutReadWriteCloser struct {
	openFunc       ReadWriteCloserOpenFunc
	closeAfterTime time.Duration
	ctx            context.Context
	cancel         context.CancelFunc
	err            error
	errMtx         sync.RWMutex
	resetTimer     func()
	rw             io.ReadWriteCloser
}

type ReadWriteCloserOpenFunc func() (io.ReadWriteCloser, error)

func NewTimeoutReadWriteCloser(openFunc ReadWriteCloserOpenFunc, closeAfterTime time.Duration) (io.ReadWriteCloser, error) {
	ac := TimeoutReadWriteCloser{openFunc: openFunc, closeAfterTime: closeAfterTime}
	if err := ac.init(); err != nil {
		return nil, err
	}
	return &ac, nil
}

func (t *TimeoutReadWriteCloser) init() error {
	var err error
	t.rw, err = t.openFunc()
	if err != nil {
		return err
	}
	t.ctx, t.cancel = context.WithCancel(context.Background())
	timer := time.AfterFunc(t.closeAfterTime, func() {
		if err := t.Close(); err != nil {
			t.errMtx.Lock()
			defer t.errMtx.Unlock()
			t.err = err
		}
	})
	t.resetTimer = func() { timer.Reset(t.closeAfterTime) }
	return nil
}

func (t *TimeoutReadWriteCloser) Close() error {
	t.errMtx.Lock()
	defer t.errMtx.Unlock()
	t.cancel()
	if t.err != nil {
		err := t.err
		t.err = nil
		return err
	}
	return t.rw.Close()
}

func (t *TimeoutReadWriteCloser) Read(p []byte) (n int, err error) {
	select {
	case <-t.ctx.Done():
		t.init()
	default:
	}
	t.resetTimer()
	return t.rw.Read(p)
}

func (t *TimeoutReadWriteCloser) Write(p []byte) (n int, err error) {
	select {
	case <-t.ctx.Done():
		t.init()
	default:
	}
	t.resetTimer()
	return t.rw.Write(p)
}
