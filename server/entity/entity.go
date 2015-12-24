package entity

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/hitbox"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
)

// An entity.
// Entities properties such as position and speed are stored here with a greater
// precision than in the package message.
type Entity struct {
	Id         message.EntityId
	Type       message.EntityType
	Position   *terrain.Position
	Speed      *terrain.Speed
	Sprite     message.Sprite
	Attributes map[message.EntityAttrId]interface{}
}

// Convert this entity to a message.Entity.
func (e *Entity) ToMessage() *message.Entity {
	ent := &message.Entity{
		Id:     e.Id,
		Type:   e.Type,
		Sprite: e.Sprite,
	}

	if e.Position != nil {
		ent.Position = e.Position.ToMessage()
	}
	if e.Speed != nil {
		ent.Speed = e.Speed.ToMessage()
	}
	if e.Attributes != nil {
		ent.Attributes = e.Attributes
	}

	return ent
}

// Apply a diff from a source entity to a destination entity.
// Changed properties will be copied from the source and overwrite the
// destination's ones.
func (e *Entity) ApplyDiff(d *message.EntityDiff, src *Entity) {
	if d.Type {
		e.Type = src.Type
	}

	if d.Position {
		e.Position = src.Position
	}

	if d.SpeedNorm || d.SpeedAngle && e.Speed == nil {
		e.Speed = &terrain.Speed{}
	}
	if d.SpeedNorm {
		e.Speed.Norm = src.Speed.Norm
	}
	if d.SpeedAngle {
		e.Speed.Angle = src.Speed.Angle
	}

	if d.Sprite {
		e.Sprite = src.Sprite
	}

	if d.Attributes {
		e.Attributes = src.Attributes
	}
}

// Get this entity's hitbox.
// [GAME-SPECIFIC]
func (e *Entity) Hitbox() hitbox.Hitbox {
	switch e.Sprite {
	case 0: // placeholder
		return hitbox.NewNull()
	case 1: // player
		return hitbox.NewCircle(e.Position, 10)
	case 2: // ball
		return hitbox.NewCircle(e.Position, 10)
	}
	return nil
}

// Create a new empty entity.
func New() *Entity {
	return &Entity{
		Position:   &terrain.Position{},
		Speed:      &terrain.Speed{},
		Attributes: map[message.EntityAttrId]interface{}{},
	}
}

// Create an entity from a message.Entity.
func NewFromMessage(src *message.Entity) *Entity {
	return &Entity{
		Id:         src.Id,
		Type:       src.Type,
		Position:   terrain.NewPositionFromMessage(src.Position),
		Speed:      terrain.NewSpeedFromMessage(src.Speed),
		Sprite:     src.Sprite,
		Attributes: src.Attributes,
	}
}
