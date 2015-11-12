// Utilities to manage changes between two states.
package delta

import (
	"container/list"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

// A delta is a change between two states.
type Delta interface {
	// The time when the change occured.
	GetTick() message.AbsoluteTick

	//Prev(current interface{}) interface{}
	//Next(current interface{}) interface{}
}

// A list of deltas.
type List struct {
	deltas *list.List
}

// Get the oldest delta in this list.
func (l *List) First() *list.Element {
	return l.deltas.Front()
}

// Get the newest delta in this list.
func (l *List) Last() *list.Element {
	return l.deltas.Back()
}

// Get the number of deltas in this list.
func (l *List) Len() int {
	return l.deltas.Len()
}

// Cleanup deltas that are too old.
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

// Insert a new delta to this list.
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

// Get the newest delta before a specific time.
func (l *List) LastBefore(t message.AbsoluteTick) *list.Element {
	for e := l.Last(); e != nil; e = e.Prev() {
		d := e.Value.(Delta)

		if d.GetTick() <= t {
			return e
		}
	}
	return nil
}

// Get the oldest delta after a specific time.
func (l *List) FirstAfter(t message.AbsoluteTick) *list.Element {
	for e := l.First(); e != nil; e = e.Next() {
		d := e.Value.(Delta)

		if d.GetTick() >= t {
			return e
		}
	}
	return nil
}

// Create a new list of deltas.
func NewList() *List {
	return &List{
		deltas: list.New(),
	}
}
