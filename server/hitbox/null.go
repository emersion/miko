package hitbox

type Null struct{}

func (hb *Null) Contour(other Hitbox) []*entity.Position {
	return []*entity.Position{}
}
