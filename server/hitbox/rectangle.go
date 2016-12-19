package hitbox

import (
	"math"

	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
)

type Rectangle struct {
	width  float64
	height float64
}

func (hb *Rectangle) bounds(center *terrain.Position) [2]*terrain.Position {
	x1 := center.X - hb.width/2
	y1 := center.Y - hb.height/2
	x2 := x1 + hb.width
	y2 := y1 + hb.height

	return [2]*terrain.Position{
		&terrain.Position{x1, y1},
		&terrain.Position{x2, y2},
	}
}

func (hb *Rectangle) Contour(center *terrain.Position) []*terrain.Position {
	bounds := hb.bounds(center)

	return []*terrain.Position{
		bounds[0],
		&terrain.Position{bounds[1].X, bounds[0].Y},
		bounds[1],
		&terrain.Position{bounds[0].X, bounds[1].Y},
	}
}

func (hb *Rectangle) Footprint(center *terrain.Position, angle float64) []*terrain.Position {
	contour := hb.Contour(center)

	angle = math.Remainder(angle, math.Pi)
	if angle == 0 {
		return []*terrain.Position{contour[1], contour[2]}
	} else if angle < math.Pi/2 { // Between 0 and pi/2
		return []*terrain.Position{contour[0], contour[2]}
	} else if angle == math.Pi/2 {
		return []*terrain.Position{contour[0], contour[1]}
	} else { // Between pi/2 and pi
		return []*terrain.Position{contour[1], contour[3]}
	}
}

func (hb *Rectangle) intersects(center *terrain.Position, other Hitbox, otherCenter *terrain.Position) (intersects bool, err error) {
	switch o := other.(type) {
	case *Rectangle:
		b1 := hb.bounds(center)
		b2 := o.bounds(otherCenter)
		intersects = (b1[0].X <= b2[1].X && b1[1].X >= b2[0].X &&
			b1[0].Y <= b2[1].Y && b1[1].Y >= b2[0].Y)
	case *Circle:
		// See http://stackoverflow.com/a/402010
		bounds := hb.bounds(center)

		dx := math.Abs(otherCenter.X - bounds[0].X)
		dy := math.Abs(otherCenter.Y - bounds[0].Y)

		if dx > hb.width/2+o.radius || dy > hb.height/2+o.radius {
			intersects = false
		} else if dx <= hb.width/2 || dy <= hb.height/2 {
			intersects = true
		} else {
			dc := math.Pow(dx-hb.width/2, 2) + math.Pow(dy-hb.height/2, 2)
			intersects = (dc <= math.Pow(o.radius, 2))
		}
	default:
		err = errUnsupportedHitbox
	}
	return
}

func NewRectangle(width, height float64) *Rectangle {
	return &Rectangle{width, height}
}
