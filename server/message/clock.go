package message

type Tick uint16

const MaxTick = 65536

type ClockService interface {
	Tick()
	GetAbsoluteTick() uint64
	GetRelativeTick() Tick
}
