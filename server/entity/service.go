// Provides functions to keep track of entities.
package entity

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/clock"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"time"
)

// The entity service
// It manages all entities by maintaining a list of them and a diff pool. The
// diff pool keeps track of created, updated and deleted entities to send
// appropriate messages to clients.
type Service struct {
	entities map[message.EntityId]*message.Entity
	diff     *message.EntityDiffPool
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
	}

	s.entities[entity.Id] = entity

	s.diff.Created = append(s.diff.Created, entity)
}

func (s *Service) Update(entity *message.Entity, diff *message.EntityDiff, t message.AbsoluteTick) {
	diff.Apply(entity, s.entities[entity.Id])
	entity = s.entities[entity.Id]

	if _, ok := s.diff.Updated[entity]; ok {
		s.diff.Updated[entity].Merge(diff)
	} else {
		s.diff.Updated[entity] = diff
	}
}

func (s *Service) Delete(id message.EntityId, t message.AbsoluteTick) {
	delete(s.entities, id)
	s.diff.Deleted = append(s.diff.Deleted, id)
}

// Check if the diff pool is empty. If not, it means that entities updates need
// to be sent to clients.
func (s *Service) IsDirty() bool {
	return len(s.diff.Created) > 0 || len(s.diff.Updated) > 0 || len(s.diff.Deleted) > 0
}

// Flush the diff pool. This returns the current one and replace it by a new one.
func (s *Service) Flush() *message.EntityDiffPool {
	diff := s.diff
	s.diff = message.NewEntityDiffPool()
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
		diff:     message.NewEntityDiffPool(),
	}
}
