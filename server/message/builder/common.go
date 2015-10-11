package builder

import (
	"io"
	"git.emersion.fr/saucisse-royale/miko/server/message"
)

func SendPing(w io.Writer) error {
	return write(w, message.Types["ping"])
}

func SendPong(w io.Writer) error {
	return write(w, message.Types["pong"])
}
