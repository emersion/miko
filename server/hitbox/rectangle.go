package hitbox

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
)

type Rectangle struct {
	center *terrain.Position
	width  float64
	height float64
}

func (hb *Rectangle) bounds() [2]*terrain.Position {
	x1 := hb.center.X - hb.width/2
	y1 := hb.center.Y - hb.height/2
	x2 := x1 + hb.width
	y2 := y1 + hb.height

	return [2]*terrain.Position{
		&terrain.Position{x1, y1},
		&terrain.Position{x2, y2},
	}
}

func (hb *Rectangle) Contour() []*terrain.Position {
	bounds := hb.bounds()

	return []*terrain.Position{
		bounds[0],
		&terrain.Position{bounds[1].X, bounds[0].Y},
		bounds[1],
		&terrain.Position{bounds[0].X, bounds[1].Y},
	}
}

func NewRectangle(center *terrain.Position, width, height float64) *Rectangle {
	return &Rectangle{center, width, height}
}
