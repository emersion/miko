// Provides functions to manage time, relative to the in-game world.
package clock

import (
	"time"

	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

// The duration of a single tick.
const TickDuration = time.Millisecond * 20

// The clock service.
type Service struct {
	startedAt time.Time
	ticks     message.AbsoluteTick
}

// Increment the clock's tick count.
func (s *Service) Tick() {
	if s.ticks == 0 {
		s.startedAt = time.Now().Add(-TickDuration)
	}
	s.ticks++
}

// Get the current absolute tick.
func (s *Service) GetAbsoluteTick() message.AbsoluteTick {
	return s.ticks
}

// Get the current relative tick.
func (s *Service) GetRelativeTick() message.Tick {
	return s.ToRelativeTick(s.ticks)
}

// Get the time at which the current tick has started.
func (s *Service) GetTickTime() time.Time {
	return s.startedAt.Add(TickDuration * time.Duration(s.ticks))
}

// Convert an absolute tick to a relative tick.
func (s *Service) ToRelativeTick(at message.AbsoluteTick) message.Tick {
	return message.Tick(at)
}

// Convert a relative tick to an absolute tick.
// Returns message.InvalidTick if an error occured.
func (s *Service) ToAbsoluteTick(rt message.Tick) message.AbsoluteTick {
	current := s.GetRelativeTick()

	at := message.AbsoluteTick(rt) + message.AbsoluteTick(s.ticks-s.ticks%message.MaxTick)
	if current < rt {
		if at < message.MaxTick {
			return message.InvalidTick // Underflow error
		}
		at -= message.AbsoluteTick(message.MaxTick)
	}

	return at
}

func (s *Service) Sync(t message.Tick) {
	s.ticks = 0
	s.ticks = s.ToAbsoluteTick(t)
}

// Create a new clock service.
func NewService() *Service {
	return &Service{}
}
