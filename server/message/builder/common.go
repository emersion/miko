package builder

import (
	"io"
	".."
)

func SendPing(w io.Writer) error {
	return write(w, message.Types["ping"])
}

func SendPong(w io.Writer) error {
	return write(w, message.Types["pong"])
}
