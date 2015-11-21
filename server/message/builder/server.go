package builder

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"io"
)

func SendLoginResp(w io.Writer, code message.LoginResponseCode, t message.Tick) error {
	if err := write(w, message.Types["login_response"]); err != nil {
		return err
	}

	if err := write(w, code); err != nil {
		return err
	}

	if code == message.LoginResponseCodes["ok"] {
		if err := write(w, t); err != nil {
			return err
		}
	}

	return nil
}

func SendRegisterResp(w io.Writer, code message.RegisterResponseCode) error {
	if err := write(w, message.Types["register_response"]); err != nil {
		return err
	}
	return write(w, code)
}

func SendPlayerJoined(w io.Writer, t message.Tick, id message.EntityId, username string) error {
	return writeAll(w, []interface{}{
		message.Types["meta_action"],
		t,
		id,
		message.MetaActionCodes["player_joined"],
		username,
	})
}
func SendPlayerLeft(w io.Writer, t message.Tick, id message.EntityId) error {
	return writeAll(w, []interface{}{
		message.Types["meta_action"],
		t,
		id,
		message.MetaActionCodes["player_left"],
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
	if err := write(w, uint16(blk.Size()-defaultTypeCount)); err != nil {
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

func SendTerrainUpdate(w io.Writer, t message.Tick, blk *message.Block) error {
	if err := write(w, message.Types["terrain_update"]); err != nil {
		return err
	}
	if err := write(w, t); err != nil {
		return err
	}
	if err := WriteBlock(w, blk); err != nil {
		return err
	}
	return nil
}

func SendEntityCreate(w io.Writer, t message.Tick, entity *message.Entity) error {
	if err := write(w, message.Types["entity_create"]); err != nil {
		return err
	}
	if err := write(w, t); err != nil {
		return err
	}

	diff := &message.EntityDiff{
		Position:   true,
		SpeedAngle: true,
		SpeedNorm:  true,
	}
	return sendEntityUpdateBody(w, entity, diff)
}

func SendEntitiesUpdate(w io.Writer, t message.Tick, entities []*message.Entity, diffs []*message.EntityDiff) error {
	if err := write(w, message.Types["entities_update"]); err != nil {
		return err
	}
	if err := write(w, t); err != nil {
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

func SendEntityDestroy(w io.Writer, t message.Tick, id message.EntityId) error {
	if err := write(w, message.Types["entity_destroy"]); err != nil {
		return err
	}
	if err := write(w, t); err != nil {
		return err
	}

	if err := write(w, id); err != nil {
		return err
	}

	return nil
}

func SendActionsDone(w io.Writer, t message.Tick, actions []*message.Action) error {
	if err := write(w, message.Types["actions_done"]); err != nil {
		return err
	}
	if err := write(w, t); err != nil {
		return err
	}
	if err := write(w, uint16(len(actions))); err != nil {
		return err
	}

	for _, action := range actions {
		if err := write(w, action.Initiator); err != nil {
			return err
		}

		if err := write(w, action.Id); err != nil {
			return err
		}

		// TODO: action params
	}

	return nil
}

func SendChatReceive(w io.Writer, t message.Tick, username string, msg string) error {
	return writeAll(w, []interface{}{
		message.Types["chat_receive"],
		t,
		username,
		msg,
	})
}

func SendConfig(w io.Writer, config *message.Config) error {
	return writeAll(w, []interface{}{
		message.Types["config"],
		config.MaxRollbackTicks,
		config.DefaultPlayerSpeed,
		config.PlayerBallCooldown,
		config.DefaultBallSpeed,
		config.DefaultBallLifespan,
	})
}

func SendEntitiesDiffToClients(w io.Writer, t message.Tick, pool *message.EntityDiffPool) error {
	// TODO: broadcast only to clients who need it

	// Created entities
	for _, entity := range pool.Created {
		err := SendEntityCreate(w, t, entity)
		if err != nil {
			return err
		}
	}

	// Updated entities
	entities := make([]*message.Entity, len(pool.Updated))
	diffs := make([]*message.EntityDiff, len(pool.Updated))
	i := 0
	for entity, diff := range pool.Updated {
		entities[i] = entity
		diffs[i] = diff
		i++
	}

	err := SendEntitiesUpdate(w, t, entities, diffs)
	if err != nil {
		return err
	}

	// Deleted entities
	for _, entityId := range pool.Deleted {
		err := SendEntityDestroy(w, t, entityId)
		if err != nil {
			return err
		}
	}

	return nil
}
