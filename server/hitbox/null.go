package hitbox

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
)

type Null struct{}

func (hb *Null) Contour() []*entity.Position {
	return []*entity.Position{}
}

func (hb *Null) intersects(other Hitbox) (intersects bool, err error) {
	intersects = false
	return
}
