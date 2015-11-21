package action

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

const frontendChanSize = 48

type ActionRequest struct {
	Action
}

type Frontend struct {
	backend *Service
	actions []*Action

	Executes chan *ActionRequest
}

func (f *Frontend) Execute(a *message.Action, t message.AbsoluteTick) {
	req := &ActionRequest{NewFromMessage(a, t)}
	f.Executes <- req
}

func newFrontend(backend *Service) *Frontend {
	return &Frontend{
		backend:  backend,
		Executes: make(chan *ActionRequest, frontendChanSize),
	}
}
