package hitbox

import (
	"errors"
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
)

type Point struct {
	center *entity.Position
}

func (hb *Point) Contour() []*entity.Position {
	return []*entity.Position{hb.center}
}

func (hb *Point) intersects(other Hitbox) (intersects bool, err error) {
	switch o := other.(type) {
	case *Point:
		intersects = hb.center.ToMessage().Equals(o.center.ToMessage())
	case *Rectangle:
		bounds := o.bounds()
		intersects = (bounds[0].X <= hb.center.X && bounds[0].Y <= hb.center.Y &&
			bounds[1].X >= hb.center.X && bounds[1].Y >= hb.center.Y)
	case *Circle:
		d := dist(hb.center, o.center)
		intersects = (d <= o.radius)
	default:
		err = errors.New("Unsupported hitbox")
	}
	return
}

func NewPoint(center *entity.Position) *Point {
	return &Point{center}
}
