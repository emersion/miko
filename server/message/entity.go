package message

type EntityId uint16

// A position contains a block coordinates and a point coordinates with that block.
type Position struct {
	BX BlockCoord
	BY BlockCoord
	X PointCoord
	Y PointCoord
}

// A speed contains an angle and a norm.
type Speed struct {
	Angle float32
	Norm float32
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

// An entity service
type EntityService interface {
	List() []*Entity
	Get(id EntityId) *Entity
	Add(entity *Entity)
	Update(entity *Entity, diff *EntityDiff)
	Delete(id EntityId)
	IsDirty() bool
	Flush() *EntityDiffPool
	Animate(trn Terrain)
}
