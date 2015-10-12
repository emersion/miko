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
		message.Types["player_meta"],
		message.MetaActionCodes["player_joined"],
		id,
		username,
	})
}
func SendPlayerLeft(w io.Writer, id message.EntityId) error {
	return writeAll(w, []interface{}{
		message.Types["player_meta"],
		message.MetaActionCodes["player_left"],
		id,
	})
}

func SendTerrainUpdate(w io.Writer, block *message.Block) error {
	// Get the point type the most used
	var defaultType message.PointType
	defaultTypeCount := -1
	typesStats := make(map[message.PointType]int)
	for _, column := range block.Points {
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

	// Send the response
	if err := write(w, message.Types["terrain_update"]); err != nil {
		return err
	}
	if err := write(w, block.X); err != nil {
		return err
	}
	if err := write(w, block.Y); err != nil {
		return err
	}
	if err := write(w, defaultType); err != nil {
		return err
	}
	if err := write(w, uint16(block.Size() - defaultTypeCount)); err != nil {
		return err
	}

	for x, column := range block.Points {
		for y, ptType := range column {
			if ptType == defaultType {
				continue
			}

			if err := write(w, message.PointCoord(x)); err != nil {
				return err
			}
			if err := write(w, message.PointCoord(y)); err != nil {
				return err
			}
			if err := write(w, ptType); err != nil {
				return err
			}
		}
	}

	return nil
}

func SendEntitiesUpdate(w io.Writer, entities []*message.Entity) error {
	if err := write(w, message.Types["entities_update"]); err != nil {
		return err
	}
	if err := write(w, uint16(len(entities))); err != nil {
		return err
	}

	for _, entity := range entities {
		if err := write(w, entity.Id); err != nil {
			return err
		}

		var bitfield uint8
		var data []interface{}

		if entity.Position != nil {
			bitfield |= 1 << 7
			data = append(data, entity.Position.BX)
			data = append(data, entity.Position.BY)
			data = append(data, entity.Position.X)
			data = append(data, entity.Position.Y)
		}
		if entity.Speed != nil {
			// TODO: null values
			if entity.Speed.Angle != 0 {
				bitfield |= 1 << 6
				data = append(data, entity.Speed.Angle)
			}
			if entity.Speed.Norm != 0 {
				bitfield |= 1 << 5
				data = append(data, entity.Speed.Norm)
			}
		}
		// TODO: entity.Data

		if err := write(w, bitfield); err != nil {
			return err
		}
		if err := writeAll(w, data); err != nil {
			return err
		}
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
