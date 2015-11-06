package message

type Tick uint16

type ClockService interface {
	Tick()
	GetTickCount() int64
	GetTick() Tick
}
