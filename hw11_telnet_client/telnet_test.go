package main

import (
	"bytes"
	"io"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
}

func TestWrongHost(t *testing.T) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	start := time.Now()
	client := NewTelnetClient("wronghost", 3*time.Second, io.NopCloser(in), out)
	end := time.Now()
	require.Error(t, client.Connect(), "Check no error thrown in wrong host")
	require.Less(t, end.Second()-start.Second(), 3, "Check timeout")
}

func TestTimeoutingHost(t *testing.T) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	start := time.Now()
	client := NewTelnetClient("mysite.ru:8080", 3*time.Second, io.NopCloser(in), out)
	end := time.Now()
	require.Error(t, client.Connect(), "Check no error thrown in wrong host")
	require.Less(t, end.Second()-start.Second(), 3, "Check timeout")
}

func TestWriteToClosedServer(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)

	in := &bytes.Buffer{}
	out := &bytes.Buffer{}

	client := NewTelnetClient(l.Addr().String(), time.Second, io.NopCloser(in), out)
	require.NoError(t, client.Connect())
	defer func() { require.NoError(t, client.Close()) }()
	require.NoError(t, l.Close())

	in.WriteString("XYZ\n")
	err = client.Send()
	require.ErrorAs(t, err, &ErrServerDisconnected)
}

func TestEOFStopsClient(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	defer func() { require.NoError(t, l.Close()) }()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient(l.Addr().String(), time.Second, io.NopCloser(in), out)
		require.NoError(t, client.Connect())
		defer func() { require.Error(t, client.Close()) }()

		in.WriteString("XYZ\n")
		err = client.Send()
		require.NoError(t, err)

		strings.NewReader("").WriteTo(in)
		err = client.Send()
		require.ErrorIs(t, err, io.EOF)
	}()

	go func() {
		defer wg.Done()

		conn, err := l.Accept()
		require.NoError(t, err)
		require.NotNil(t, conn)
		defer func() { require.NoError(t, conn.Close()) }()

		request := make([]byte, 1024)
		n, err := conn.Read(request)
		require.NoError(t, err)
		require.Equal(t, "XYZ\n", string(request)[:n])
	}()

	wg.Wait()
}
