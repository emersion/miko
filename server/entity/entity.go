package entity

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

// An entity.
// Entities properties such as position and speed are stored here with a greater
// precision than in the package message.
type Entity struct {
	Id       message.EntityId
	Position *Position
	Speed    *Speed
}

// Convert this entity to a message.Entity.
func (e *Entity) ToMessage() *message.Entity {
	return &message.Entity{
		Id:       e.Id,
		Position: e.Position.ToMessage(),
		Speed:    e.Speed.ToMessage(),
	}
}

// Apply a diff from a source entity to a destination entity.
// Changed properties will be copied from the source and overwrite the
// destination's ones.
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

// Create a new empty entity.
func New() *Entity {
	return &Entity{
		Position: &Position{},
		Speed:    &Speed{},
	}
}

// Create an entity from a message.Entity.
func NewFromMessage(src *message.Entity) *Entity {
	return &Entity{
		Id:       src.Id,
		Position: NewPositionFromMessage(src.Position),
		Speed:    NewSpeedFromMessage(src.Speed),
	}
}
