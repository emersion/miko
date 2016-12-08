package message

import (
	"time"
)

// An absolute tick. It starts at 1 when the server is ready, and is incremented
// each time a tick is performed. The special value 0 indicates that the field
// is not set.
type AbsoluteTick uint64

// A relative tick. Absolute ticks are not designed to be sent over the network
// because of their size. Instead, a smaller relative tick is used. Because of
// its capacity, it reaches much more quickly its maximum value, so it is
// regularly reset.
type Tick uint16

// The maximum value for a relative tick.
const MaxTick = 65536

// The maximum interval of time that can be rewinded by a Rewindable.
const MaxRewind AbsoluteTick = 20

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

	// Get the time at which the current tick has started.
	GetTickTime() time.Time

	// Synchronize the internal clock with another one.
	// Most of the time, it is used to synchronize the client's clock with the
	// server's one.
	Sync(t Tick)
}

// TODO: move Rewindable elsewhere?

// Exposes functions to restore its state in the past, and replay changes to get
// back in the present.
type Rewindable interface {
	// Get the current time.
	GetTick() AbsoluteTick

	// Go back in the past.
	Rewind(dt AbsoluteTick) error

	// Get back in the future (relative to the past!).
	FastForward(dt AbsoluteTick) error
}

// A Unix timestamp in milliseconds.
type Timestamp uint64

// Convert a Time into a Timestamp.
func TimeToTimestamp(t time.Time) Timestamp {
	return Timestamp(t.Unix())*1000000 + Timestamp(t.Nanosecond()/1000)
}

// Convert a Timestamp into a Time.
func TimestampToTime(t Timestamp) time.Time {
	sec := int64(t / 1000000)
	nsec := int64((t - Timestamp(sec)*1000000) * 1000)
	return time.Unix(sec, nsec)
}
