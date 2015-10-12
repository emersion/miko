package message

type EntityId uint16

type Position struct {
	BX BlockCoord
	BY BlockCoord
	X PointCoord
	Y PointCoord
}

type Speed struct {
	Angle float32
	Norm float32
}

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

type EntityDiffPool struct {
	Created []*Entity
	Updated map[EntityId]*EntityDiff
	Deleted []EntityId
}

type EntityService interface {
	Get(id EntityId) *Entity
	Add(entity *Entity)
	Update(entity *Entity, diff *EntityDiff)
	Delete(id EntityId)
	IsDirty() bool
	Flush() *EntityDiffPool
}
