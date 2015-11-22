package message

// An entity unique identifier.
type EntityId uint16

// A position contains a block coordinates and a point coordinates with that block.
type Position struct {
	BX BlockCoord // Block X coordinate
	BY BlockCoord // Block Y coordinate
	X  PointCoord // X coordinate inside block
	Y  PointCoord // Y coordinate inside block
}

// Check if this position is the same as another one.
func (p *Position) Equals(other *Position) bool {
	return (p.BX == other.BX && p.BY == other.BY && p.X == other.X && p.Y == other.Y)
}

// A speed contains an angle and a norm.
type Speed struct {
	Angle float32
	Norm  float32
}

// Check if this speed is the same as another one.
func (s *Speed) Equals(other *Speed) bool {
	return (s.Angle == other.Angle && s.Norm == other.Norm)
}

// An entity type identifier.
type EntityType uint16

// A sprite index.
type Sprite uint16

// An entity attribute id.
type EntityAttrId uint16

// An entity.
// Can be a player, an object, a monster, and so on.
type Entity struct {
	Id         EntityId
	Type       EntityType
	Position   *Position
	Speed      *Speed
	Sprite     Sprite
	Attributes map[EntityAttrId]interface{}
}

// Check if this entity is the same as another one. Only attributes specified in
// the diff will be checked against.
func (e *Entity) EqualsWithDiff(other *Entity, diff *EntityDiff) bool {
	if diff.Type && e.Type != other.Type {
		return false
	}
	if diff.Position && !e.Position.Equals(other.Position) {
		return false
	}
	if diff.SpeedAngle && e.Speed.Angle != other.Speed.Angle {
		return false
	}
	if diff.SpeedNorm && e.Speed.Norm != other.Speed.Norm {
		return false
	}
	if diff.Sprite && e.Sprite != other.Sprite {
		return false
	}
	if diff.Attributes {
		// TODO
	}
	return true
}

// Initialize a new entity.
func NewEntity() *Entity {
	return &Entity{
		Position: &Position{},
		Speed:    &Speed{},
	}
}

// Contains a list of fields that have changed in an entity.
type EntityDiff struct {
	Position   bool
	SpeedAngle bool
	SpeedNorm  bool
	Type       bool
	Sprite     bool
	Attributes bool
}

// Get this diff's bitfield.
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
	if d.Type {
		bitfield |= 1 << 2
	}
	if d.Sprite {
		bitfield |= 1 << 1
	}
	if d.Attributes {
		bitfield |= 1 << 0
	}
	return bitfield
}

// Merge two diffs.
func (d *EntityDiff) Merge(other *EntityDiff) *EntityDiff {
	d.Position = d.Position || other.Position
	d.SpeedAngle = d.SpeedAngle || other.SpeedAngle
	d.SpeedNorm = d.SpeedNorm || other.SpeedNorm
	d.Type = d.Type || other.Type
	d.Sprite = d.Sprite || other.Sprite
	d.Attributes = d.Attributes || other.Attributes
	return d
}

// Create a new entity diff.
func NewEntityDiff() *EntityDiff {
	return &EntityDiff{}
}

// Parse a diff bitfield.
func NewEntityDiffFromBitfield(bitfield uint8) *EntityDiff {
	return &EntityDiff{
		bitfield&(1<<7) > 0,
		bitfield&(1<<6) > 0,
		bitfield&(1<<5) > 0,
		bitfield&(1<<2) > 0,
		bitfield&(1<<1) > 0,
		bitfield&(1<<0) > 0,
	}
}

// Create a new diff filled with one value.
func NewFilledEntityDiff(val bool) *EntityDiff {
	var bitfield uint8
	if val {
		bitfield = ^bitfield
	}
	return NewEntityDiffFromBitfield(bitfield)
}

// An entity diff pool.
// Contains three lists for created, updated and deleted entities.
type EntityDiffPool struct {
	Created []*Entity
	Updated map[*Entity]*EntityDiff
	Deleted []EntityId
}

// Check if this diff pool is empty.
func (dp *EntityDiffPool) IsEmpty() bool {
	return len(dp.Created) == 0 && len(dp.Updated) == 0 && len(dp.Deleted) == 0
}

// Create a new diff pool.
func NewEntityDiffPool() *EntityDiffPool {
	return &EntityDiffPool{Updated: map[*Entity]*EntityDiff{}}
}

// An entity service.
type EntityService interface {
	// Get a specific entity.
	Get(id EntityId) *Entity

	// Add a new entity.
	Add(entity *Entity, t AbsoluteTick)

	// Update an entity.
	Update(entity *Entity, diff *EntityDiff, t AbsoluteTick)

	// Delete an entity.
	Delete(id EntityId, t AbsoluteTick)

	// Check if some entities have been added, updated or deleted.
	IsDirty() bool

	// Flush the current diff pool. Return the current one and replace it by a new
	// empty one.
	Flush() *EntityDiffPool
}
