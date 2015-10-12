package entity

import(
	"git.emersion.fr/saucisse-royale/miko/server/message"
)

type entityServiceDiff struct {
	created []*message.Entity
	updated []*message.Entity
	deleted []*message.Entity
}

type EntityService struct {
	entities []*message.Entity
	diff entityServiceDiff
}

func (s *EntityService) Get(id message.EntityId) *message.Entity {
	if int(id) >= len(s.entities) {
		return nil
	}
	return s.entities[id]
}

func (s *EntityService) Append(entity *message.Entity) {
	entity.Id = message.EntityId(len(s.entities))
	s.entities = append(s.entities, entity)
}

func (s *EntityService) Update(entity *message.Entity, diff *message.EntityDiff) {
	// TODO
}

func (s *EntityService) Delete(id *message.EntityId) {
	// TODO
}

func NewService() *EntityService {
	return &EntityService{}
}
