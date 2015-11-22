package hitbox

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
)

type Null struct{}

func (hb *Null) Contour(other Hitbox) []*entity.Position {
	return []*entity.Position{}
}
