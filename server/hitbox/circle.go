package hitbox

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
)

type Circle struct {
	center *terrain.Position
	radius float64
}

func (hb *Circle) Contour() []*terrain.Position {
	contour := []*terrain.Position{}
	// TODO
	return contour
}

func NewCircle(center *terrain.Position, radius float64) *Circle {
	return &Circle{center, radius}
}
