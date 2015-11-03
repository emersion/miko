package message

type ClockService interface {
	Tick()
	GetTicks() int64
}
