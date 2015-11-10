// Provides functions to keep track of entities.
package entity

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/clock"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"time"
)

// A delta (and not a @delthas) stores all data about an entity change: not only
// its diff if it has been updated, but also the time when the change has been
// made and the entity itself.
type delta struct {
	Tick   message.AbsoluteTick
	Entity *message.Entity

	Created bool
	Deleted bool
	Diff    *message.EntityDiff
}

// The entity service.
// It manages all entities by maintaining a list of them and a diff pool. The
// diff pool keeps track of created, updated and deleted entities to send
// appropriate messages to clients.
type Service struct {
	entities map[message.EntityId]*message.Entity
	deltas   []*delta
}

func (s *Service) List() (entities []*message.Entity) {
	for _, entity := range s.entities {
		entities = append(entities, entity)
	}
	return
}

func (s *Service) Get(id message.EntityId) *message.Entity {
	if entity, ok := s.entities[id]; ok {
		return entity
	}
	return nil
}

func (s *Service) Add(entity *message.Entity, t message.AbsoluteTick) {
	if int(entity.Id) == 0 {
		nextId := len(s.entities)
		if nextId == 0 {
			nextId = 1
		}
		entity.Id = message.EntityId(nextId)
	} else if s.Get(entity.Id) != nil {
		return // TODO: error handling?
	}

	s.entities[entity.Id] = entity

	s.deltas = append(s.deltas, &delta{
		Tick:    t,
		Entity:  entity,
		Created: true,
	})
}

func (s *Service) Update(entity *message.Entity, diff *message.EntityDiff, t message.AbsoluteTick) {
	diff.Apply(entity, s.entities[entity.Id])
	entity = s.entities[entity.Id]

	s.deltas = append(s.deltas, &delta{
		Tick:   t,
		Entity: entity,
		Diff:   diff,
	})
}

func (s *Service) Delete(id message.EntityId, t message.AbsoluteTick) {
	entity := s.entities[id]
	delete(s.entities, id)

	s.deltas = append(s.deltas, &delta{
		Tick:    t,
		Entity:  entity,
		Deleted: true,
	})
}

// Check if the diff pool is empty. If not, it means that entities updates need
// to be sent to clients.
func (s *Service) IsDirty() bool {
	return len(s.deltas) > 0
}

// Flush the diff pool. This returns the current one and replace it by a new one.
func (s *Service) Flush() *message.EntityDiffPool {
	diff := message.NewEntityDiffPool()

	// TODO: sort deltas by tick

	for _, d := range s.deltas {
		entity := s.entities[d.Entity.Id]
		if d.Created {
			diff.Created = append(diff.Created, entity)
		} else if d.Deleted {
			diff.Deleted = append(diff.Deleted, entity.Id)
		} else {
			// Entity has been updated
			if _, ok := diff.Updated[entity]; ok {
				diff.Updated[entity].Merge(d.Diff)
			} else {
				diff.Updated[entity] = d.Diff
			}
		}
	}

	s.deltas = []*delta{}

	return diff
}

func (s *Service) Animate(trn message.Terrain, clk message.ClockService) {
	mover := NewMover(trn, clk)

	for {
		start := time.Now().UnixNano()
		clk.Tick()

		for _, entity := range s.entities {
			diff := mover.UpdateEntity(entity)
			if diff != nil {
				s.Update(entity, diff, clk.GetAbsoluteTick())
			}
		}

		end := time.Now().UnixNano()
		time.Sleep(clock.TickDuration - time.Nanosecond*time.Duration(end-start))
	}
}

func NewService() message.EntityService {
	return &Service{
		entities: map[message.EntityId]*message.Entity{},
	}
}
