package builder

import (
	"io"
	"git.emersion.fr/saucisse-royale/miko/server/message"
)

func SendVersion(w io.Writer) error {
	return writeAll(w, []interface{}{
		message.Types["version"],
		message.CURRENT_VERSION,
	})
}

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

func SendEntityUpdate(w io.Writer, entity *message.Entity, diff *message.EntityDiff) error {
	if err := write(w, message.Types["entity_update"]); err != nil {
		return err
	}

	return sendEntityUpdateBody(w, entity, diff)
}

func SendActionDo(w io.Writer, action *message.Action) error {
	if err := write(w, message.Types["action_do"]); err != nil {
		return err
	}

	if err := write(w, action.Id); err != nil {
		return err
	}

	// TODO: action params

	return nil
}

func SendChatSend(w io.Writer, msg string) error {
	if err := write(w, message.Types["chat_send"]); err != nil {
		return err
	}
	return writeString(w, msg)
}

func SendEntitiesDiffToServer(w io.Writer, pool *message.EntityDiffPool) error {
	// Updated entities
	for entity, diff := range pool.Updated {
		err := SendEntityUpdate(w, entity, diff)
		if err != nil {
			return err
		}
	}

	return nil
}
