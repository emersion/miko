package entity

import(
	"git.emersion.fr/saucisse-royale/miko/server/message"
)

// The entity service
// It manages all entities by maintaining a list of them and a diff pool. The
// diff pool keeps track of created, updated and deleted entities to send
// appropriate messages to clients.
type EntityService struct {
	entities []*message.Entity
	diff *message.EntityDiffPool
}

func (s *EntityService) Get(id message.EntityId) *message.Entity {
	if int(id) >= len(s.entities) {
		return nil
	}
	return s.entities[id]
}

func (s *EntityService) Add(entity *message.Entity) {
	entity.Id = message.EntityId(len(s.entities))
	s.entities = append(s.entities, entity)

	s.diff.Created = append(s.diff.Created, entity)
}

func (s *EntityService) Update(entity *message.Entity, diff *message.EntityDiff) {
	s.entities[entity.Id] = entity

	if _, ok := s.diff.Updated[entity.Id]; ok {
		s.diff.Updated[entity.Id].Merge(diff)
	} else {
		s.diff.Updated[entity.Id] = diff
	}
}

func (s *EntityService) Delete(id message.EntityId) {
	s.entities[id] = nil

	s.diff.Deleted = append(s.diff.Deleted, id)
}

// Check if the diff pool is empty. If not, it means that entities updates need
// to be sent to clients.
func (s *EntityService) IsDirty() bool {
	return len(s.diff.Created) > 0 || len(s.diff.Updated) > 0 || len(s.diff.Deleted) > 0
}

// Flush the diff pool. This returns the current one and replace it by a new one.
func (s *EntityService) Flush() *message.EntityDiffPool {
	diff := s.diff
	s.diff = &message.EntityDiffPool{}
	return diff
}

func NewService() *EntityService {
	return &EntityService{
		diff: &message.EntityDiffPool{},
	}
}
