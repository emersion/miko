package hitbox

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
)

type Null struct{}

func (hb *Null) Contour(center *terrain.Position) []*terrain.Position {
	return []*terrain.Position{}
}

func (hb *Null) Footprint(center *terrain.Position) []*terrain.Position {
	return []*terrain.Position{}
}

func (hb *Null) intersects(center *terrain.Position, other Hitbox, otherCenter *terrain.Position) (intersects bool, err error) {
	intersects = false
	return
}

func NewNull() *Null {
	return &Null{}
}
