package entity_test

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"testing"
)

func TestService_Rewind(t *testing.T) {
	s := entity.NewService()
	f := s.Frontend()

	accepted := make(chan entity.Request)

	// Accept all requests
	go (func() {
		for {
			var req entity.Request
			select {
			case req = <-f.Creates:
			case req = <-f.Updates:
			case req = <-f.Deletes:
			}

			s.AcceptRequest(req)
			accepted <- req
		}
	})()

	t0 := message.AbsoluteTick(1)

	// Create a new entity
	ent := message.NewEntity()
	ent.Position.X = 10
	ent.Position.Y = 15
	f.Add(ent, t0)

	created := (<-accepted).(*entity.CreateRequest)
	id := created.Entity.Id
	if id == 0 {
		t.Fatal("Entity doesn't have an id after being added")
	}

	ent0 := s.Get(id)

	t1 := t0 + 5

	// Update the entity
	update := message.NewEntity()
	update.Id = id
	update.Position.X = 33
	update.Position.Y = 15
	diff := &message.EntityDiff{Position: true}
	f.Update(update, diff, t1)

	<-accepted

	// Rewind
	err := s.Rewind(s.GetTick() - t0)
	if err != nil {
		t.Fatal(err)
	}

	rewindEnt := s.Get(id)
	if rewindEnt.Position.X != ent0.Position.X {
		t.Fatal("Invalid rewind, expected position", ent0.Position.X, "but got", rewindEnt.Position.X)
	}
}
