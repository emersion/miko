package handler

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"io"
	"log"
)

func ReadBlock(r io.Reader) *message.Block {
	blk := &message.Block{}

	var defaultType message.PointType
	var size uint16
	read(r, &blk.X)
	read(r, &blk.Y)
	read(r, &defaultType)
	read(r, &size)

	blk.Points = new(message.BlockPoints)
	blk.Fill(defaultType)

	var x, y message.PointCoord
	var t message.PointType
	for i := 0; i < int(size); i++ {
		read(r, &x)
		read(r, &y)
		read(r, &t)

		blk.Points[x][y] = t
	}

	return blk
}

func ReadActionDone(r io.Reader) (action *message.Action) {
	action = message.NewAction()

	read(r, &action.Initiator)
	read(r, &action.Id)
	// TODO: action params

	return
}

func ReadLoginResponse(r io.Reader) (t message.Tick, code message.LoginResponseCode) {
	read(r, &code)

	if code == message.LoginResponseCodes["ok"] {
		read(r, &t)
	}

	return
}

func ReadRegisterResponse(r io.Reader) (code message.RegisterResponseCode) {
	read(r, &code)
	return
}

func ReadMetaAction(r io.Reader) (t message.Tick, entityId message.EntityId, code message.MetaActionCode, username string) {
	read(r, &t)
	read(r, &entityId)
	read(r, &code)

	if code == message.MetaActionCodes["player_joined"] {
		username = readString(r)
	}

	return
}

func ReadChunkUpdate(r io.Reader) (t message.Tick, blk *message.Block) {
	read(r, &t)
	blk = ReadBlock(r)
	return
}

func ReadChunksUpdate(r io.Reader) (t message.Tick, blks []*message.Block) {
	read(r, &t)

	var size uint16
	read(r, &size)

	blks = make([]*message.Block, size)
	for i := 0; i < int(size); i++ {
		blks[i] = ReadBlock(r)
	}
	return
}

func ReadEntityCreate(r io.Reader) (t message.Tick, entity *message.Entity) {
	read(r, &t)
	entity, _ = ReadEntity(r)
	return
}

func ReadEntitiesUpdate(r io.Reader) (t message.Tick, entities []*message.Entity, diffs []*message.EntityDiff) {
	read(r, &t)

	var size uint16
	read(r, &size)

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
	read(r, &t)
	read(r, &id)
	return
}

func ReadActionsDone(r io.Reader) (t message.Tick, actions []*message.Action) {
	read(r, &t)

	var size uint16
	read(r, &size)

	actions = make([]*message.Action, size)
	for i := 0; i < int(size); i++ {
		actions[i] = ReadActionDone(r)
	}

	return
}

func ReadChatReceive(r io.Reader) (t message.Tick, username, msg string) {
	read(r, &t)
	username = readString(r)
	msg = readString(r)
	return
}

func ReadConfig(r io.Reader) (config *message.Config) {
	config = &message.Config{}
	read(r, &config.MaxRollbackTicks)
	read(r, &config.DefaultPlayerSpeed)
	read(r, &config.PlayerBallCooldown)
	read(r, &config.DefaultBallSpeed)
	read(r, &config.DefaultBallLifespan)
	return
}

func ReadEntityIdChange(r io.Reader) (oldId message.EntityId, newId message.EntityId) {
	read(r, &oldId)
	read(r, &newId)
	return
}

var clientHandlers = &map[message.Type]TypeHandler{
	message.Types["exit"]: func(ctx *message.Context, io *message.IO) error {
		code := ReadExit(io)
		log.Println("Server closed connection, reason:", code)

		return io.Close()
	},
	message.Types["login_response"]: func(ctx *message.Context, io *message.IO) error {
		t, code := ReadLoginResponse(io)

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
		read(io, &size)

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
		read(io, &entityId)
		ctx.Entity.Delete(entityId, t)
		ctx.Entity.Flush()

		log.Println("Entity destroyed:", entityId)
		return nil
	},
	message.Types["actions_done"]: func(ctx *message.Context, io *message.IO) error {
		readTick(io) // TODO: properly handle this tick

		var size uint16
		read(io, &size)

		for i := 0; i < int(size); i++ {
			action := ReadActionDone(io)
			// TODO: do something with this action
			log.Println("Received action with ID", action.Id, "from", action.Initiator)
		}

		return nil
	},
	message.Types["chat_receive"]: func(ctx *message.Context, io *message.IO) error {
		username := readString(io)
		msg := readString(io)

		log.Println("Chat:", username, msg)
		return nil
	},
}
