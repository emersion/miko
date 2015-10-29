package entity

import (
	"git.emersion.fr/saucisse-royale/miko/server/message"
)

type Entity struct {
	Id message.EntityId
	Position *Position
	Speed *Speed
}

func (e *Entity) ToMessage() *message.Entity {
	return &message.Entity{
		Id: e.Id,
		Position: e.Position.ToMessage(),
		Speed: e.Speed.ToMessage(),
	}
}

func New() *Entity {
	return &Entity{
		Position: &Position{},
		Speed: &Speed{},
	}
}
