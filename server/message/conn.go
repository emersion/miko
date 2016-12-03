package message

import (
	"io"
)

type flusher interface {
	Flush() error
}

type State int

const (
	Connected    State = 1 << iota // New client connected
	Accepted                       // Client sent its version and has been accepted
	LoggedIn                       // Client is logged in
	Ready                          // Client is ready to play
	Disconnected                   // Client has disconnected
)

// A connection.
type Conn struct {
	io.Reader
	io.Writer
	closer          io.Closer
	broadcastWriter io.Writer
	Id              int
	Version         ProtocolVersion
	State           State
}

// Implements io.Closer.
func (c *Conn) Close() error {
	c.State = Disconnected

	if closer, ok := c.Writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func (c *Conn) Flush() error {
	if flusher, ok := c.Writer.(flusher); ok {
		return flusher.Flush()
	}
	return nil
}

func (c *Conn) Broadcaster() io.Writer {
	return nil // TODO
}

func NewConn(id int, r io.Reader, w io.Writer) *Conn {
	return &Conn{
		Reader:          r,
		Writer:          w,
		Id:              id,
		State:           Connected,
	}
}
