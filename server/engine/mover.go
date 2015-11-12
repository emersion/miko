package engine

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/clock"
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"time"
)

func CheckRoute(route entity.Route, ent *message.Entity, trn message.Terrain) *entity.Position {
	var last entity.RouteStep
	for _, step := range route {
		t, err := trn.GetPointAt(step[0], step[1])

		if err != nil || t != message.PointType(0) {
			return &entity.Position{float64(last[0]), float64(last[1])}
		}

		last = step
	}

	return nil
}

// A service that moves entities.
type Mover struct {
	engine      *Engine
	lastUpdates map[message.EntityId]message.AbsoluteTick
	positions   map[message.EntityId]*entity.Position // TODO: positions may be outdated, refresh from entity.Service!
}

// Compute an entity's new position.
// Returns an EntityDiff if the entity has changed, nil otherwise.
func (m *Mover) UpdateEntity(ent *message.Entity) *message.EntityDiff {
	now := m.engine.ctx.Clock.GetAbsoluteTick()

	var last message.AbsoluteTick
	var ok bool
	if last, ok = m.lastUpdates[ent.Id]; !ok {
		last = now
	}

	m.lastUpdates[ent.Id] = now
	dt := time.Duration(now-last) * clock.TickDuration // Convert to seconds
	if dt == 0 {
		return nil
	}

	speed := entity.NewSpeedFromMessage(ent.Speed)
	var pos *entity.Position
	if pos, ok = m.positions[ent.Id]; !ok {
		pos = entity.NewPositionFromMessage(ent.Position)
	}

	nextPos := speed.GetNextPosition(pos, dt)
	if nextPos == nil {
		return nil
	}

	// Check terrain
	route := entity.GetRouteBetween(pos, nextPos)

	stoppedAt := CheckRoute(route, ent, m.engine.ctx.Terrain)
	if stoppedAt != nil {
		// The entity could has been stopped while moving
		nextPos = stoppedAt
	}

	m.positions[ent.Id] = nextPos
	ent.Position = nextPos.ToMessage()

	return &message.EntityDiff{Position: true}
}

func NewMover(engine *Engine) *Mover {
	return &Mover{
		engine:      engine,
		lastUpdates: map[message.EntityId]message.AbsoluteTick{},
		positions:   map[message.EntityId]*entity.Position{},
	}
}
