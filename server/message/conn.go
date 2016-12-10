package message

import (
	"io"
	"sync"
)

type WriteFunc func(io.Writer) error

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

	writer    io.Writer
	locker    sync.Locker
	broadcast func(WriteFunc) error

	Id        int
	Version   ProtocolVersion
	State     State
}

// Implements io.Closer.
func (c *Conn) Close() error {
	c.State = Disconnected

	if closer, ok := c.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func (c *Conn) Write(write WriteFunc) error {
	c.locker.Lock()
	defer c.locker.Unlock()

	if err := write(c.writer); err != nil {
		return err
	}

	if flusher, ok := c.writer.(flusher); ok {
		return flusher.Flush()
	}
	return nil
}

func (c *Conn) Broadcast(write WriteFunc) error {
	return c.broadcast(write)
}

func NewConn(id int, r io.Reader, w io.Writer, brd func (WriteFunc) error) *Conn {
	return &Conn{
		Reader:    r,
		writer:    w,
		locker:    &sync.Mutex{},
		broadcast: brd,
		Id:        id,
		State:     Connected,
	}
}
