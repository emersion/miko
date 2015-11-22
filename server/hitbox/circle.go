package hitbox

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
)

type Circle struct {
	center *entity.Position
	radius float64
}

func (hb *Circle) Contour() []*entity.Position {
	contour := []*entity.Position{}
	// TODO
	return contour
}

func NewCircle(center *entity.Position, radius float64) *Circle {
	return &Circle{center, radius}
}
