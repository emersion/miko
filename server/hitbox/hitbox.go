// Provides functions to compute intersections between hitboxes.
package hitbox

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
	"math"
)

func round(f float64) int {
	return int(math.Floor(f + .5))
}

func dist(a, b *entity.Position) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}

type Hitbox interface {
	Contour() []*entity.Position
}

type hitbox interface {
	intersects(other Hitbox) (bool, error)
}

func Intersects(a, b Hitbox) bool {
	if ahb, ok := a.(hitbox); ok {
		intersects, err := ahb.intersects(b)
		if err == nil {
			return intersects
		}
	}
	if bhb, ok := b.(hitbox); ok {
		intersects, err := bhb.intersects(a)
		if err == nil {
			return intersects
		}
	}

	// TODO
	return false
}
