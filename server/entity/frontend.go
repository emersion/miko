package entity

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

const frontendChanSize = 48

type Request struct {
	Tick message.AbsoluteTick
}

type CreateRequest struct {
	Request
	Entity *message.Entity
}

type UpdateRequest struct {
	Request
	Entity *message.Entity
	Diff   *message.EntityDiff
}

type DeleteRequest struct {
	Request
	EntityId message.EntityId
}

type Frontend struct {
	backend *Service
	pool    *message.EntityDiffPool

	Creates chan *CreateRequest
	Updates chan *UpdateRequest
	Deletes chan *DeleteRequest
}

func (f *Frontend) List() map[message.EntityId]*message.Entity {
	return f.backend.List()
}

func (f *Frontend) Get(id message.EntityId) *message.Entity {
	return f.backend.Get(id)
}

func (f *Frontend) Add(entity *message.Entity, t message.AbsoluteTick) {
	req := &CreateRequest{}
	req.Tick = t
	req.Entity = entity

	f.Creates <- req
}

func (f *Frontend) Update(entity *message.Entity, diff *message.EntityDiff, t message.AbsoluteTick) {
	req := &UpdateRequest{}
	req.Tick = t
	req.Entity = entity
	req.Diff = diff

	f.Updates <- req
}

func (f *Frontend) Delete(id message.EntityId, t message.AbsoluteTick) {
	req := &DeleteRequest{}
	req.Tick = t
	req.EntityId = id

	f.Deletes <- req
}

func (f *Frontend) IsDirty() bool {
	return f.backend.IsDirty()
}

func (f *Frontend) Flush() *message.EntityDiffPool {
	return f.backend.Flush()
}

func newFrontend(backend *Service) *Frontend {
	return &Frontend{
		backend: backend,
		Creates: make(chan *CreateRequest, frontendChanSize),
		Updates: make(chan *UpdateRequest, frontendChanSize),
		Deletes: make(chan *DeleteRequest, frontendChanSize),
	}
}
