// Package hitbox provides functions to compute intersections between hitboxes.
package hitbox

import (
	"errors"
	"math"

	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
)

var errUnsupportedHitbox = errors.New("Unsupported hitbox")

func round(f float64) int {
	return int(math.Floor(f + .5))
}

func dist(a, b *terrain.Position) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}

// Hitbox is an object's shape.
type Hitbox interface {
	// Contour returns an approximation of the contour of this hitbox, if centered
	// at the provided position.
	Contour(center *terrain.Position) []*terrain.Position

	// Footprint returns 0 to 2 points, representing the hitbox's shape if
	// centered at the provided position and going in the provided direction.
	Footprint(center *terrain.Position, angle float64) []*terrain.Position
}

type hitbox interface {
	intersects(other Hitbox) (bool, error)
}

// Intersects checks if two hitboxes overlap.
func Intersects(a, b Hitbox) bool {
	if a, ok := a.(hitbox); ok {
		if intersects, err := a.intersects(b); err == nil {
			return intersects
		}
	}
	if b, ok := b.(hitbox); ok {
		if intersects, err := b.intersects(a); err == nil {
			return intersects
		}
	}

	// TODO: use Contour
	return false
}
