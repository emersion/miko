package hitbox

import (
	"errors"
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
