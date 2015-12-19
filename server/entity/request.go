package entity

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

type request struct {
	tick      message.AbsoluteTick
	requested bool
	accepted  bool
	wait      chan error
}

func (r *request) GetTick() message.AbsoluteTick {
	return r.tick
}

func (r *request) Wait() error {
	return <-r.wait
}

func (r *request) Done(err error) {
	select {
	case r.wait <- err:
	default:
	}
}

func newRequest(t message.AbsoluteTick) *request {
	return &request{
		tick: t,
		wait: make(chan error),
	}
}

func newClientRequest(t message.AbsoluteTick) *request {
	req := newRequest(t)
	req.requested = true
	return req
}

// TODO: move this somewhere else?
func NewCreateRequest(t message.AbsoluteTick, entity *Entity) *CreateRequest {
	return &CreateRequest{newRequest(t), entity}
}

// TODO: move this somewhere else?
func NewUpdateRequest(t message.AbsoluteTick, entity *Entity, diff *message.EntityDiff) *UpdateRequest {
	return &UpdateRequest{newRequest(t), entity, diff}
}

// TODO: move this somewhere else?
func NewDeleteRequest(t message.AbsoluteTick, entityId message.EntityId) *DeleteRequest {
	return &DeleteRequest{newRequest(t), entityId}
}

type Request interface {
	message.Request
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
