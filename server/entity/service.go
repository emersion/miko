// Provides functions to keep track of entities.
package entity

import (
	"errors"
	"fmt"
	"git.emersion.fr/saucisse-royale/miko.git/server/delta"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

// A delta (not to be confused with @delthas) stores all data about an entity
// change: not only its diff if it has been updated, but also the time when the
// change has been made and the entity itself.
type Delta struct {
	tick      message.AbsoluteTick
	requested bool
	EntityId  message.EntityId

	Diff *message.EntityDiff
	From *message.Entity
	To   *message.Entity
}

func (d *Delta) GetTick() message.AbsoluteTick {
	return d.tick
}

func (d *Delta) Requested() bool {
	return d.requested
}

func (d *Delta) Request() Request {
	req := newRequest(d.tick)
	req.requested = d.requested
	req.accepted = true

	if d.From != nil && d.To != nil { // Update
		if d.Diff == nil {
			panic("Cannot build update request from delta: no diff available")
		}

		return &UpdateRequest{req, d.To, d.Diff}
	} else if d.To != nil { // Create
		return &CreateRequest{req, d.To}
	} else if d.From != nil { // Delete
		return &DeleteRequest{req, d.From.Id}
	}

	panic("Cannot build request from delta: from and to are empty")
}

func copyFromDiff(src *message.Entity, diff *message.EntityDiff) *message.Entity {
	dst := &message.Entity{}
	diff.Apply(src, dst)
	return dst
}

// The entity service.
// It manages all entities by maintaining a list of them and a diff pool. The
// diff pool keeps track of created, updated and deleted entities to send
// appropriate messages to clients.
type Service struct {
	entities map[message.EntityId]*message.Entity
	deltas   *delta.List
	tick     message.AbsoluteTick
	frontend *Frontend
}

func (s *Service) GetTick() message.AbsoluteTick {
	return s.tick
}

func (s *Service) List() map[message.EntityId]*message.Entity {
	return s.entities
}

func (s *Service) Get(id message.EntityId) *message.Entity {
	if entity, ok := s.entities[id]; ok {
		return entity
	}
	return nil
}

func (s *Service) AcceptRequest(req Request) error {
	switch r := req.(type) {
	case *CreateRequest:
		return s.acceptCreate(r)
	case *UpdateRequest:
		return s.acceptUpdate(r)
	case *DeleteRequest:
		return s.acceptDelete(r)
	default:
		panic("Cannot accept request: not a request")
	}
}

func (s *Service) acceptCreate(req *CreateRequest) error {
	entity := req.Entity

	if int(entity.Id) == 0 {
		nextId := len(s.entities)
		if nextId == 0 {
			nextId = 1
		}
		entity.Id = message.EntityId(nextId)
	} else if s.Get(entity.Id) != nil {
		return errors.New("Cannot create entity: entity already exists")
	}

	s.entities[entity.Id] = entity

	if !req.accepted {
		s.deltas.Insert(&Delta{
			tick:      req.tick,
			requested: req.requested,
			EntityId:  entity.Id,
			From:      nil,
			To:        entity,
		})
	}

	s.tick = req.tick
	req.accepted = true

	return nil
}

func (s *Service) acceptUpdate(req *UpdateRequest) error {
	entity := req.Entity
	diff := req.Diff

	current := s.Get(entity.Id)
	if current == nil {
		return errors.New("Cannot create entity: no such entity")
	}

	// Calculate delta
	var d *Delta
	if !req.accepted {
		d = &Delta{
			tick:      req.tick,
			requested: req.requested,
			EntityId:  entity.Id,
			Diff:      diff,
			From:      copyFromDiff(current, diff),
			To:        entity,
		}
	}

	// Apply diff
	diff.Apply(entity, current)

	// Add delta to history
	if d != nil {
		s.deltas.Insert(d)
	}

	s.tick = req.tick
	req.accepted = true

	return nil
}

func (s *Service) acceptDelete(req *DeleteRequest) error {
	id := req.EntityId

	entity := s.Get(id)
	if entity == nil {
		return errors.New("Cannot delete entity: no such entity")
	}

	delete(s.entities, id)

	s.deltas.Insert(&Delta{
		tick:      req.tick,
		requested: req.requested,
		EntityId:  entity.Id,
		From:      entity,
		To:        nil,
	})

	s.tick = req.tick
	req.accepted = true

	return nil
}

func (s *Service) Rewind(dt message.AbsoluteTick) error {
	if dt > s.tick {
		return errors.New(fmt.Sprintf("Cannot rewind by %d: negative tick", dt))
	} else if dt < 0 {
		return errors.New(fmt.Sprintf("Cannot rewind by %d: negative rewind", dt))
	} else if dt == 0 {
		return nil
	}

	target := s.tick - dt

	// Browse deltas list backward (from newer to older)
	for e := s.deltas.LastBefore(s.tick); e != nil; e = e.Prev() {
		d := e.Value.(*Delta)

		if d.tick < target {
			// Reached target, stop here
			break
		}

		// Revert delta
		if d.From != nil {
			current := s.Get(d.EntityId)

			if current != nil {
				// The entity has been updated, revert update
				diff := d.Diff
				if diff == nil {
					diff = &message.EntityDiff{true, true, true, true}
				}

				diff.Apply(d.From, current)
			} else {
				// The entity has been deleted, restore it
				s.entities[d.EntityId] = d.From
			}
		} else {
			// The entity has been created, delete it
			delete(s.entities, d.EntityId)
		}
	}

	s.tick = target
	return nil
}

func (s *Service) Deltas() *delta.List {
	return s.deltas
}

// TODO: remove this method
func (s *Service) Redo(d *Delta) error {
	if d.tick > s.tick {
		// Newer delta, update internal tick
		s.tick = d.tick
	} else if d.tick < s.tick {
		return errors.New("Cannot redo an action in the past, rewind before")
	}

	if d.To != nil {
		if d.From != nil {
			// Update the entity
			current := s.Get(d.EntityId)
			if current == nil {
				return errors.New("Cannot update entity: no such entity")
			}
			if d.Diff == nil {
				return errors.New("Cannot update entity: no diff provided")
			}

			d.Diff.Apply(d.To, current)
		} else {
			// Create the entity
			if s.Get(d.EntityId) == nil {
				s.entities[d.EntityId] = d.To
			} else {
				return errors.New("Cannot create entity: entity already exists")
			}
		}
	} else {
		// Delete the entity
		if s.Get(d.EntityId) != nil {
			delete(s.entities, d.EntityId)
		} else {
			return errors.New("Cannot delete entity: no such entity")
		}
	}

	return nil
}

func (s *Service) Frontend() *Frontend {
	if s.frontend == nil {
		s.frontend = newFrontend(s)
	}

	return s.frontend
}

func NewService() *Service {
	return &Service{
		entities: map[message.EntityId]*message.Entity{},
		deltas:   delta.NewList(),
	}
}
