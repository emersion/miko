package hitbox

import (
	"errors"
	"math"

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

func (hb *Rectangle) intersects(other Hitbox) (intersects bool, err error) {
	switch o := other.(type) {
	case *Rectangle:
		b1 := hb.bounds()
		b2 := o.bounds()
		intersects = (b1[0].X <= b2[1].X && b1[1].X >= b2[0].X &&
			b1[0].Y <= b2[1].Y && b1[1].Y >= b2[0].Y)
	case *Circle:
		// See http://stackoverflow.com/a/402010
		bounds := hb.bounds()

		dx := math.Abs(o.center.X - bounds[0].X)
		dy := math.Abs(o.center.Y - bounds[0].Y)

		if dx > hb.width/2+o.radius || dy > hb.height/2+o.radius {
			intersects = false
		} else if dx <= hb.width/2 || dy <= hb.height/2 {
			intersects = true
		} else {
			dc := math.Pow(dx-hb.width/2, 2) + math.Pow(dy-hb.height/2, 2)
			intersects = (dc <= math.Pow(o.radius, 2))
		}
	default:
		err = errors.New("Unsupported hitbox")
	}
	return
}

func NewRectangle(center *terrain.Position, width, height float64) *Rectangle {
	return &Rectangle{center, width, height}
}
