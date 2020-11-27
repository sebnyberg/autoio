package autoio_test

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/sebnyberg/autoio"
	"github.com/stretchr/testify/require"
)

type TestWriteCloser struct {
	Open bool
}

func (c *TestWriteCloser) Write(b []byte) (int, error) {
	fmt.Println(string(b))
	return 1, nil
}

func (c *TestWriteCloser) Close() error {
	fmt.Println("writer was closed")
	c.Open = false
	return nil
}

func Test_TimeoutWriteCloser(t *testing.T) {
	w, _ := autoio.NewTimeoutWriteCloser(func() (io.WriteCloser, error) {
		return &TestWriteCloser{}, nil
	}, 100*time.Millisecond)
	require.Equal(t, false, w.Closed())
	w.Write([]byte("test"))
	require.Equal(t, false, w.Closed())
	time.Sleep(150 * time.Millisecond)
	w.Write([]byte("test"))
	require.Equal(t, true, w.Closed())
}
