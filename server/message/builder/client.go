package builder

import (
	"io"
	"git.emersion.fr/saucisse-royale/miko/server/message"
)

func SendLogin(w io.Writer, username string, password string) error {
	return writeAll(w, []interface{}{
		message.Types["login"],
		username,
		password,
	})
}

func SendRegister(w io.Writer, username string, password string) error {
	return writeAll(w, []interface{}{
		message.Types["register"],
		username,
		password,
	})
}


func SendChatSend(w io.Writer, msg string) error {
	if err := write(w, message.Types["chat_send"]); err != nil {
		return err
	}
	return writeString(w, msg)
}
