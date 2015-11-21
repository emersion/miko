package action

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/delta"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

type Action interface {
	delta.Delta

	Execute() []delta.Delta
}

type action struct {
	*message.Action
	tick message.AbsoluteTick
}

func (a *action) GetTick() message.AbsoluteTick {
	return a.tick
}

func (a *action) Execute() []delta.Delta {
	return []delta.Delta{}
}

func (a *action) Inverse() delta.Delta {
	return &inversedAction{a}
}

func NewFromMessage(src *message.Action, t message.AbsoluteTick) *action {
	return &action{
		Action: src,
		tick:   t,
	}
}

type inversedAction struct {
	*action
}

func (a *inversedAction) Execute() []delta.Delta {
	inversed := []delta.Delta{}
	deltas := a.action.Execute()

	for i := len(deltas) - 1; i >= 0; i-- {
		inversed = append(inversed, deltas[i].Inverse())
	}

	return inversed
}

func (a *inversedAction) Inverse() delta.Delta {
	return a.action
}
