package hitbox

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
)

type Rectangle struct {
	center *entity.Position
	width  float64
	height float64
}

func (hb *Rectangle) Contour() []*entity.Position {
	contour := []*entity.Position{}
	// TODO
	return contour
}

func NewRectangle(center *entity.Position, width, height float64) *Rectangle {
	return &Rectangle{center, width, height}
}
