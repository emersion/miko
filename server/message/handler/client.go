package handler

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"io"
	"log"
)

func ReadEntity(r io.Reader) (*message.Entity, *message.EntityDiff) {
	entity := message.NewEntity()

	var bitfield uint8
	read(r, &entity.Id)
	read(r, &bitfield)

	diff := message.NewEntityDiffFromBitfield(bitfield)

	if diff.Position {
		read(r, &entity.Position.BX)
		read(r, &entity.Position.BY)
		read(r, &entity.Position.X)
		read(r, &entity.Position.Y)
	}
	if diff.SpeedAngle {
		read(r, &entity.Speed.Angle)
	}
	if diff.SpeedNorm {
		read(r, &entity.Speed.Norm)
	}

	// TODO: entity.Data

	return entity, diff
}

func ReadBlock(r io.Reader) *message.Block {
	blk := &message.Block{}

	var defaultType message.PointType
	var size uint16
	read(r, &blk.X)
	read(r, &blk.Y)
	read(r, &defaultType)
	read(r, &size)

	log.Println("Receiving terrain, size:", size)

	blk.Points = new(message.BlockPoints)
	blk.Fill(defaultType)

	var x, y message.PointCoord
	var t message.PointType
	for i := 0; i < int(size); i++ {
		read(r, &x)
		read(r, &y)
		read(r, &t)

		log.Println(" Point at:", x, y, t)
		blk.Points[x][y] = t
	}

	return blk
}

func ReadActionDone(r io.Reader) (action *message.Action) {
	read(r, &action.Initiator)
	read(r, &action.Id)
	// TODO: action params
	return
}

func ReadLoginResponse(r io.Reader) (code message.LoginResponseCode, t message.Tick) {
	read(r, &code)

	if code == message.LoginResponseCodes["ok"] {
		read(r, &t)
	}

	return
}

func ReadVersionResponse(r io.Reader) (code message.VersionResponseCode) {
	read(r, &code)
	return
}

var clientHandlers = &map[message.Type]TypeHandler{
	message.Types["exit"]: func(ctx *message.Context, io *message.IO) error {
		code := ReadExit(io.Reader)
		log.Println("Server closed connection, reason:", code)

		return io.Writer.Close()
	},
	message.Types["login_response"]: func(ctx *message.Context, io *message.IO) error {
		code, t := ReadLoginResponse(io.Reader)

		if code == message.LoginResponseCodes["ok"] {
			ctx.Clock.Sync(t)
		}

		log.Println("Login response:", code)
		return nil
	},
	message.Types["meta_action"]: func(ctx *message.Context, io *message.IO) error {
		var entityId message.EntityId
		var code message.MetaActionCode

		readTick(io.Reader) // TODO: properly handle this tick
		read(io.Reader, &entityId)
		read(io.Reader, &code)

		if code == message.MetaActionCodes["player_joined"] {
			username := readString(io.Reader)
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
	message.Types["terrain_update"]: func(ctx *message.Context, io *message.IO) error {
		readTick(io.Reader) // TODO: properly handle this tick
		blk := ReadBlock(io.Reader)
		ctx.Terrain.SetBlock(blk)
		return nil
	},
	message.Types["entities_update"]: func(ctx *message.Context, io *message.IO) error {
		t := ctx.Clock.ToAbsoluteTick(readTick(io.Reader))

		var size uint16
		read(io.Reader, &size)

		for i := 0; i < int(size); i++ {
			entity, diff := ReadEntity(io.Reader)
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
		t := ctx.Clock.ToAbsoluteTick(readTick(io.Reader))

		entity, _ := ReadEntity(io.Reader)
		ctx.Entity.Add(entity, t)
		ctx.Entity.Flush()

		log.Println("Received new entity with ID:", entity.Id)
		return nil
	},
	message.Types["entity_destroy"]: func(ctx *message.Context, io *message.IO) error {
		t := ctx.Clock.ToAbsoluteTick(readTick(io.Reader))

		var entityId message.EntityId
		read(io.Reader, &entityId)
		ctx.Entity.Delete(entityId, t)
		ctx.Entity.Flush()

		log.Println("Entity destroyed:", entityId)
		return nil
	},
	message.Types["actions_done"]: func(ctx *message.Context, io *message.IO) error {
		readTick(io.Reader) // TODO: properly handle this tick

		var size uint16
		read(io.Reader, &size)

		for i := 0; i < int(size); i++ {
			action := ReadActionDone(io.Reader)
			// TODO: do something with this action
			log.Println("Received action with ID", action.Id, "from", action.Initiator)
		}

		return nil
	},
	message.Types["chat_receive"]: func(ctx *message.Context, io *message.IO) error {
		username := readString(io.Reader)
		msg := readString(io.Reader)

		log.Println("Chat:", username, msg)
		return nil
	},
	message.Types["version_response"]: func(ctx *message.Context, io *message.IO) error {
		code := ReadVersionResponse(io.Reader)
		log.Println("Version response:", code)
		return nil
	},
}
