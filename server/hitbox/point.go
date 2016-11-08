package hitbox

import (
	"errors"
	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
)

type Point struct{}

func (hb *Point) Contour(center *terrain.Position) []*terrain.Position {
	return []*terrain.Position{center}
}

func (hb *Point) Footprint(center *terrain.Position, angle float64) []*terrain.Position {
	return []*terrain.Position{center}
}

func (hb *Point) intersects(center *terrain.Position, other Hitbox, otherCenter *terrain.Position) (intersects bool, err error) {
	switch o := other.(type) {
	case *Point:
		intersects = center.ToMessage().Equals(otherCenter.ToMessage())
	case *Rectangle:
		bounds := o.bounds(otherCenter)
		intersects = (bounds[0].X <= center.X && bounds[0].Y <= center.Y &&
			bounds[1].X >= center.X && bounds[1].Y >= center.Y)
	case *Circle:
		d := dist(center, otherCenter)
		intersects = (d <= o.radius)
	default:
		err = errors.New("Unsupported hitbox")
	}
	return
}

func NewPoint() *Point {
	return &Point{}
}
