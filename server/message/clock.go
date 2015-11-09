package message

// An absolute tick. It starts at 0 when the server is ready, and is incremented
// each time a tick is performed.
type AbsoluteTick uint64

// A relative tick. Absolute ticks are not designed to be sent over the network
// because of their size. Instead, a smaller relative tick is used. Because of
// its capacity, it raises much more quickly its maximum value, so it is
// regularly reset.
type Tick uint16

// The maximum value for a relative tick.
const MaxTick = 65536

// A clock service keeps track of the server's current in-game time.
type ClockService interface {
	// Trigger a new tick.
	Tick()
	// Get the current absolute tick.
	GetAbsoluteTick() AbsoluteTick
	// Get the current relative tick.
	GetRelativeTick() Tick
}
