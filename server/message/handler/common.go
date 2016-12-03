package handler

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/builder"
	"io"
	"log"
)

func ReadType(r io.Reader) (t message.Type) {
	Read(r, &t)
	return
}

func readTick(r io.Reader) (t message.Tick) {
	Read(r, &t)
	return
}

func ReadExit(r io.Reader) (code message.ExitCode) {
	Read(r, &code)
	return
}

func ReadEntity(r io.Reader) (*message.Entity, *message.EntityDiff) {
	entity := message.NewEntity()

	var bitfield uint8
	Read(r, &entity.Id, &bitfield)

	diff := message.NewEntityDiffFromBitfield(bitfield)

	if diff.Position {
		Read(r, &entity.Position.BX, &entity.Position.BY)
		Read(r, &entity.Position.X, &entity.Position.Y)
	}
	if diff.SpeedAngle {
		Read(r, &entity.Speed.Angle)
	}
	if diff.SpeedNorm {
		Read(r, &entity.Speed.Norm)
	}

	if diff.Type {
		Read(r, &entity.Type)
	}
	if diff.Sprite {
		Read(r, &entity.Sprite)
	}

	if diff.Attributes {
		var size uint16
		Read(r, &size)

		// TODO: move this somewhere else
		for i := 0; i < int(size); i++ {
			var attrId message.EntityAttrId
			Read(r, &attrId)

			// TODO: do something of the data
			var attrVal interface{}
			switch attrId {
			case 0:
				var ticksLeft uint16
				Read(r, &ticksLeft)
				attrVal = ticksLeft
			case 1:
				var health uint16
				Read(r, &health)
				attrVal = health
			case 2:
				var sender message.EntityId
				Read(r, &sender)
				attrVal = sender
			case 30000:
				var cooldownOne uint16
				Read(r, &cooldownOne)
				attrVal = cooldownOne
			default:
				attrVal = nil
			}

			if attrVal != nil {
				entity.Attributes[attrId] = attrVal
			}
		}
	}

	return entity, diff
}

var commonHandlers = &map[message.Type]TypeHandler{
	message.Types["ping"]: func(ctx *message.Context, conn *message.Conn) error {
		log.Println("Ping received!")
		return builder.SendPong(conn)
	},
	message.Types["pong"]: func(ctx *message.Context, conn *message.Conn) error {
		log.Println("Pong received!")
		return nil
	},
}
