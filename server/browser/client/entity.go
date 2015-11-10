package client

import (
	"github.com/emersion/go-js-canvas"
	"github.com/gopherjs/gopherjs/js"

	"git.emersion.fr/saucisse-royale/miko.git/server/entity"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

type EntityService struct {
	message.EntityService
	Canvas *canvas.Canvas
}

func (s *EntityService) DrawEntity(entity *message.Entity) {
	s.Canvas.SetFillStyle("red")
	s.Canvas.FillRect(int(entity.Position.X)*res, int(entity.Position.Y)*res, res, res)
}

func (s *EntityService) Draw() {
	s.Canvas.ClearRect(0, 0, 500, 500)
	for _, entity := range s.List() {
		s.DrawEntity(entity)
	}
}

func (s *EntityService) Add(entity *message.Entity, t message.AbsoluteTick) {
	s.EntityService.Add(entity, t)
	s.DrawEntity(entity)
}

func (s *EntityService) Update(entity *message.Entity, diff *message.EntityDiff, t message.AbsoluteTick) {
	s.EntityService.Update(entity, diff, t)
	s.Draw()
}

func (s *EntityService) Delete(id message.EntityId, t message.AbsoluteTick) {
	s.EntityService.Delete(id, t)
	s.Draw()
}

func NewEntityService(el *js.Object) *EntityService {
	return &EntityService{entity.NewService(), canvas.New(el)}
}
