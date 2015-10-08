package message

import (
	"io"
)

type IO struct {
	Reader io.Reader
	Writer io.WriteCloser
	BroadcastWriter io.Writer
}
