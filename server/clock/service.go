// Provides functions to manage time, relative to the in-game world.
package clock

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"log"
	"time"
)

const TickDuration = time.Millisecond * 50

type Service struct {
	ticks uint64
}

func (s *Service) Tick() {
	s.ticks++
}

func (s *Service) GetAbsoluteTick() uint64 {
	return s.ticks
}

func (s *Service) GetRelativeTick() message.Tick {
	return s.ToRelativeTick(s.ticks)
}

func (s *Service) ToRelativeTick(at uint64) message.Tick {
	return message.Tick(at)
}

func (s *Service) ToAbsoluteTick(rt message.Tick) uint64 {
	current := s.GetRelativeTick()

	at := uint64(rt) + uint64(s.ticks-s.ticks%message.MaxTick)
	log.Println(rt, current)
	if current < rt {
		at -= uint64(message.MaxTick)
	}

	return at
}

func NewService() *Service {
	return &Service{}
}
