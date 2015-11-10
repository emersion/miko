package message

// An absolute tick. It starts at 0 when the server is ready, and is incremented
// each time a tick is performed.
type AbsoluteTick uint64

// A relative tick. Absolute ticks are not designed to be sent over the network
// because of their size. Instead, a smaller relative tick is used. Because of
// its capacity, it reaches much more quickly its maximum value, so it is
// regularly reset.
type Tick uint16

// The maximum value for a relative tick.
const MaxTick = 65536

// A clock service keeps track of the server's current in-game time.
type ClockService interface {
	// Trigger a new tick.
	Tick()

	// Convert a relative tick to an absolute one.
	ToRelativeTick(at AbsoluteTick) Tick

	// Convert an absolute tick to a relative one.
	ToAbsoluteTick(rt Tick) AbsoluteTick

	// Get the current absolute tick.
	GetAbsoluteTick() AbsoluteTick

	// Get the current relative tick.
	GetRelativeTick() Tick

	// Synchronize the internal clock with another one.
	// Most of the time, it is used to synchronize the client's clock with the
	// server's one.
	Sync(t Tick)
}
