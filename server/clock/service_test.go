package clock_test

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/clock"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"testing"
)

func TestService(t *testing.T) {
	s := clock.NewService()

	if s.GetAbsoluteTick() != 0 {
		t.Error("Absolute tick not initialized at 0")
	}
	if s.GetRelativeTick() != 0 {
		t.Error("Relative tick not initialized at 0")
	}

	s.Tick()

	if s.GetAbsoluteTick() != 1 {
		t.Error("Absolute tick not incremented, should be 1 but is", s.GetAbsoluteTick())
	}
	if s.GetRelativeTick() != 1 {
		t.Error("Relative tick not incremented, should be 1 but is", s.GetRelativeTick())
	}

	for i := 0; i < message.MaxTick; i++ {
		s.Tick()
	}

	if s.GetAbsoluteTick() != 1+message.MaxTick {
		t.Error("Invalid absolute tick after message.MaxTick ticks, should be", 1+message.MaxTick, "but is", s.GetAbsoluteTick())
	}
	if s.GetRelativeTick() != 1 {
		t.Error("Invalid relative tick after message.MaxTick ticks, should be ", 1, "but is", s.GetRelativeTick())
	}
}
