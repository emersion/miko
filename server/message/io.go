package message

import (
	"io"
)

// An input/output for a specific client
type IO struct {
	Reader io.Reader
	Writer io.WriteCloser
	BroadcastWriter io.Writer
	Id int
	Version ProtocolVersion
}
