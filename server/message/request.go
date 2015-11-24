package message

type Request interface {
	GetTick() AbsoluteTick
	Wait() error
	Done(error)
}
