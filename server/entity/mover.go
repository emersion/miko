package entity

import (
	"git.emersion.fr/saucisse-royale/miko/server/message"
	"time"
)

// A service that moves entities
type Mover struct {
	terrain     message.Terrain
	lastUpdates map[message.EntityId]int64
	positions   map[message.EntityId]*Position
}

// Compute an entity's new position
// Returns true if the position has changed, false otherwise
func (m *Mover) UpdateEntity(entity *message.Entity) *message.EntityDiff {
	last := m.lastUpdates[entity.Id]
	now := time.Now().UnixNano()
	m.lastUpdates[entity.Id] = now
	dt := float64(now-last) / 1000 / 1000 // Convert to seconds

	speed := NewSpeedFromMessage(entity.Speed) // TODO
	var pos *Position
	var ok bool
	if pos, ok = m.positions[entity.Id]; !ok {
		pos = NewPositionFromMessage(entity.Position)
	}

	nextPos := speed.GetNextPosition(pos, dt)
	if nextPos == nil {
		return nil
	}

	// Check terrain
	canMove := true
	pts := GetRouteBetween(pos, nextPos)
	for _, pt := range pts {
		t := m.terrain.GetPointAt(pt[0], pt[1])

		if t != message.PointType(0) {
			canMove = false
		}
	}

	if !canMove {
		return nil
	}

	m.positions[entity.Id] = nextPos
	entity.Position = nextPos.ToMessage()

	return &message.EntityDiff{Position: true}
}

func NewMover(trn message.Terrain) *Mover {
	return &Mover{
		terrain:     trn,
		lastUpdates: map[message.EntityId]int64{},
		positions:   map[message.EntityId]*Position{},
	}
}
