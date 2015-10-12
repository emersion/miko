package entity

import(
	"git.emersion.fr/saucisse-royale/miko/server/message"
)

type EntityService struct {
	entities []*message.Entity
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

func NewService() *EntityService {
	return &EntityService{}
}
