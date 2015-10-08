package message

import (
	"io"
	"bufio"
)

type IO struct {
	Reader *bufio.Reader
	Writer io.WriteCloser
	BroadcastWriter io.Writer
}
