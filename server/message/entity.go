package message

import "math/cmplx"

type EntityId uint16

// A position contains a block coordinates and a point coordinates with that block.
type Position struct {
	BX BlockCoord
	BY BlockCoord
	X PointCoord
	Y PointCoord
}

// Get absolute coordinates of this position.
func (p *Position) AbsoluteCoords() (x, y int) {
	x = int(p.BX) * BLOCK_LEN + int(p.X)
	y = int(p.BY) * BLOCK_LEN + int(p.Y)
	return
}

// Create a new position from absolute coordinates.
// A @delthas powered function name.
func NewPositionFromAbsoluteCoords(x, y int) *Position {
	pos := &Position{}

	pos.X = PointCoord(x % BLOCK_LEN)
	pos.Y = PointCoord(y % BLOCK_LEN)
	pos.BX = BlockCoord((x - int(pos.X)) / BLOCK_LEN)
	pos.BY = BlockCoord((y - int(pos.Y)) / BLOCK_LEN)

	return pos
}

// A speed contains an angle and a norm.
type Speed struct {
	Angle float32
	Norm float32
}

// Get the position reached by an object at t+dt if it has this speed during dt.
func (s *Speed) GetNextPosition(current *Position, dt float64) *Position {
	if s.Norm == 0 {
		return nil
	}

	speed := cmplx.Rect(float64(s.Norm), float64(s.Angle))
	x, y := current.AbsoluteCoords()
	pos := complex(float64(x), float64(y))

	pos += speed * complex(dt, 0)

	return NewPositionFromAbsoluteCoords(int(real(pos)), int(imag(pos)))
}

// An entity
// Can be a player, an object, a monster, and so on.
type Entity struct {
	Id EntityId
	Position *Position
	Speed *Speed
}

func NewEntity() *Entity {
	return &Entity{
		Position: &Position{},
		Speed: &Speed{},
	}
}

// Contains a list of fields that have changed in an entity.
type EntityDiff struct {
	Position bool
	SpeedAngle bool
	SpeedNorm bool
}

func (d *EntityDiff) GetBitfield() uint8 {
	var bitfield uint8
	if d.Position {
		bitfield |= 1 << 7
	}
	if d.SpeedAngle {
		bitfield |= 1 << 6
	}
	if d.SpeedNorm {
		bitfield |= 1 << 5
	}
	return bitfield
}

func (d *EntityDiff) Merge(other *EntityDiff) {
	d.Position = d.Position || other.Position
	d.SpeedAngle = d.SpeedAngle || other.SpeedAngle
	d.SpeedNorm = d.SpeedNorm || other.SpeedNorm
}

func (d *EntityDiff) Apply(src *Entity, dst *Entity) {
	if d.Position {
		dst.Position = src.Position
	}
	if d.SpeedNorm {
		dst.Speed.Norm = src.Speed.Norm
	}
	if d.SpeedAngle {
		dst.Speed.Angle = src.Speed.Angle
	}
}

func NewEntityDiffFromBitfield(bitfield uint8) *EntityDiff {
	return &EntityDiff{
		bitfield & (1 << 7) > 0,
		bitfield & (1 << 6) > 0,
		bitfield & (1 << 5) > 0,
	}
}

// An entity diff pool
// Contains three lists for created, updated and deleted entities.
type EntityDiffPool struct {
	Created []*Entity
	Updated map[*Entity]*EntityDiff
	Deleted []EntityId
}

func NewEntityDiffPool() *EntityDiffPool {
	return &EntityDiffPool{Updated: map[*Entity]*EntityDiff{}}
}

type EntityMover interface {
	UpdateEntity(entity *Entity) *Position
}

// An entity service
type EntityService interface {
	List() []*Entity
	Get(id EntityId) *Entity
	Add(entity *Entity)
	Update(entity *Entity, diff *EntityDiff)
	Delete(id EntityId)
	IsDirty() bool
	Flush() *EntityDiffPool
	Mover() EntityMover
}
