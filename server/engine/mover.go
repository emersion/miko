package engine

import (
	"time"

	"git.emersion.fr/saucisse-royale/miko.git/server/clock"
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
)

func CheckRoute(route terrain.Route, ent *entity.Entity, trn message.Terrain) *terrain.Position {
	var last terrain.RouteStep
	for _, step := range route {
		t, err := trn.GetPointAt(step[0], step[1])

		if err != nil || t != message.PointType(0) {
			return &terrain.Position{float64(last[0]), float64(last[1])}
		}

		last = step
	}

	return nil
}

// A service that moves entities.
type Mover struct {
	engine      *Engine
	lastUpdates map[message.EntityId]message.AbsoluteTick
}

// Compute an entity's new position.
// Returns an EntityDiff if the entity has changed, nil otherwise.
func (m *Mover) UpdateEntity(ent *entity.Entity, now message.AbsoluteTick) *entity.UpdateRequest {
	// TODO: remove Mover.lastUpdates?
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
	if dt < 0 {
		return nil // TODO: figure out if it's the right thing to do here
	}

	speed := ent.Speed
	pos := ent.Position

	nextPos := speed.GetNextPosition(pos, dt)
	if nextPos == nil {
		return nil
	}

	// Check terrain
	route := terrain.GetRouteBetween(pos, nextPos)

	stoppedAt := CheckRoute(route, ent, m.engine.ctx.Terrain)
	if stoppedAt != nil {
		// The entity could has been stopped while moving
		nextPos = stoppedAt
	}

	newEnt := entity.New()
	newEnt.Id = ent.Id
	newEnt.Position = nextPos
	diff := &message.EntityDiff{Position: true}

	return entity.NewUpdateRequest(now, newEnt, diff)
}

func NewMover(engine *Engine) *Mover {
	return &Mover{
		engine:      engine,
		lastUpdates: map[message.EntityId]message.AbsoluteTick{},
	}
}
