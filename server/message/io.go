package message

import (
	"io"
	"sync"
)

type bufferedWriter interface {
	io.Writer
	Flush() error
}

// An input/output for a specific client
type IO struct {
	reader          io.Reader
	writer          io.Writer
	closer          io.Closer
	broadcastWriter io.Writer
	Id              int
	Version         ProtocolVersion

	locker *sync.Mutex
	locked bool
}

func (io *IO) Read(p []byte) (int, error) {
	return io.reader.Read(p)
}

func (io *IO) Write(p []byte) (int, error) {
	return io.writer.Write(p)
}

func (io *IO) Close() error {
	return io.closer.Close()
}

func (io *IO) Broadcaster() io.Writer {
	return io.broadcastWriter
}

func (io *IO) Lock() {
	io.locker.Lock()
	io.locked = true
}

func (io *IO) Unlock() {
	if w, ok := io.writer.(bufferedWriter); ok {
		w.Flush()
	}
	io.locker.Unlock()
	io.locked = false
}

func (io *IO) Locked() bool {
	return io.locked
}

func NewIO(id int, reader io.Reader, writer io.Writer, closer io.Closer, broadcastWriter io.Writer) *IO {
	return &IO{
		reader:          reader,
		writer:          writer,
		closer:          closer,
		broadcastWriter: broadcastWriter,
		Id:              id,

		locker: &sync.Mutex{},
	}
}
