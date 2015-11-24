package action

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/delta"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

type Action struct {
	*message.Action
	tick message.AbsoluteTick
	wait chan error
}

func (a *Action) GetTick() message.AbsoluteTick {
	return a.tick
}

func (a *Action) Wait() error {
	return <-a.wait
}

func (a *Action) Done(err error) {
	select {
	case a.wait <- err:
	default:
	}
}

func (a *Action) Execute() []delta.Delta {
	return []delta.Delta{}
}

func (a *Action) Inverse() delta.Delta {
	return &inversedAction{a}
}

func (a *Action) ToMessage() *message.Action {
	return a.Action
}

func NewFromMessage(src *message.Action, t message.AbsoluteTick) *Action {
	return &Action{
		Action: src,
		tick:   t,
		wait:   make(chan error),
	}
}

type inversedAction struct {
	*Action
}

func (a *inversedAction) Execute() []delta.Delta {
	inversed := []delta.Delta{}
	deltas := a.Action.Execute()

	for i := len(deltas) - 1; i >= 0; i-- {
		inversed = append(inversed, deltas[i].Inverse())
	}

	return inversed
}

func (a *inversedAction) Inverse() delta.Delta {
	return a.Action
}
