// Provides functions to compute intersections between hitboxes.
package hitbox

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
)

type Hitbox interface {
	Contour() []*entity.Position
}
