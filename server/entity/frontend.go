package entity

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

const frontendChanSize = 48

type request struct {
	tick      message.AbsoluteTick
	requested bool
	accepted  bool
}

func (r *request) GetTick() message.AbsoluteTick {
	return r.tick
}

func newRequest(t message.AbsoluteTick) *request {
	return &request{
		tick: t,
	}
}

func newClientRequest(t message.AbsoluteTick) *request {
	req := newRequest(t)
	req.requested = true
	return req
}

// TODO: move this somewhere else?
func NewUpdateRequest(t message.AbsoluteTick, entity *Entity, diff *message.EntityDiff) *UpdateRequest {
	return &UpdateRequest{newRequest(t), entity, diff}
}

type Request interface {
	GetTick() message.AbsoluteTick
}

type CreateRequest struct {
	*request
	Entity *Entity
}

type UpdateRequest struct {
	*request
	Entity *Entity
	Diff   *message.EntityDiff
}

type DeleteRequest struct {
	*request
	EntityId message.EntityId
}

type Frontend struct {
	backend *Service
	pool    *message.EntityDiffPool

	Creates chan *CreateRequest
	Updates chan *UpdateRequest
	Deletes chan *DeleteRequest
}

func (f *Frontend) Get(id message.EntityId) *message.Entity {
	return f.backend.Get(id).ToMessage()
}

func (f *Frontend) Add(entity *message.Entity, t message.AbsoluteTick) {
	req := &CreateRequest{newClientRequest(t), NewFromMessage(entity)}
	f.Creates <- req
}

func (f *Frontend) Update(entity *message.Entity, diff *message.EntityDiff, t message.AbsoluteTick) {
	req := &UpdateRequest{newClientRequest(t), NewFromMessage(entity), diff}
	f.Updates <- req
}

func (f *Frontend) Delete(id message.EntityId, t message.AbsoluteTick) {
	req := &DeleteRequest{newClientRequest(t), id}
	f.Deletes <- req
}

// Check if the diff pool is empty. If not, it means that entities updates need
// to be sent to clients.
func (f *Frontend) IsDirty() bool {
	return !f.pool.IsEmpty()
}

// Flush the diff pool. This returns the current one and replace it by a new one.
func (f *Frontend) Flush() *message.EntityDiffPool {
	pool := f.pool
	f.pool = message.NewEntityDiffPool()
	return pool
}

func newFrontend(backend *Service) *Frontend {
	return &Frontend{
		backend: backend,
		pool:    message.NewEntityDiffPool(),
		Creates: make(chan *CreateRequest, frontendChanSize),
		Updates: make(chan *UpdateRequest, frontendChanSize),
		Deletes: make(chan *DeleteRequest, frontendChanSize),
	}
}
