package action

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

const frontendChanSize = 48

type Frontend struct {
	backend *Service
	actions []Action

	Executes chan *Request
}

func (f *Frontend) Execute(a *message.Action, t message.AbsoluteTick) {
	req := &Request{NewFromMessage(a, t)}
	f.Executes <- req
}

func (f *Frontend) IsDirty() bool {
	return (len(f.actions) > 0)
}

func (f *Frontend) Flush() []*message.Action {
	actions := make([]*message.Action, len(f.actions))
	for i, a := range f.actions {
		actions[i] = a.(*action).ToMessage()
	}

	f.actions = []Action{}
	return actions
}

func newFrontend(backend *Service) *Frontend {
	return &Frontend{
		backend:  backend,
		Executes: make(chan *Request, frontendChanSize),
	}
}
