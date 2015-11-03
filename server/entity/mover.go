package entity

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/clock"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"time"
)

// A service that moves entities
type Mover struct {
	terrain     message.Terrain
	clock       message.ClockService
	lastUpdates map[message.EntityId]int64
	positions   map[message.EntityId]*Position
}

// Compute an entity's new position
// Returns true if the position has changed, false otherwise
func (m *Mover) UpdateEntity(entity *message.Entity) *message.EntityDiff {
	last := m.lastUpdates[entity.Id]
	now := m.clock.GetTicks()
	m.lastUpdates[entity.Id] = now
	dt := time.Duration(now-last) * clock.TickDuration // Convert to seconds

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

func NewMover(trn message.Terrain, clk message.ClockService) *Mover {
	return &Mover{
		terrain:     trn,
		clock:       clk,
		lastUpdates: map[message.EntityId]int64{},
		positions:   map[message.EntityId]*Position{},
	}
}
