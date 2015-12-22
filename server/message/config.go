package message

import (
	"io"
)

// The server config.
type Config interface {
	io.ReaderFrom
	io.WriterTo
}
