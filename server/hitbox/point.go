package hitbox

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
)

type Point struct {
	center *entity.Position
}

func (hb *Point) Contour() []*entity.Position {
	return []*entity.Position{hb.center}
}

func NewPoint(center *entity.Position) *Point {
	return &Point{center}
}
