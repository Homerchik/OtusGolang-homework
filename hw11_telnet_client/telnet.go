package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var ErrServerDisconnected = errors.New("server disconnected")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type MyClient struct {
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func (c *MyClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("error connecting to %s: %w", c.address, err)
	}
	fmt.Fprint(os.Stderr, "...Connected to ", c.address, "\n")
	c.conn = conn
	return nil
}

func (c *MyClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *MyClient) Send() error {
	buf := make([]byte, 1024)
	n, err := bufio.NewReader(c.in).Read(buf)
	if err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "error reading from input: %v\n", err)
		c.Close()
		return err
	}
	if n > 0 {
		if _, errCon := c.conn.Write(buf[:n]); errCon != nil {
			fmt.Fprintf(os.Stderr, "...Connection was closed by peer\n")
			return ErrServerDisconnected
		}
	}
	if err == io.EOF {
		fmt.Fprintf(os.Stderr, "^D\n...EOF\n")
		c.Close()
	}
	return err
}

func (c *MyClient) Receive() error {
	_, err := bufio.NewReader(c.conn).WriteTo(c.out)
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &MyClient{address: address, timeout: timeout, in: in, out: out}
}
