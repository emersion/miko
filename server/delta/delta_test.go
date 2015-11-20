package delta_test

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/delta"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"testing"
)

type mockDelta struct {
	tick message.AbsoluteTick
}

func (d *mockDelta) GetTick() message.AbsoluteTick {
	return d.tick
}

func TestList(t *testing.T) {
	l := delta.NewList()

	unsorted := []message.AbsoluteTick{5, 13, 11, 2, 8}
	sorted := []message.AbsoluteTick{2, 5, 8, 11, 13}

	for _, tick := range unsorted {
		l.Insert(&mockDelta{tick: tick})
	}

	if l.Len() != len(unsorted) {
		t.Error("Wrong list length, expected", len(unsorted), "but got", l.Len())
	}

	i := 0
	for el := l.First(); el != nil; el = el.Next() {
		tick := el.Value.(*mockDelta).tick
		if tick != sorted[i] {
			t.Error("Wrong delta at index", i, ", expected", sorted[i], "but got", tick)
		}
		i++
	}

	if l.LastBefore(9).Value.(*mockDelta).tick != 8 {
		t.Error("Wrong last delta before tick", 9)
	}
	if l.LastBefore(5).Value.(*mockDelta).tick != 5 {
		t.Error("Wrong last delta before tick", 5)
	}
	if l.LastBefore(1) != nil {
		t.Error("Wrong last delta before tick", 1)
	}

	if l.FirstAfter(10).Value.(*mockDelta).tick != 11 {
		t.Error("Wrong first delta after tick", 10)
	}
	if l.FirstAfter(11).Value.(*mockDelta).tick != 13 {
		t.Error("Wrong first delta after tick", 11)
	}
	if l.FirstAfter(15) != nil {
		t.Error("Wrong first delta after tick", 15)
	}

	l.Cleanup(8 + message.MaxRewind)
	if l.Len() != 3 {
		t.Error("Wrong list length after cleanup, expected", 3, "but got", l.Len())
	}
}
