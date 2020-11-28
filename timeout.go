package autoio

import (
	"context"
	"io"
	sync "sync"
	"time"
)

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

func NewTimeoutWriteCloser(openFunc WriteCloserOpenFunc, closeAfterTime time.Duration) (*TimeoutWriteCloser, error) {
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

func (t *TimeoutWriteCloser) Closed() bool {
	select {
	case <-t.ctx.Done():
		return true
	default:
		return false
	}
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
