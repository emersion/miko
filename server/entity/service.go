// Provides functions to keep track of entities.
package entity

import (
	"errors"
	"fmt"

	"git.emersion.fr/saucisse-royale/miko.git/server/delta"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"

	//"log"
)

func copyEntityFromDiff(src *Entity, diff *message.EntityDiff) *Entity {
	dst := &Entity{}
	dst.ApplyDiff(diff, src)
	return dst
}

// The entity service.
// It manages all entities by maintaining a list of them and a diff pool. The
// diff pool keeps track of created, updated and deleted entities to send
// appropriate messages to clients.
type Service struct {
	entities map[message.EntityId]*Entity
	lastId   message.EntityId
	deltas   *delta.List
	tick     message.AbsoluteTick
	frontend *Frontend
}

func (s *Service) GetTick() message.AbsoluteTick {
	return s.tick
}

func (s *Service) List() map[message.EntityId]*Entity {
	return s.entities
}

func (s *Service) Get(id message.EntityId) *Entity {
	if entity, ok := s.entities[id]; ok {
		return entity
	}
	return nil
}

func (s *Service) AcceptRequest(req Request) (err error) {
	switch r := req.(type) {
	case *CreateRequest:
		err = s.acceptCreate(r)
	case *UpdateRequest:
		err = s.acceptUpdate(r)
	case *DeleteRequest:
		err = s.acceptDelete(r)
	default:
		panic("Cannot accept request: not an entity request")
	}

	req.Done(err)
	return err
}

func (s *Service) allocateId() message.EntityId {
	s.lastId++
	return s.lastId
}

func (s *Service) acceptCreate(req *CreateRequest) error {
	entity := req.Entity

	if entity.Id == 0 {
		entity.Id = s.allocateId()
	} else if s.Get(entity.Id) != nil {
		return errors.New("Cannot create entity: entity already exists")
	}

	s.entities[entity.Id] = entity

	d := &Delta{
		tick:      req.tick,
		requested: req.requested,
		EntityId:  entity.Id,
		From:      nil,
		To:        entity,
	}
	if req.requested {
		s.deltas.Insert(d)
		//log.Println("Accepted create request, deltas count:", s.deltas.Len())
	}
	if s.frontend != nil {
		s.frontend.locker.Lock()
		s.frontend.deltas.Insert(d)
		s.frontend.locker.Unlock()
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
		return errors.New("Cannot update entity: no such entity")
	}

	// Calculate delta
	d := &Delta{
		tick:      req.tick,
		requested: req.requested,
		EntityId:  entity.Id,
		Diff:      diff,
		From:      copyEntityFromDiff(current, diff),
		To:        entity,
	}

	// Apply diff
	current.ApplyDiff(diff, entity)

	// Add delta to history
	if req.requested {
		s.deltas.Insert(d)
		//log.Println("Accepted update request for entity", entity.Id)
	}
	// TODO: fix this
	// && !current.ToMessage().EqualsWithDiff(entity.ToMessage(), diff)
	if s.frontend != nil {
		s.frontend.locker.Lock()
		s.frontend.deltas.Insert(d)
		s.frontend.locker.Unlock()
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

	d := &Delta{
		tick:      req.tick,
		requested: req.requested,
		EntityId:  entity.Id,
		From:      entity,
		To:        nil,
	}
	if req.requested {
		s.deltas.Insert(d)
		//log.Println("Accepted delete request, deltas count:", s.deltas.Len())
	}
	if s.frontend != nil {
		s.frontend.locker.Lock()
		s.frontend.deltas.Insert(d)
		s.frontend.locker.Unlock()
	}

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

		if d.tick <= target {
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
					diff = message.NewFilledEntityDiff(true)
				}

				current.ApplyDiff(diff, d.From)
			} else {
				// The entity has been deleted, restore it
				s.entities[d.EntityId] = d.From
			}
		} else {
			// The entity has been created, delete it
			delete(s.entities, d.EntityId)
		}

		if s.frontend != nil {
			s.frontend.locker.Lock()
			s.frontend.deltas.Insert(d.Inverse().(*Delta))
			s.frontend.locker.Unlock()
		}
	}

	s.tick = target

	return nil
}

// Cleanup data that is too old.
func (s *Service) Cleanup(t message.AbsoluteTick) {
	s.deltas.Cleanup(t)
}

// Get this service's deltas.
func (s *Service) Deltas() *delta.List {
	return s.deltas
}

// TODO: remove this method
func (s *Service) redo(d *Delta) error {
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

			current.ApplyDiff(d.Diff, d.To)
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

// Get this service's frontend.
func (s *Service) Frontend() *Frontend {
	if s.frontend == nil {
		s.frontend = newFrontend(s)
	}

	return s.frontend
}

// Create a new entity service.
func NewService() *Service {
	return &Service{
		entities: map[message.EntityId]*Entity{},
		deltas:   delta.NewList(),
	}
}
