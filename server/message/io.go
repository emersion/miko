package message

import (
	"io"
	"sync"
)

// An input/output for a specific client
type IO struct {
	reader          io.Reader
	writer          io.WriteCloser
	broadcastWriter io.Writer
	Id              int
	Version         ProtocolVersion

	locker *sync.Mutex
}

func (io *IO) Read(p []byte) (int, error) {
	return io.reader.Read(p)
}

func (io *IO) Write(p []byte) (int, error) {
	return io.writer.Write(p)
}

func (io *IO) Close() error {
	return io.writer.Close()
}

func (io *IO) Broadcaster() io.Writer {
	return io.broadcastWriter
}

func (io *IO) Lock() {
	io.locker.Lock()
}

func (io *IO) Unlock() {
	io.locker.Unlock()
}

func NewIO(id int, reader io.Reader, writer io.WriteCloser, broadcastWriter io.Writer) *IO {
	return &IO{
		reader:          reader,
		writer:          writer,
		broadcastWriter: broadcastWriter,
		Id:              id,

		locker: &sync.Mutex{},
	}
}
