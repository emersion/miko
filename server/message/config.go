package message

import (
	"io"
)

type ConfigBase struct {
	MaxRollbackTicks uint16
	TimeServerPort   uint16
}

// The server config.
type Config interface {
	io.ReaderFrom
	io.WriterTo
}
