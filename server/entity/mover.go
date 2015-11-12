package entity

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/clock"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"time"
)

func CheckRoute(route Route, entity *message.Entity, trn message.Terrain) *Position {
	var last RouteStep
	for _, step := range route {
		t, err := trn.GetPointAt(step[0], step[1])

		if err != nil || t != message.PointType(0) {
			return &Position{float64(last[0]), float64(last[1])}
		}

		last = step
	}

	return nil
}

// A service that moves entities
type Mover struct {
	terrain     message.Terrain
	clock       message.ClockService
	lastUpdates map[message.EntityId]message.AbsoluteTick
	positions   map[message.EntityId]*Position
}

// Compute an entity's new position
// Returns an EntityDiff if the entity has changed, nil otherwise.
func (m *Mover) UpdateEntity(entity *message.Entity) *message.EntityDiff {
	now := m.clock.GetAbsoluteTick()

	var last message.AbsoluteTick
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
	route := GetRouteBetween(pos, nextPos)

	stoppedAt := CheckRoute(route, entity, m.terrain)
	if stoppedAt != nil {
		// The entity could has been stopped while moving
		nextPos = stoppedAt
	}

	m.positions[entity.Id] = nextPos
	entity.Position = nextPos.ToMessage()

	return &message.EntityDiff{Position: true}
}

func NewMover(trn message.Terrain, clk message.ClockService) *Mover {
	return &Mover{
		terrain:     trn,
		clock:       clk,
		lastUpdates: map[message.EntityId]message.AbsoluteTick{},
		positions:   map[message.EntityId]*Position{},
	}
}
