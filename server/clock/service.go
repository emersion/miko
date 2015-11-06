package clock

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"time"
)

const TickDuration = time.Millisecond * 500

type Service struct {
	ticks int64
}

func (s *Service) Tick() {
	s.ticks++
}

func (s *Service) GetTickCount() int64 {
	return s.ticks
}

func (s *Service) GetTick() message.Tick {
	return message.Tick(s.ticks)
}

func NewService() *Service {
	return &Service{}
}
