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

	ticks := 42
	for i := 0; i < 42; i++ {
		s.Tick()
	}

	if int(s.GetAbsoluteTick()) != ticks {
		t.Error("Absolute tick not incremented, expected", ticks, "but got", s.GetAbsoluteTick())
	}
	if int(s.GetRelativeTick()) != ticks {
		t.Error("Relative tick not incremented, expected", ticks, "but got", s.GetRelativeTick())
	}

	// |--------------------|------------>
	//           ^     ^
	//         tick   now
	relTick := message.Tick(ticks - 6)
	absTick := s.ToAbsoluteTick(relTick)
	if int(absTick) != int(relTick) {
		t.Error("Invalid absolute tick conversion, expected", relTick, "but got", absTick)
	}

	for i := 0; i < message.MaxTick; i++ {
		s.Tick()
	}

	if int(s.GetAbsoluteTick()) != ticks+message.MaxTick {
		t.Error("Invalid absolute tick after message.MaxTick ticks, expected", 1+message.MaxTick, "but got", s.GetAbsoluteTick())
	}
	if int(s.GetRelativeTick()) != ticks {
		t.Error("Invalid relative tick after message.MaxTick ticks, expected", 1, "but got", s.GetRelativeTick())
	}

	// |--------------------|------------>
	//                          ^     ^
	//                        tick   now
	relTick = message.Tick(6)
	absTick = s.ToAbsoluteTick(relTick)
	if int(absTick) != message.MaxTick+6 {
		t.Error("Invalid absolute tick conversion, expected", message.MaxTick+6, "but got", absTick)
	}

	// |--------------------|------------>
	//                 ^               ^
	//               tick             now
	relTick = message.Tick(message.MaxTick - 42)
	absTick = s.ToAbsoluteTick(relTick)
	if int(absTick) != int(relTick) {
		t.Error("Invalid absolute tick conversion, expected", relTick, "but got", absTick)
	}
}
