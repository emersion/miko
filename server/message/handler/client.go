package handler

import (
	"io"
	"log"
	"git.emersion.fr/saucisse-royale/miko/server/message"
)

func readEntityUpdateBody(r io.Reader) *message.Entity {
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

	return entity
}

var clientHandlers = &map[message.Type]TypeHandler{
	message.Types["exit"]: func(ctx *message.Context, io *message.IO) error {
		var code message.ExitCode
		read(io.Reader, &code)
		log.Println("Server closed connection, reason:", code)

		return io.Writer.Close()
	},
	message.Types["login_response"]: func(ctx *message.Context, io *message.IO) error {
		var code message.LoginResponseCode
		read(io.Reader, &code)
		log.Println("Login response:", code)
		return nil
	},
	message.Types["meta_action"]: func(ctx *message.Context, io *message.IO) error {
		var code message.MetaActionCode
		var entityId message.EntityId
		read(io.Reader, &code)
		read(io.Reader, &entityId)
		username := readString(io.Reader)

		log.Println("Player joined/left:", code, entityId, username)
		// TODO

		return nil
	},
	message.Types["terrain_update"]: func(ctx *message.Context, io *message.IO) error {
		var blk message.Block
		var defaultType message.PointType
		var size uint16
		read(io.Reader, &blk.X)
		read(io.Reader, &blk.Y)
		read(io.Reader, &defaultType)
		read(io.Reader, &size)

		log.Println("Receiving terrain, size:", size)

		blk.Points = new(message.BlockPoints)

		var x, y message.PointCoord
		var t message.PointType
		for i := 0; i < int(size); i++ {
			read(io.Reader, &x)
			read(io.Reader, &y)
			read(io.Reader, &t)

			blk.Points[x][y] = t
			log.Println(" Point at:", x, y, t)
		}

		// TODO: do something with terrain

		return nil
	},
	message.Types["entities_update"]: func(ctx *message.Context, io *message.IO) error {
		var size uint16
		read(io.Reader, &size)

		for i := 0; i < int(size); i++ {
			entity := readEntityUpdateBody(io.Reader)
			// TODO: do something with entity
			log.Println("Received entity update with ID:", entity.Id)
		}

		return nil
	},
	message.Types["entity_create"]: func(ctx *message.Context, io *message.IO) error {
		entity := readEntityUpdateBody(io.Reader)

		// TODO: entity.Data

		// TODO: do something with entity

		log.Println("Received new entity with ID:", entity.Id)
		return nil
	},
	message.Types["chat_receive"]: func(ctx *message.Context, io *message.IO) error {
		username := readString(io.Reader)
		msg := readString(io.Reader)

		log.Println("Chat:", username, msg)
		return nil
	},
}
