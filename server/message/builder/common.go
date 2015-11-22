// Provides functions to send messages to remotes.
package builder

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"io"
	"sort"
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

func sendEntityUpdateBody(w io.Writer, entity *message.Entity, diff *message.EntityDiff) error {
	if err := write(w, entity.Id); err != nil {
		return err
	}
	if err := write(w, diff.GetBitfield()); err != nil {
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

		attrIds := make([]int, len(entity.Attributes))
		i := 0
		for id, _ := range entity.Attributes {
			attrIds[i] = int(id)
			i++
		}
		sort.Ints(attrIds)

		for _, id := range attrIds {
			attrId := message.EntityAttrId(id)
			data = append(data, attrId, entity.Attributes[attrId])
		}
	}

	return writeAll(w, data)
}
