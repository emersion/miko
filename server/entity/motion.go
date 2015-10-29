package entity

import (
	"math"
	"git.emersion.fr/saucisse-royale/miko/server/message"
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
		X: message.PointCoord(x),
		Y: message.PointCoord(y),
	}
}

func NewPositionFromMessage(coords *message.Position) *Position {
	return &Position{
		X: float64(int(coords.BX) * message.BLOCK_LEN + int(coords.X)),
		Y: float64(int(coords.BY) * message.BLOCK_LEN + int(coords.Y)),
	}
}

type Speed struct {
	Angle float64
	Norm float64
}

// Get the position reached by an object at t+dt if it has this speed during dt.
// Returns nil if it hasn't changed.
func (s *Speed) GetNextPosition(current *Position, dt float64) *Position {
	if s.Norm == 0 {
		return nil
	}

	sx, sy := s.Norm*math.Cos(s.Angle), s.Norm*math.Sin(s.Angle)

	return &Position{
		X: current.X + sx * dt,
		Y: current.Y + sy * dt,
	}
}

func (s *Speed) ToMessage() *message.Speed {
	return &message.Speed{
		Angle: float32(s.Angle),
		Norm: float32(s.Norm),
	}
}

func NewSpeedFromMessage(speed *message.Speed) *Speed {
	return &Speed{
		Angle: float64(speed.Angle),
		Norm: float64(speed.Norm),
	}
}
