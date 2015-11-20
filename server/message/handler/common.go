package handler

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/builder"
	"io"
	"log"
)

func readTick(r io.Reader) (t message.Tick) {
	read(r, &t)
	return
}

func ReadExit(r io.Reader) (code message.ExitCode) {
	read(r, &code)
	return
}

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

	if diff.Type {
		read(r, &entity.Type)
	}
	if diff.Sprite {
		read(r, &entity.Sprite)
	}

	return entity, diff
}

var commonHandlers = &map[message.Type]TypeHandler{
	message.Types["ping"]: func(ctx *message.Context, io *message.IO) error {
		log.Println("Ping received!")
		return builder.SendPong(io.Writer)
	},
	message.Types["pong"]: func(ctx *message.Context, io *message.IO) error {
		log.Println("Pong received!")
		return nil
	},
}
