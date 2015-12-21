package builder

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"io"
	"log"
)

func SendLoginResp(w io.Writer, code message.LoginResponseCode, t message.Tick) error {
	data := []interface{}{
		message.Types["login_response"],
		code,
	}

	if code == message.LoginResponseCodes["ok"] {
		data = append(data, t)
	}

	return sendAll(w, data)
}

func SendRegisterResp(w io.Writer, code message.RegisterResponseCode) error {
	return sendAll(w, []interface{}{
		message.Types["register_response"],
		code,
	})
}

func SendPlayerJoined(w io.Writer, t message.Tick, id message.EntityId, username string) error {
	return sendAll(w, []interface{}{
		message.Types["meta_action"],
		t,
		id,
		message.MetaActionCodes["player_joined"],
		username,
	})
}
func SendPlayerLeft(w io.Writer, t message.Tick, id message.EntityId) error {
	return sendAll(w, []interface{}{
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

func SendChunkUpdate(w io.Writer, t message.Tick, blk *message.Block) error {
	lock(w)
	defer unlock(w)

	if err := write(w, message.Types["chunk_update"]); err != nil {
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

func SendChunksUpdate(w io.Writer, t message.Tick, blks []*message.Block) error {
	lock(w)
	defer unlock(w)

	err := writeAll(w, []interface{}{
		message.Types["chunks_update"],
		t,
		uint16(len(blks)),
	})
	if err != nil {
		return err
	}

	for _, blk := range blks {
		if err := WriteBlock(w, blk); err != nil {
			return err
		}
	}

	return nil
}

func SendEntityCreate(w io.Writer, t message.Tick, entity *message.Entity) error {
	lock(w)
	defer unlock(w)

	if err := write(w, message.Types["entity_create"]); err != nil {
		return err
	}
	if err := write(w, t); err != nil {
		return err
	}

	return writeEntityUpdateBody(w, entity, message.NewFilledEntityDiff(true))
}

func SendEntitiesUpdate(w io.Writer, t message.Tick, entities []*message.Entity, diffs []*message.EntityDiff) error {
	lock(w)
	defer unlock(w)

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

		if err := writeEntityUpdateBody(w, entity, diff); err != nil {
			return err
		}
	}

	return nil
}

func SendEntityDestroy(w io.Writer, t message.Tick, id message.EntityId) error {
	return sendAll(w, []interface{}{
		message.Types["entity_destroy"],
		t,
		id,
	})
}

func SendActionsDone(w io.Writer, t message.Tick, actions []*message.Action) error {
	lock(w)
	defer unlock(w)

	err := writeAll(w, []interface{}{
		message.Types["actions_done"],
		t,
		uint16(len(actions)),
	})
	if err != nil {
		return err
	}

	for _, action := range actions {
		err := writeAll(w, []interface{}{
			action.Initiator,
			action.Id,
		})
		if err != nil {
			return err
		}

		err = writeAll(w, action.Params)
		if err != nil {
			return err
		}
	}

	return nil
}

func SendChatReceive(w io.Writer, t message.Tick, username string, msg string) error {
	return sendAll(w, []interface{}{
		message.Types["chat_receive"],
		t,
		username,
		msg,
	})
}

func SendConfig(w io.Writer, config message.Config) error {
	lock(w)
	defer unlock(w)

	if err := write(w, message.Types["config"]); err != nil {
		return err
	}

	return writeAll(w, config.Export())
}

func SendEntityIdChange(w io.Writer, oldId message.EntityId, newId message.EntityId) error {
	return sendAll(w, []interface{}{
		message.Types["entity_id_change"],
		oldId,
		newId,
	})
}

func SendEntitiesDiffToClients(w io.Writer, t message.Tick, pool *message.EntityDiffPool) error {
	if pool.IsEmpty() {
		return nil
	}
	log.Println("Sending entities diff:", pool)

	// TODO: broadcast only to clients who need it

	// Created entities
	for _, entity := range pool.Created {
		err := SendEntityCreate(w, t, entity)
		if err != nil {
			return err
		}
	}

	// Updated entities
	if len(pool.Updated) > 0 {
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
