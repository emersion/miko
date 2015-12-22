package hitbox

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
)

type Null struct{}

func (hb *Null) Contour() []*terrain.Position {
	return []*terrain.Position{}
}

func (hb *Null) intersects(other Hitbox) (intersects bool, err error) {
	intersects = false
	return
}

func NewNull() *Null {
	return &Null{}
}
