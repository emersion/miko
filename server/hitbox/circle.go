package hitbox

import (
	"errors"
	"math"

	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
)

const CircleContourLen = 4

type Circle struct {
	center *terrain.Position
	radius float64
}

func (hb *Circle) Contour() []*terrain.Position {
	// TODO: change base angle according to direction
	contour := make([]*terrain.Position, CircleContourLen)

	for i := 0; i < CircleContourLen; i++ {
		angle := 2 * math.Pi * float64(i) / CircleContourLen
		x := hb.center.X + hb.radius*math.Cos(angle)
		y := hb.center.Y + hb.radius*math.Sin(angle)
		contour[i] = &terrain.Position{x, y}
	}

	return contour
}

func (hb *Circle) intersects(other Hitbox) (intersects bool, err error) {
	switch o := other.(type) {
	case *Circle:
		d := dist(hb.center, o.center)
		intersects = (d <= hb.radius+o.radius)
	default:
		err = errors.New("Unsupported hitbox")
	}
	return
}

func NewCircle(center *terrain.Position, radius float64) *Circle {
	return &Circle{center, radius}
}
