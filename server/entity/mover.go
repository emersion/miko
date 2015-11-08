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
	lastUpdates map[message.EntityId]uint64
	positions   map[message.EntityId]*Position
}

// Compute an entity's new position
// Returns an EntityDiff if the entity has changed, nil otherwise.
func (m *Mover) UpdateEntity(entity *message.Entity) *message.EntityDiff {
	now := m.clock.GetAbsoluteTick()

	var last uint64
	var ok bool
	if last, ok = m.lastUpdates[entity.Id]; !ok {
		last = now
	}

	m.lastUpdates[entity.Id] = now
	dt := time.Duration(now-last) * clock.TickDuration // Convert to seconds
	if dt == 0 {
		return nil
	}

	speed := NewSpeedFromMessage(entity.Speed)
	var pos *Position
	if pos, ok = m.positions[entity.Id]; !ok {
		pos = NewPositionFromMessage(entity.Position)
	}

	nextPos := speed.GetNextPosition(pos, dt)
	if nextPos == nil {
		return nil
	}

	// Check terrain
	pts := GetRouteBetween(pos, nextPos)
	var lastPt [2]int
	for _, pt := range pts {
		t, err := m.terrain.GetPointAt(pt[0], pt[1])
		if err != nil {
			return nil // TODO: trigger a more severe error
		}

		if t != message.PointType(0) {
			nextPos = &Position{float64(lastPt[0]), float64(lastPt[1])}
			break
		}

		lastPt = pt
	}

	m.positions[entity.Id] = nextPos
	entity.Position = nextPos.ToMessage()

	return &message.EntityDiff{Position: true}
}

func NewMover(trn message.Terrain, clk message.ClockService) *Mover {
	return &Mover{
		terrain:     trn,
		clock:       clk,
		lastUpdates: map[message.EntityId]uint64{},
		positions:   map[message.EntityId]*Position{},
	}
}
