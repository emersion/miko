package handler

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"io"
	"log"
	"time"
)

func ReadBlock(r io.Reader) *message.Block {
	blk := &message.Block{}

	var defaultType message.PointType
	var size uint16
	Read(r, &blk.X, &blk.Y)
	Read(r, &defaultType)
	Read(r, &size)

	blk.Points = new(message.BlockPoints)
	blk.Fill(defaultType)

	var x, y message.PointCoord
	var t message.PointType
	for i := 0; i < int(size); i++ {
		Read(r, &x, &y, &t)

		blk.Points[x][y] = t
	}

	return blk
}

func ReadActionDone(r io.Reader) (action *message.Action) {
	action = message.NewAction()

	Read(r, &action.Initiator, &action.Id)
	// TODO: action params

	return
}

func ReadLoginResponse(r io.Reader) (code message.LoginResponseCode, tick message.Tick, t time.Time) {
	Read(r, &code)

	if code == message.LoginResponseCodes["ok"] {
		Read(r, &tick)

		var timestamp message.Timestamp
		Read(r, &timestamp)
		t = message.TimestampToTime(timestamp)
	}

	return
}

func ReadRegisterResponse(r io.Reader) (code message.RegisterResponseCode) {
	Read(r, &code)
	return
}

func ReadMetaAction(r io.Reader) (t message.Tick, entityId message.EntityId, code message.MetaActionCode, username string) {
	Read(r, &t, &entityId, &code)

	if code == message.MetaActionCodes["player_joined"] {
		Read(r, &username)
	}

	return
}

func ReadChunkUpdate(r io.Reader) (t message.Tick, blk *message.Block) {
	Read(r, &t)
	blk = ReadBlock(r)
	return
}

func ReadChunksUpdate(r io.Reader) (t message.Tick, blks []*message.Block) {
	Read(r, &t)

	var size uint16
	Read(r, &size)

	blks = make([]*message.Block, size)
	for i := 0; i < int(size); i++ {
		blks[i] = ReadBlock(r)
	}
	return
}

func ReadEntityCreate(r io.Reader) (t message.Tick, entity *message.Entity) {
	Read(r, &t)
	entity, _ = ReadEntity(r)
	return
}

func ReadEntitiesUpdate(r io.Reader) (t message.Tick, entities []*message.Entity, diffs []*message.EntityDiff) {
	Read(r, &t)

	var size uint16
	Read(r, &size)

	entities = make([]*message.Entity, size)
	diffs = make([]*message.EntityDiff, size)

	for i := 0; i < int(size); i++ {
		entity, diff := ReadEntity(r)

		entities[i] = entity
		diffs[i] = diff
	}

	return
}

func ReadEntityDestroy(r io.Reader) (t message.Tick, id message.EntityId) {
	Read(r, &t, &id)
	return
}

func ReadActionsDone(r io.Reader) (t message.Tick, actions []*message.Action) {
	Read(r, &t)

	var size uint16
	Read(r, &size)

	actions = make([]*message.Action, size)
	for i := 0; i < int(size); i++ {
		actions[i] = ReadActionDone(r)
	}

	return
}

func ReadChatReceive(r io.Reader) (t message.Tick, username, msg string) {
	Read(r, &t, &username, &msg)
	return
}

func ReadConfig(r io.Reader, config message.Config) error {
	config.ReadFrom(r)
	return nil
}

func ReadEntityIdChange(r io.Reader) (oldId message.EntityId, newId message.EntityId) {
	Read(r, &oldId, &newId)
	return
}

var clientHandlers = &map[message.Type]TypeHandler{
	message.Types["exit"]: func(ctx *message.Context, io *message.IO) error {
		code := ReadExit(io)
		log.Println("Server closed connection, reason:", code)

		return io.Close()
	},
	message.Types["login_response"]: func(ctx *message.Context, io *message.IO) error {
		code, t, _ := ReadLoginResponse(io)

		if code == message.LoginResponseCodes["ok"] {
			ctx.Clock.Sync(t)
		}

		log.Println("Login response:", code)
		return nil
	},
	message.Types["meta_action"]: func(ctx *message.Context, io *message.IO) error {
		// TODO: properly handle this tick
		_, entityId, code, username := ReadMetaAction(io)

		if code == message.MetaActionCodes["player_joined"] {
			log.Println("Player joined:", entityId, username)

			if username == ctx.Me.Username {
				log.Println("My entity ID:", entityId)
				ctx.Me.Entity = ctx.Entity.Get(entityId)
			}
		}
		if code == message.MetaActionCodes["player_left"] {
			log.Println("Player left:", entityId)
		}

		// TODO

		return nil
	},
	message.Types["chunk_update"]: func(ctx *message.Context, io *message.IO) error {
		t, blk := ReadChunkUpdate(io)
		ctx.Terrain.SetBlock(blk, ctx.Clock.ToAbsoluteTick(t))
		return nil
	},
	message.Types["entities_update"]: func(ctx *message.Context, io *message.IO) error {
		t := ctx.Clock.ToAbsoluteTick(readTick(io))

		var size uint16
		Read(io, &size)

		for i := 0; i < int(size); i++ {
			entity, diff := ReadEntity(io)
			// TODO: do something with entity
			log.Println("Received entity update with ID:", entity.Id)

			ctx.Entity.Update(entity, diff, t)
		}

		// These are changes triggered by the server, do not fill the diff pool
		// TODO: if the pool wasn't empty, do not destroy all of it
		ctx.Entity.Flush()

		return nil
	},
	message.Types["entity_create"]: func(ctx *message.Context, io *message.IO) error {
		t := ctx.Clock.ToAbsoluteTick(readTick(io))

		entity, _ := ReadEntity(io)
		ctx.Entity.Add(entity, t)
		ctx.Entity.Flush()

		log.Println("Received new entity with ID:", entity.Id)
		return nil
	},
	message.Types["entity_destroy"]: func(ctx *message.Context, io *message.IO) error {
		t := ctx.Clock.ToAbsoluteTick(readTick(io))

		var entityId message.EntityId
		Read(io, &entityId)
		ctx.Entity.Delete(entityId, t)
		ctx.Entity.Flush()

		log.Println("Entity destroyed:", entityId)
		return nil
	},
	message.Types["actions_done"]: func(ctx *message.Context, io *message.IO) error {
		readTick(io) // TODO: properly handle this tick

		var size uint16
		Read(io, &size)

		for i := 0; i < int(size); i++ {
			action := ReadActionDone(io)
			// TODO: do something with this action
			log.Println("Received action with ID", action.Id, "from", action.Initiator)
		}

		return nil
	},
	message.Types["chat_receive"]: func(ctx *message.Context, io *message.IO) error {
		var username, msg string
		Read(io, &username)
		Read(io, &msg)

		log.Println("Chat:", username, msg)
		return nil
	},
}
