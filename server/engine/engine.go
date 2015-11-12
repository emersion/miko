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
	mover := NewMover(e)

	for {
		start := time.Now().UnixNano()
		e.ctx.Clock.Tick()

		// Process requests from clients
		// TODO: security checks
		minTick := e.ctx.Clock.GetAbsoluteTick() - message.MaxRewind
		minAcceptedTick := e.ctx.Clock.GetAbsoluteTick()
		accepted := []interface{}{}
		for {
			select {
			case req := <-entityFrontend.Creates:
				if req.Tick < minTick {
					continue
				}
				if req.Tick < minAcceptedTick {
					minAcceptedTick = req.Tick
				}
				accepted = append(accepted, req)
			case req := <-entityFrontend.Updates:
				if req.Tick < minTick {
					continue
				}
				if req.Tick < minAcceptedTick {
					minAcceptedTick = req.Tick
				}
				accepted = append(accepted, req)
			case req := <-entityFrontend.Deletes:
				if req.Tick < minTick {
					continue
				}
				if req.Tick < minAcceptedTick {
					minAcceptedTick = req.Tick
				}
				accepted = append(accepted, req)
			default:
				break
			}
		}

		if minAcceptedTick < e.ctx.Clock.GetAbsoluteTick() {
			// TODO: initiate lag compensation
		}

		// Pass all accepted requests to entity service
		for _, item := range accepted {
			switch req := item.(type) {
			case entity.CreateRequest:
				e.entity.Add(req.Entity, req.Tick)
			case entity.UpdateRequest:
				e.entity.Update(req.Entity, req.Diff, req.Tick)
			case entity.DeleteRequest:
				e.entity.Delete(req.EntityId, req.Tick)
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

	ctx.Config = &message.Config{
		MaxRollbackTicks: uint16(message.MaxRewind),
	}

	return e
}
