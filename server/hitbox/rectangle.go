package hitbox

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
)

type Rectangle struct {
	center *entity.Position
	width  float64
	height float64
}

func (hb *Rectangle) bounds() [2]*entity.Position {
	x1 := hb.center.X - hb.width/2
	y1 := hb.center.Y - hb.height/2
	x2 := x1 + hb.width
	y2 := y1 + hb.height

	return [2]*entity.Position{
		&entity.Position{x1, y1},
		&entity.Position{x2, y2},
	}
}

func (hb *Rectangle) Contour() []*entity.Position {
	bounds := hb.bounds()

	return []*entity.Position{
		bounds[0],
		&entity.Position{bounds[1].X, bounds[0].Y},
		bounds[1],
		&entity.Position{bounds[0].X, bounds[1].Y},
	}
}

func NewRectangle(center *entity.Position, width, height float64) *Rectangle {
	return &Rectangle{center, width, height}
}
