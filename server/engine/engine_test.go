package engine_test

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/clock"
	"git.emersion.fr/saucisse-royale/miko.git/server/engine"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"testing"
	"time"
)

func TestEngine(t *testing.T) {
	e := engine.New(nil)
	ctx := e.Context()

	go e.Start()

	time.Sleep(5 * clock.TickDuration)

	// Create a new entity
	createdAt := ctx.Clock.GetAbsoluteTick()
	t.Log("Adding new entity at tick", createdAt)
	ent := message.NewEntity()
	ent.Id = 1
	ent.Position.X = 10
	ent.Position.Y = 15
	req := ctx.Entity.Add(ent, createdAt)
	if err := req.Wait(); err != nil {
		t.Fatal("Cannot create entity:", err)
	}

	time.Sleep(5 * clock.TickDuration)

	// Update the entity
	updatedAt := ctx.Clock.GetAbsoluteTick()
	t.Log("Updating the entity at tick", updatedAt)
	update := message.NewEntity()
	update.Id = ent.Id
	update.Speed.Norm = 10
	diff := &message.EntityDiff{SpeedNorm: true}
	req = ctx.Entity.Update(update, diff, updatedAt)
	if err := req.Wait(); err != nil {
		t.Fatal("Cannot update entity:", err)
	}

	time.Sleep(2 * clock.TickDuration)

	// Update the entity a second time, in the past
	reupdatedAt := ctx.Clock.GetAbsoluteTick() - 5
	t.Log("Updating the entity a second time at tick", reupdatedAt)
	update = message.NewEntity()
	update.Id = ent.Id
	update.Speed.Norm = 20
	req = ctx.Entity.Update(update, diff, reupdatedAt)
	if err := req.Wait(); err != nil {
		t.Fatal("Cannot update entity:", err)
	}

	time.Sleep(clock.TickDuration)

	pool := ctx.Entity.Flush()
	if len(pool.Created) != 1 {
		t.Fatal("Created an entity, but it isn't in the diff pool")
	}
	if len(pool.Updated) != 0 {
		t.Fatal("Updated an entity, but it shouldn't be in the diff pool since it just has been created")
	}
	if len(pool.Deleted) != 0 {
		t.Fatal("No entity deleted, but there is one in the diff pool")
	}

	if pool.Created[0].Id != ent.Id {
		t.Fatal("The entity in the diff pool has a wrong id")
	}
	if pool.Created[0].Speed.Norm != 10 {
		t.Fatal("The entity in the diff pool has a wrong speed norm, got", pool.Created[0].Speed.Norm, "instead of", 10)
	}

	// Update the entity another time
	update = message.NewEntity()
	update.Id = ent.Id
	update.Speed.Norm = 15
	ctx.Entity.Update(update, diff, ctx.Clock.GetAbsoluteTick())

	time.Sleep(clock.TickDuration)

	pool = ctx.Entity.Flush()
	if len(pool.Created) != 0 {
		t.Fatal("No entity created, but there is one in the diff pool")
	}
	if len(pool.Updated) != 1 {
		t.Fatal("Updated an entity, but it isn't in the diff pool")
	}
	if len(pool.Deleted) != 0 {
		t.Fatal("No entity deleted, but there is one in the diff pool")
	}

	for updated, diff := range pool.Updated {
		if updated.Id != ent.Id {
			t.Fatal("The entity in the diff pool has a wrong id")
		}
		if updated.Speed.Norm != 15 {
			t.Fatal("The entity in the diff pool has a wrong speed norm")
		}
		if !diff.SpeedNorm {
			t.Fatal("The diff doesn't mark the speed norm as outdated, but it has just been updated")
		}
	}

	e.Stop()
}
