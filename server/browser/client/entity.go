package client

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/emersion/go-js-canvas"

	"git.emersion.fr/saucisse-royale/miko/server/message"
	"git.emersion.fr/saucisse-royale/miko/server/entity"
)

type EntityService struct {
	entity.Service
	Canvas *canvas.Canvas
}

func (s *EntityService) DrawEntity(entity *message.Entity) {
	s.Canvas.SetFillStyle("red")
	s.Canvas.FillRect(int(entity.Position.X) * res, int(entity.Position.Y) * res, res, res)
}

func (s *EntityService) Draw() {
	s.Canvas.ClearRect(0, 0, 500, 500)
	for _, entity := range s.List() {
		s.DrawEntity(entity)
	}
}

func (s *EntityService) Add(entity *message.Entity) {
	s.Service.Add(entity)
	s.DrawEntity(entity)
}

func NewEntityService(el *js.Object) *EntityService {
	return &EntityService{*entity.NewService(),canvas.New(el)}
}
