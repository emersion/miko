package clock

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
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

func (s *Service) ToRelativeTick(t uint64) message.Tick {
	return message.Tick(t)
}

func (s *Service) ToAbsoluteTick(mt message.Tick) uint64 {
	return uint64(mt) + uint64(s.ticks/(message.MaxTick+1))
}

func NewService() *Service {
	return &Service{}
}
