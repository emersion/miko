package message

import (
	"io"
	"bufio"
)

type IO struct {
	Reader *bufio.Reader
	Writer io.Writer
	BroadcastWriter io.Writer
}
