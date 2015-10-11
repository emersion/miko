package builder

import (
	"io"
	"git.emersion.fr/saucisse-royale/miko/server/message"
)

func SendChatSend(w io.Writer, msg string) error {
	if err := write(w, message.Types["chat_send"]); err != nil {
		return err
	}
	return writeString(w, msg)
}
