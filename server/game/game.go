// Game-specific data.
package game

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

const (
	PlayerEntity message.EntityType = 0
	BallEntity                      = 1
)

const (
	PlaceholderSprite message.Sprite = 0
	PlayerSprite                     = 1
	BallSprite                       = 2
)

const (
	TicksLeftAttr   message.EntityAttrId = 0
	HealthAttr                           = 1
	SenderAttr                           = 2
	CooldownOneAttr                      = 30000
)

type TicksLeft uint16
type Health uint16
type Cooldown uint16

const (
	ThrowBallAction message.ActionId = 0
)
