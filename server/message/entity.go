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

type EntityUpdate struct {
	EntityId EntityId
	Diff *EntityDiff
}

func (u *EntityUpdate) merge(other *EntityUpdate) {
	if u.EntityId != other.EntityId {
		return
	}
	u.Diff.Merge(other.Diff)
}

type EntityDiffPool struct {
	Created []*Entity
	Updated []*EntityUpdate
	Deleted []EntityId
}

func (d *EntityDiffPool) MergeUpdated() {
	var updatesIds map[EntityId]int
	for i, update := range d.Updated {
		if j, ok := updatesIds[update.EntityId]; ok {
			// Another update for this entity, merge
			d.Updated[j].merge(update)

			// Remove this update from the list
			d.Updated[i] = d.Updated[len(d.Updated)-1]
			d.Updated = d.Updated[:len(d.Updated)-1]
		} else {
			// First update of this entity
			updatesIds[update.EntityId] = i
		}
	}
}

type EntityService interface {
	Get(id EntityId) *Entity
	Add(entity *Entity)
	Update(entity *Entity, diff *EntityDiff)
	Delete(id EntityId)
	IsDirty() bool
	Flush() *EntityDiffPool
}
