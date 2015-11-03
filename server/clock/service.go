package clock

import "time"

const TickDuration = time.Microsecond * 20

type Service struct {
	ticks int64
}

func (s *Service) Tick() {
	s.ticks++
}

func (s *Service) GetTicks() int64 {
	return s.ticks
}

func NewService() *Service {
	return &Service{}
}
