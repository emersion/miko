package entity

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/requests"
)

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
	requests.Request
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
