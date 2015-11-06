package entity

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"math"
	"time"
)

func round(f float64) int {
	return int(math.Floor(f + .5))
}

type Position struct {
	X float64
	Y float64
}

func (p *Position) ToMessage() *message.Position {
	x := round(p.X) % message.BLOCK_LEN
	y := round(p.Y) % message.BLOCK_LEN
	bx := round((p.X - float64(x)) / message.BLOCK_LEN)
	by := round((p.Y - float64(y)) / message.BLOCK_LEN)

	return &message.Position{
		BX: message.BlockCoord(bx),
		BY: message.BlockCoord(by),
		X:  message.PointCoord(x),
		Y:  message.PointCoord(y),
	}
}

func NewPositionFromMessage(coords *message.Position) *Position {
	return &Position{
		X: float64(int(coords.BX)*message.BLOCK_LEN + int(coords.X)),
		Y: float64(int(coords.BY)*message.BLOCK_LEN + int(coords.Y)),
	}
}

type Speed struct {
	Angle float64
	Norm  float64
}

// Get the position reached by an object at t+dt if it has this speed during dt.
// Returns nil if it hasn't changed.
func (s *Speed) GetNextPosition(current *Position, dt time.Duration) *Position {
	if s.Norm == 0 {
		return nil
	}

	sx, sy := s.Norm*math.Cos(s.Angle), s.Norm*math.Sin(s.Angle)

	return &Position{
		X: current.X + sx*float64(dt)/float64(time.Second),
		Y: current.Y + sy*float64(dt)/float64(time.Second),
	}
}

func (s *Speed) ToMessage() *message.Speed {
	return &message.Speed{
		Angle: float32(s.Angle),
		Norm:  float32(s.Norm),
	}
}

func NewSpeedFromMessage(speed *message.Speed) *Speed {
	return &Speed{
		Angle: float64(speed.Angle),
		Norm:  float64(speed.Norm),
	}
}

func GetRouteBetween(from, to *Position) (route [][2]int) {
	// Distance between points
	dx := math.Abs(to.X - from.X)
	dy := math.Abs(to.Y - from.Y)

	// Direction of x & y
	xSign := 1
	ySign := 1
	if to.X < from.X {
		xSign = -1
	}
	if to.Y < from.Y {
		ySign = -1
	}

	// We will define a parameter t which will take all integer values between 0 and
	// the greatest distance. One coordinate will evolve with t, the other has a
	// shorter distance to go through. Thus, while one coordinate will increase with
	// a speed of 1, the other will have a lower speed. kx and ky are these two
	// speeds.

	var kx, ky float64
	if dx > dy {
		kx = 1
		ky = dy / dx
	} else if dx < dy {
		kx = dx / dy
		ky = 1
	} else {
		kx = 1
		ky = 1
	}

	for t := 0; t <= round(math.Max(dx, dy)); t++ {
		// Calculate the two coordinates for this parameter
		x := round(from.X + float64(xSign*t)*kx)
		y := round(from.Y + float64(ySign*t)*ky)

		route = append(route, [2]int{x, y})
	}

	return
}
