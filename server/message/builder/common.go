// Provides functions to send messages to remotes.
//
// Writing and sending are different: sending is thread-safe and ensures that
// the message won't be mixed with another one being sent in another process.
package builder

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"io"
)

func SendPing(w io.Writer) error {
	return Write(w, message.Types["ping"])
}

func SendPong(w io.Writer) error {
	return Write(w, message.Types["pong"])
}

func SendExit(w io.Writer, code message.ExitCode) error {
	return Write(w, message.Types["exit"], code)
}

func writeEntityUpdateBody(w io.Writer, entity *message.Entity, diff *message.EntityDiff) error {
	if err := Write(w, entity.Id, diff.GetBitfield()); err != nil {
		return err
	}

	var data []interface{}
	if diff.Position {
		data = append(data, entity.Position.BX, entity.Position.BY, entity.Position.X, entity.Position.Y)
	}
	if diff.SpeedAngle {
		data = append(data, entity.Speed.Angle)
	}
	if diff.SpeedNorm {
		data = append(data, entity.Speed.Norm)
	}

	if diff.Type {
		data = append(data, entity.Type)
	}
	if diff.Sprite {
		data = append(data, entity.Sprite)
	}

	if diff.Attributes {
		data = append(data, uint16(len(entity.Attributes)))

		for id, val := range entity.Attributes {
			data = append(data, id, val)
		}
	}

	return Write(w, data...)
}
