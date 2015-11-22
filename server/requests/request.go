package requests

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

type Request interface {
	GetTick() message.AbsoluteTick
}
