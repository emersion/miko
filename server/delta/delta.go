package delta

import (
	"container/list"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

type Delta interface {
	GetTick() message.AbsoluteTick
	//Prev(current interface{}) interface{}
	//Next(current interface{}) interface{}
}

type List struct {
	deltas *list.List
}

func (l *List) First() *list.Element {
	return l.deltas.Front()
}

func (l *List) Last() *list.Element {
	return l.deltas.Back()
}

func (l *List) Cleanup(t message.AbsoluteTick) {
	minTick := t - message.MaxRewind

	for e := l.First(); e != nil; e = e.Next() {
		d := e.Value.(Delta)

		if d.GetTick() < minTick {
			l.deltas.Remove(e)
		} else {
			break
		}
	}
}

func (l *List) Insert(v Delta) *list.Element {
	for e := l.Last(); e != nil; e = e.Prev() {
		d := e.Value.(Delta)

		// Make sure to insert the new delta is at the right position: we want to
		// keep the list ordered.
		if d.GetTick() < v.GetTick() {
			return l.deltas.InsertAfter(v, e)
		}
	}

	return l.deltas.PushFront(v)
}

func (l *List) LastBefore(t message.AbsoluteTick) *list.Element {
	for e := l.Last(); e != nil; e = e.Prev() {
		d := e.Value.(Delta)

		if d.GetTick() <= t {
			return e
		}
	}
	return nil
}

func (l *List) FirstAfter(t message.AbsoluteTick) *list.Element {
	for e := l.First(); e != nil; e = e.Next() {
		d := e.Value.(Delta)

		if d.GetTick() >= t {
			return e
		}
	}
	return nil
}

func NewList() *List {
	return &List{
		deltas: list.New(),
	}
}
