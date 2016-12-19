// Game-specific data.
package game

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/hitbox"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

const (
	PlayerEntity message.EntityType = iota
	BallEntity
)

const (
	PlaceholderSprite message.Sprite = iota
	PlayerSprite
	BallSprite
)

const (
	TicksLeftAttr   message.EntityAttrId = iota
	HealthAttr
	SenderAttr
	CooldownOneAttr                      = 30000
)

const (
	PointTypeEmpty message.PointType = iota
)

type TicksLeft uint16
type Health uint16
type Cooldown uint16

const (
	ThrowBallAction message.ActionId = iota
)

func GetHitbox(sprite message.Sprite) hitbox.Hitbox {
	switch sprite {
	case PlaceholderSprite:
		return hitbox.NewNull()
	case PlayerSprite:
		return hitbox.NewCircle(10)
	case BallSprite:
		return hitbox.NewCircle(10)
	}
	return nil
}

func Collides(a, b message.EntityType) bool {
	return true
}
