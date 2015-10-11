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

func SendExit(w io.Writer, code message.ExitCode) error {
	if err := write(w, message.Types["exit"]); err != nil {
		return err
	}
	return write(w, code)
}
