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
	ac := TimeoutCloser{openFunc: openFunc, closeAfterTime: closeAfterTime}
	if err := ac.init(); err != nil {
		return nil, err
	}
	return &ac, nil
}

func (ac *TimeoutCloser) init() error {
	var err error
	ac.c, err = ac.openFunc()
	if err != nil {
		return err
	}
	ac.ctx, ac.cancel = context.WithCancel(context.Background())
	timer := time.AfterFunc(ac.closeAfterTime, func() {
		if err := ac.Close(); err != nil {
			ac.errMtx.Lock()
			defer ac.errMtx.Unlock()
			ac.err = err
		}
	})
	ac.resetTimer = func() { timer.Reset(ac.closeAfterTime) }
	return nil
}

func (ac *TimeoutCloser) Close() error {
	ac.errMtx.Lock()
	defer ac.errMtx.Unlock()
	ac.cancel()
	if ac.err != nil {
		err := ac.err
		ac.err = nil
		return err
	}
	return ac.c.Close()
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
	ac := TimeoutReadCloser{openFunc: openFunc, closeAfterTime: closeAfterTime}
	if err := ac.init(); err != nil {
		return nil, err
	}
	return &ac, nil
}

func (ac *TimeoutReadCloser) init() error {
	var err error
	ac.rc, err = ac.openFunc()
	if err != nil {
		return err
	}
	ac.ctx, ac.cancel = context.WithCancel(context.Background())
	timer := time.AfterFunc(ac.closeAfterTime, func() {
		if err := ac.Close(); err != nil {
			ac.errMtx.Lock()
			defer ac.errMtx.Unlock()
			ac.err = err
		}
	})
	ac.resetTimer = func() { timer.Reset(ac.closeAfterTime) }
	return nil
}

func (ac *TimeoutReadCloser) Close() error {
	ac.errMtx.Lock()
	defer ac.errMtx.Unlock()
	ac.cancel()
	if ac.err != nil {
		err := ac.err
		ac.err = nil
		return err
	}
	return ac.rc.Close()
}

func (ac *TimeoutReadCloser) Read(p []byte) (n int, err error) {
	select {
	case <-ac.ctx.Done():
		ac.init()
	default:
	}
	ac.resetTimer()
	return ac.rc.Read(p)
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
	ac := TimeoutWriteCloser{openFunc: openFunc, closeAfterTime: closeAfterTime}
	if err := ac.init(); err != nil {
		return nil, err
	}
	return &ac, nil
}

func (ac *TimeoutWriteCloser) init() error {
	var err error
	ac.wc, err = ac.openFunc()
	if err != nil {
		return err
	}
	ac.ctx, ac.cancel = context.WithCancel(context.Background())
	timer := time.AfterFunc(ac.closeAfterTime, func() {
		if err := ac.Close(); err != nil {
			ac.errMtx.Lock()
			defer ac.errMtx.Unlock()
			ac.err = err
		}
	})
	ac.resetTimer = func() { timer.Reset(ac.closeAfterTime) }
	return nil
}

func (ac *TimeoutWriteCloser) Close() error {
	ac.errMtx.Lock()
	defer ac.errMtx.Unlock()
	ac.cancel()
	if ac.err != nil {
		err := ac.err
		ac.err = nil
		return err
	}
	return ac.wc.Close()
}

func (ac *TimeoutWriteCloser) Write(p []byte) (n int, err error) {
	select {
	case <-ac.ctx.Done():
		ac.init()
	default:
	}
	ac.resetTimer()
	return ac.wc.Write(p)
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

func (ac *TimeoutReadWriteCloser) init() error {
	var err error
	ac.rw, err = ac.openFunc()
	if err != nil {
		return err
	}
	ac.ctx, ac.cancel = context.WithCancel(context.Background())
	timer := time.AfterFunc(ac.closeAfterTime, func() {
		if err := ac.Close(); err != nil {
			ac.errMtx.Lock()
			defer ac.errMtx.Unlock()
			ac.err = err
		}
	})
	ac.resetTimer = func() { timer.Reset(ac.closeAfterTime) }
	return nil
}

func (ac *TimeoutReadWriteCloser) Close() error {
	ac.errMtx.Lock()
	defer ac.errMtx.Unlock()
	ac.cancel()
	if ac.err != nil {
		err := ac.err
		ac.err = nil
		return err
	}
	return ac.rw.Close()
}

func (ac *TimeoutReadWriteCloser) Read(p []byte) (n int, err error) {
	select {
	case <-ac.ctx.Done():
		ac.init()
	default:
	}
	ac.resetTimer()
	return ac.rw.Read(p)
}

func (ac *TimeoutReadWriteCloser) Write(p []byte) (n int, err error) {
	select {
	case <-ac.ctx.Done():
		ac.init()
	default:
	}
	ac.resetTimer()
	return ac.rw.Write(p)
}
