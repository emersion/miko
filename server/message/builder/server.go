package builder

import (
	"io"
	"git.emersion.fr/saucisse-royale/miko/server/message"
)

func SendLoginResp(w io.Writer, code message.LoginResponseCode) error {
	if err := write(w, message.Types["login_response"]); err != nil {
		return err
	}
	return write(w, code)
}

func SendRegisterResp(w io.Writer, code message.RegisterResponseCode) error {
	if err := write(w, message.Types["register_response"]); err != nil {
		return err
	}
	return write(w, code)
}

func SendPlayerJoined(w io.Writer, id message.EntityId, username string) error {
	return writeAll(w, []interface{}{
		message.Types["meta_action"],
		message.MetaActionCodes["player_joined"],
		id,
		username,
	})
}
func SendPlayerLeft(w io.Writer, id message.EntityId) error {
	return writeAll(w, []interface{}{
		message.Types["meta_action"],
		message.MetaActionCodes["player_left"],
		id,
	})
}

func WriteBlock(w io.Writer, blk *message.Block) error {
	// Get the point type the most used
	var defaultType message.PointType
	defaultTypeCount := -1
	typesStats := make(map[message.PointType]int)
	for _, column := range blk.Points {
		for _, ptType := range column {
			if _, present := typesStats[ptType]; !present {
				typesStats[ptType] = 1
			} else {
				typesStats[ptType]++
			}

			if defaultTypeCount < typesStats[ptType] {
				defaultTypeCount = typesStats[ptType]
				defaultType = ptType
			}
		}
	}

	if err := write(w, blk.X); err != nil {
		return err
	}
	if err := write(w, blk.Y); err != nil {
		return err
	}
	if err := write(w, defaultType); err != nil {
		return err
	}
	if err := write(w, uint16(blk.Size() - defaultTypeCount)); err != nil {
		return err
	}

	for i := range blk.Points {
		for j := range blk.Points[i] {
			ptType := blk.Points[i][j]
			if ptType == defaultType {
				continue
			}

			if err := write(w, message.PointCoord(i)); err != nil {
				return err
			}
			if err := write(w, message.PointCoord(j)); err != nil {
				return err
			}
			if err := write(w, ptType); err != nil {
				return err
			}
		}
	}

	return nil
}

func SendTerrainUpdate(w io.Writer, blk *message.Block) error {
	if err := write(w, message.Types["terrain_update"]); err != nil {
		return err
	}
	if err := WriteBlock(w, blk); err != nil {
		return err
	}
	return nil
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

	// TODO: entity.Data

	if err := writeAll(w, data); err != nil {
		return err
	}

	return nil
}

// TODO: check proto
func SendEntityCreate(w io.Writer, entity *message.Entity) error {
	if err := write(w, message.Types["entity_create"]); err != nil {
		return err
	}

	diff := &message.EntityDiff{true, true, true}
	return sendEntityUpdateBody(w, entity, diff)
}

func SendEntitiesUpdate(w io.Writer, entities []*message.Entity, diffs []*message.EntityDiff) error {
	if err := write(w, message.Types["entities_update"]); err != nil {
		return err
	}
	if err := write(w, uint16(len(entities))); err != nil {
		return err
	}

	for i, entity := range entities {
		diff := diffs[i]

		if err := sendEntityUpdateBody(w, entity, diff); err != nil {
			return err
		}
	}

	return nil
}

func SendEntityDestroy(w io.Writer, id message.EntityId) error {
	if err := write(w, message.Types["entity_destroy"]); err != nil {
		return err
	}

	if err := write(w, id); err != nil {
		return err
	}

	return nil
}

func SendChatReceive(w io.Writer, username string, msg string) error {
	if err := write(w, message.Types["chat_receive"]); err != nil {
		return err
	}
	if err := writeString(w, username); err != nil {
		return err
	}
	return writeString(w, msg)
}
