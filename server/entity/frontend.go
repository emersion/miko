package entity

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"log"
)

const frontendChanSize = 128

type Frontend struct {
	backend *Service
	deltas  []*Delta

	Creates chan *CreateRequest
	Updates chan *UpdateRequest
	Deletes chan *DeleteRequest
}

func (f *Frontend) Get(id message.EntityId) *message.Entity {
	return f.backend.Get(id).ToMessage()
}

func (f *Frontend) Add(entity *message.Entity, t message.AbsoluteTick) message.Request {
	entity.Id = f.backend.allocateId()
	req := &CreateRequest{newClientRequest(t), NewFromMessage(entity)}
	f.Creates <- req
	return req
}

func (f *Frontend) Update(entity *message.Entity, diff *message.EntityDiff, t message.AbsoluteTick) message.Request {
	req := &UpdateRequest{newClientRequest(t), NewFromMessage(entity), diff}
	f.Updates <- req
	return req
}

func (f *Frontend) Delete(id message.EntityId, t message.AbsoluteTick) message.Request {
	req := &DeleteRequest{newClientRequest(t), id}
	f.Deletes <- req
	return req
}

// Check if the diff pool is empty. If not, it means that entities updates need
// to be sent to clients.
func (f *Frontend) IsDirty() bool {
	return (len(f.deltas) > 0)
}

// Flush the diff pool. This returns the current one and replace it by a new one.
func (f *Frontend) Flush() *message.EntityDiffPool {
	log.Println("Flushing entity frontend, deltas count:", len(f.deltas))
	pool := deltasToDiffPool(flattenDeltas(f.deltas))
	f.deltas = []*Delta{}
	return pool
}

func newFrontend(backend *Service) *Frontend {
	return &Frontend{
		backend: backend,
		Creates: make(chan *CreateRequest, frontendChanSize),
		Updates: make(chan *UpdateRequest, frontendChanSize),
		Deletes: make(chan *DeleteRequest, frontendChanSize),
	}
}
