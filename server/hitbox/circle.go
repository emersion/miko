package hitbox

import (
	"math"

	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
)

const CircleContourLen = 4

type Circle struct {
	radius float64
}

func (hb *Circle) Contour(center *terrain.Position) []*terrain.Position {
	// TODO: change base angle according to direction
	contour := make([]*terrain.Position, CircleContourLen)

	for i := 0; i < CircleContourLen; i++ {
		angle := 2 * math.Pi * float64(i) / CircleContourLen
		x := center.X + hb.radius*math.Cos(angle)
		y := center.Y + hb.radius*math.Sin(angle)
		contour[i] = &terrain.Position{x, y}
	}

	return contour
}

func (hb *Circle) Footprint(center *terrain.Position, angle float64) []*terrain.Position {
	angle1 := angle + math.Pi/2
	angle2 := angle - math.Pi/2
	pt1 := &terrain.Position{center.X + math.Cos(angle1), center.Y + math.Sin(angle1)}
	pt2 := &terrain.Position{center.X + math.Cos(angle2), center.Y + math.Sin(angle2)}
	return []*terrain.Position{pt1, pt2}
}

func (hb *Circle) intersects(center *terrain.Position, other Hitbox, otherCenter *terrain.Position) (intersects bool, err error) {
	switch o := other.(type) {
	case *Circle:
		d := dist(center, otherCenter)
		intersects = (d <= hb.radius+o.radius)
	default:
		err = errUnsupportedHitbox
	}
	return
}

func NewCircle(radius float64) *Circle {
	return &Circle{radius}
}
