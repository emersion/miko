package engine_test

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/clock"
	"git.emersion.fr/saucisse-royale/miko.git/server/engine"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/server"
	"testing"
	"time"
)

func TestEngine(t *testing.T) {
	srv := server.New("")
	e := engine.New(srv)
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

	time.Sleep(clock.TickDuration)

	pool := ctx.Entity.Flush()
	if len(pool.Created) != 1 {
		t.Error("Created an entity, but it isn't in the diff pool")
	}
	if len(pool.Updated) != 0 {
		t.Error("Updated an entity, but it shouldn't be in the diff pool since it just has been created")
	}
	if len(pool.Deleted) != 0 {
		t.Error("No entity deleted, but there is one in the diff pool")
	}

	if pool.Created[0].Id != ent.Id {
		t.Error("The entity in the diff pool has a wrong id")
	}
	if pool.Created[0].Speed.Norm != 10 {
		t.Error("The entity in the diff pool has a wrong speed norm")
	}

	// Update the entity another time
	update = message.NewEntity()
	update.Id = ent.Id
	update.Speed.Norm = 15
	ctx.Entity.Update(update, diff, ctx.Clock.GetAbsoluteTick())

	time.Sleep(clock.TickDuration)

	pool = ctx.Entity.Flush()
	if len(pool.Created) != 0 {
		t.Error("No entity created, but there is one in the diff pool")
	}
	if len(pool.Updated) != 1 {
		t.Error("Updated an entity, but it isn't in the diff pool")
	}
	if len(pool.Deleted) != 0 {
		t.Error("No entity deleted, but there is one in the diff pool")
	}

	for updated, diff := range pool.Updated {
		if updated.Id != ent.Id {
			t.Error("The entity in the diff pool has a wrong id")
		}
		if updated.Speed.Norm != 15 {
			t.Error("The entity in the diff pool has a wrong speed norm")
		}
		if !diff.SpeedNorm {
			t.Error("The diff doesn't mark the speed norm as outdated, but it has just been updated")
		}
	}

	e.Stop()
}
