package builder

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"io"
)

func SendVersion(w io.Writer) error {
	return writeAll(w, []interface{}{
		message.Types["version"],
		message.CurrentVersion,
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

func SendEntityUpdate(w io.Writer, t message.Tick, entity *message.Entity, diff *message.EntityDiff) error {
	if err := write(w, message.Types["entity_update"]); err != nil {
		return err
	}
	if err := write(w, t); err != nil {
		return err
	}

	return writeEntityUpdateBody(w, entity, diff)
}

func SendActionDo(w io.Writer, t message.Tick, action *message.Action) error {
	if err := write(w, message.Types["action_do"]); err != nil {
		return err
	}
	if err := write(w, t); err != nil {
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

func SendEntitiesDiffToServer(w io.Writer, t message.Tick, pool *message.EntityDiffPool) error {
	// Updated entities
	for entity, diff := range pool.Updated {
		err := SendEntityUpdate(w, t, entity, diff)
		if err != nil {
			return err
		}
	}

	return nil
}
