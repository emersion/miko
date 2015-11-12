// The game engine.
package engine

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/auth"
	"git.emersion.fr/saucisse-royale/miko.git/server/clock"
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/terrain"
	"time"
)

type Engine struct {
	entity *entity.Service
	ctx    *message.Context
}

func (e *Engine) Start() {
	entityFrontend := e.entity.Frontend()
	mover := entity.NewMover(e.ctx.Terrain, e.ctx.Clock)

	for {
		start := time.Now().UnixNano()
		e.ctx.Clock.Tick()

		// Process requests from clients
		// TODO: security checks
		// TODO: initiate lag compensation
		for {
			select {
			case req := <-entityFrontend.Creates:
				e.entity.Add(req.Entity, req.Tick)
			case req := <-entityFrontend.Updates:
				e.entity.Update(req.Entity, req.Diff, req.Tick)
			case req := <-entityFrontend.Deletes:
				e.entity.Delete(req.EntityId, req.Tick)
			default:
				break
			}
		}

		// Animate entities
		for _, entity := range e.ctx.Entity.List() {
			diff := mover.UpdateEntity(entity)
			if diff != nil {
				e.ctx.Entity.Update(entity, diff, e.ctx.Clock.GetAbsoluteTick())
			}
		}

		end := time.Now().UnixNano()
		time.Sleep(clock.TickDuration - time.Nanosecond*time.Duration(end-start))
	}
}

func (e *Engine) Context() *message.Context {
	return e.ctx
}

func New() *Engine {
	// Create a new context
	ctx := message.NewServerContext()

	// Create the engine
	e := &Engine{
		ctx:    ctx,
		entity: entity.NewService(),
	}

	// Populate context
	ctx.Auth = auth.NewService()
	ctx.Entity = e.entity.Frontend()
	ctx.Terrain = terrain.New()
	ctx.Clock = clock.NewService()

	return e
}
