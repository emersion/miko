package engine

import (
	"time"

	"git.emersion.fr/saucisse-royale/miko.git/server/clock"
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
	"git.emersion.fr/saucisse-royale/miko.git/server/game"
	"git.emersion.fr/saucisse-royale/miko.git/server/hitbox"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
)

func checkRoute(route terrain.Route, ent *entity.Entity, trn message.Terrain, entities []*message.Entity) (*terrain.Position, interface{}) {
	hb := game.GetHitbox(ent.Sprite)

	var last terrain.RouteStep
	for _, step := range route {
		// Check if terrain collides
		t, err := trn.GetPointAt(step[0], step[1])
		if err != nil || t != game.PointTypeEmpty {
			return &terrain.Position{float64(last[0]), float64(last[1])}, t
		}

		// Check if entities collide
		// TODO: optimize this
		for _, e := range entities {
			if e.Id == ent.Id || !game.Collides(ent.Type, e.Type) {
				continue
			}

			if hitbox.Intersects(hb, game.GetHitbox(e.Sprite)) {
				return &terrain.Position{float64(last[0]), float64(last[1])}, e
			}
		}

		last = step
	}

	return nil, nil
}

// A service that moves entities.
type Mover struct {
	engine      *Engine
	lastUpdates map[message.EntityId]message.AbsoluteTick
}

// Compute an entity's new position.
// Returns an EntityDiff if the entity has changed, nil otherwise.
func (m *Mover) UpdateEntity(ent *entity.Entity, now message.AbsoluteTick) (req *entity.UpdateRequest, collidesWith interface{}) {
	// TODO: remove Mover.lastUpdates?
	var last message.AbsoluteTick
	var ok bool
	if last, ok = m.lastUpdates[ent.Id]; !ok {
		last = now
	}

	m.lastUpdates[ent.Id] = now
	dt := time.Duration(now-last) * clock.TickDuration // Convert to seconds
	if dt == 0 {
		return
	}
	if dt < 0 {
		return // TODO: figure out if it's the right thing to do here
	}

	speed := ent.Speed
	pos := ent.Position

	nextPos := speed.GetNextPosition(pos, dt)
	if nextPos == nil {
		return
	}

	// Check terrain
	// TODO: use brensenham algorithm
	// See http://tech-algorithm.com/articles/drawing-line-using-bresenham-algorithm/
	route := terrain.GetRouteBetween(pos, nextPos)

	stoppedAt, collidesWith := checkRoute(route, ent, m.engine.ctx.Terrain, m.engine.ctx.Entity.List())
	if stoppedAt != nil {
		// The entity has been stopped while moving
		nextPos = stoppedAt
	}

	newEnt := entity.New()
	newEnt.Id = ent.Id
	newEnt.Position = nextPos
	diff := &message.EntityDiff{Position: true}

	req = entity.NewUpdateRequest(now, newEnt, diff)
	return
}

func NewMover(engine *Engine) *Mover {
	return &Mover{
		engine:      engine,
		lastUpdates: map[message.EntityId]message.AbsoluteTick{},
	}
}
