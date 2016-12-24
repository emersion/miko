package entity

import (
	"sync"

	"git.emersion.fr/saucisse-royale/miko.git/server/delta"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

const frontendChanSize = 128

type Frontend struct {
	backend *Service
	deltas  *delta.List
	locker  sync.Mutex

	Creates chan *CreateRequest
	Updates chan *UpdateRequest
	Deletes chan *DeleteRequest
}

func (f *Frontend) Get(id message.EntityId) *message.Entity {
	return f.backend.Get(id).ToMessage()
}

func (f *Frontend) List() []*message.Entity {
	l := []*message.Entity{}

	for _, entity := range f.backend.entities {
		l = append(l, entity.ToMessage())
	}

	return l
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
	return (f.deltas.Len() > 0)
}

// Flush the diff pool. This returns the current one and replace it by a new one.
func (f *Frontend) Flush() *message.EntityDiffPool {
	f.locker.Lock()
	flattened := flattenDeltas(f.deltas)
	f.deltas.Reset()
	f.locker.Unlock()

	pool := deltasToDiffPool(flattened)
	pool.Tick = f.backend.tick.ToRelative()
	return pool
}

func newFrontend(backend *Service) *Frontend {
	return &Frontend{
		backend: backend,
		deltas:  delta.NewList(),
		Creates: make(chan *CreateRequest, frontendChanSize),
		Updates: make(chan *UpdateRequest, frontendChanSize),
		Deletes: make(chan *DeleteRequest, frontendChanSize),
	}
}
