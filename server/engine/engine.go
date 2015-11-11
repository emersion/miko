// The game engine.
package engine

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/clock"
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"time"
)

type Engine struct {
	ctx *message.Context
}

func (e *Engine) Start() {
	mover := entity.NewMover(e.ctx.Terrain, e.ctx.Clock)

	for {
		start := time.Now().UnixNano()
		e.ctx.Clock.Tick()

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

func New(ctx *message.Context) *Engine {
	return &Engine{ctx}
}
