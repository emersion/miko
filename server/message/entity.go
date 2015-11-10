package message

// An entity unique identifier.
type EntityId uint16

// A position contains a block coordinates and a point coordinates with that block.
type Position struct {
	BX BlockCoord
	BY BlockCoord
	X  PointCoord
	Y  PointCoord
}

// A speed contains an angle and a norm.
type Speed struct {
	Angle float32
	Norm  float32
}

// A sprite index.
type Sprite uint16

// An entity.
// Can be a player, an object, a monster, and so on.
type Entity struct {
	Id       EntityId
	Position *Position
	Speed    *Speed
	Sprite   Sprite
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
	Sprite     bool
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
	if d.Sprite {
		bitfield |= 1 << 1
	}
	return bitfield
}

// Merge two diffs.
func (d *EntityDiff) Merge(other *EntityDiff) {
	d.Position = d.Position || other.Position
	d.SpeedAngle = d.SpeedAngle || other.SpeedAngle
	d.SpeedNorm = d.SpeedNorm || other.SpeedNorm
	d.Sprite = d.Sprite || other.Sprite
}

// Apply a diff from a source entity to a destination entity.
// Changed properties will be copied from the source and overwrite the
// destination's ones.
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
	if d.Sprite {
		dst.Sprite = src.Sprite
	}
}

// Parse a diff bitfield.
func NewEntityDiffFromBitfield(bitfield uint8) *EntityDiff {
	return &EntityDiff{
		bitfield&(1<<7) > 0,
		bitfield&(1<<6) > 0,
		bitfield&(1<<5) > 0,
		bitfield&(1<<1) > 0,
	}
}

// An entity diff pool.
// Contains three lists for created, updated and deleted entities.
type EntityDiffPool struct {
	Created []*Entity
	Updated map[*Entity]*EntityDiff
	Deleted []EntityId
}

// Create a new diff pool.
func NewEntityDiffPool() *EntityDiffPool {
	return &EntityDiffPool{Updated: map[*Entity]*EntityDiff{}}
}

// An entity service.
type EntityService interface {
	// Get a list of all entities.
	List() []*Entity

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

	// Run a loop that will animate entities, ie. increment the clock and calculate
	// their new positions.
	Animate(trn Terrain, clk ClockService)
}
