// Provides functions to manage actions.
package action

import (
	"errors"
	"fmt"
	"git.emersion.fr/saucisse-royale/miko.git/server/delta"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

// An action service.
type Service struct {
	frontend *Frontend
	actions  *delta.List
	tick     message.AbsoluteTick
}

func (s *Service) GetTick() message.AbsoluteTick {
	return s.tick
}

func (s *Service) AcceptRequest(req *Request) (err error) {
	s.actions.Insert(req.Action)
	s.tick = req.GetTick()

	req.Done(err)

	if s.frontend != nil {
		s.frontend.actions = append(s.frontend.actions, req.Action)
	}

	return
}

func (s *Service) Rewind(dt message.AbsoluteTick) (deltas []delta.Delta, err error) {
	if dt > s.tick {
		err = errors.New(fmt.Sprintf("Cannot rewind by %d: negative tick", dt))
		return
	} else if dt < 0 {
		err = errors.New(fmt.Sprintf("Cannot rewind by %d: negative rewind", dt))
		return
	} else if dt == 0 {
		return
	}

	target := s.tick - dt

	// Browse actions list backward (from newer to older)
	for e := s.actions.LastBefore(s.tick); e != nil; e = e.Prev() {
		a := e.Value.(Action)
		inversed := a.Inverse().(*Action)

		for _, d := range inversed.Execute() {
			deltas = append(deltas, d)
		}

		if s.frontend != nil {
			// TODO: cancel action
		}
	}

	s.tick = target

	return
}

// Cleanup data that is too old.
func (s *Service) Cleanup(t message.AbsoluteTick) {
	s.actions.Cleanup(t)
}

func (s *Service) Actions() *delta.List {
	return s.actions
}

// Get this service's frontend.
func (s *Service) Frontend() *Frontend {
	if s.frontend == nil {
		s.frontend = newFrontend(s)
	}

	return s.frontend
}

// Create a new action service.
func NewService() *Service {
	return &Service{
		actions: delta.NewList(),
	}
}
