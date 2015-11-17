package entity

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

type Entity struct {
	Id       message.EntityId
	Position *Position
	Speed    *Speed
}

func (e *Entity) ToMessage() *message.Entity {
	return &message.Entity{
		Id:       e.Id,
		Position: e.Position.ToMessage(),
		Speed:    e.Speed.ToMessage(),
	}
}

func (e *Entity) ApplyDiff(d *message.EntityDiff, src *Entity) {
	if d.Position {
		e.Position = src.Position
	}

	if d.SpeedNorm || d.SpeedAngle && e.Speed == nil {
		e.Speed = &Speed{}
	}
	if d.SpeedNorm {
		e.Speed.Norm = src.Speed.Norm
	}
	if d.SpeedAngle {
		e.Speed.Angle = src.Speed.Angle
	}
}

func New() *Entity {
	return &Entity{
		Position: &Position{},
		Speed:    &Speed{},
	}
}

func NewFromMessage(src *message.Entity) *Entity {
	return &Entity{
		Id:       src.Id,
		Position: NewPositionFromMessage(src.Position),
		Speed:    NewSpeedFromMessage(src.Speed),
	}
}
