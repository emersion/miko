package engine_test

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/clock"
	"git.emersion.fr/saucisse-royale/miko.git/server/engine"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"testing"
	"time"
)

func TestEngine(t *testing.T) {
	e := engine.New()
	ctx := e.Context()

	go e.Start()

	time.Sleep(5 * clock.TickDuration)

	// Create a new entity
	ent := message.NewEntity()
	ent.Id = 1
	ent.Position.X = 10
	ent.Position.Y = 15
	ctx.Entity.Add(ent, ctx.Clock.GetAbsoluteTick())

	time.Sleep(5 * clock.TickDuration)

	// Update the entity
	update := message.NewEntity()
	update.Id = ent.Id
	update.Speed.Norm = 10
	diff := &message.EntityDiff{SpeedNorm: true}
	ctx.Entity.Update(update, diff, ctx.Clock.GetAbsoluteTick())

	time.Sleep(2 * clock.TickDuration)

	// Update the entity a second time, in the past
	update = message.NewEntity()
	update.Id = ent.Id
	update.Speed.Norm = 20
	ctx.Entity.Update(update, diff, ctx.Clock.GetAbsoluteTick()-10)

	e.Stop()
}
