// Provides functions to manage the game terrain.
package terrain

import (
	"container/list"
	"errors"
	"fmt"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

// A terrain's point change.
type delta struct {
	X    int
	Y    int
	Tick message.AbsoluteTick
	From message.PointType
	To   message.PointType
}

// The terrain.
// It is a database of the game map. It offers functions to retrieve specific
// parts of it.
type Terrain struct {
	points [][]message.PointType
	tick   message.AbsoluteTick
	deltas *list.List
}

// Get the last terrain update time.
func (t *Terrain) GetTick() message.AbsoluteTick {
	return t.tick
}

func (t *Terrain) GetBounds() (int, int, int, int) {
	return 0, 0, len(t.points), len(t.points[0])
}

// Check if the terrain has a point at given coordinates.
func (t *Terrain) hasPointAt(x, y int) bool {
	x1, y1, x2, y2 := t.GetBounds()
	return x >= x1 && y >= y1 && x < x2 && y < y2
}

// Check if the terrain has a block at given coordinates.
func (t *Terrain) hasBlockAt(bx, by message.BlockCoord) bool {
	x := int(bx) * message.BlockLen
	y := int(by) * message.BlockLen

	return t.hasPointAt(x, y) && t.hasPointAt(x+message.BlockLen-1, y+message.BlockLen-1)
}

// Get the block at a specific position.
func (t *Terrain) GetBlockAt(bx, by message.BlockCoord) (*message.Block, error) {
	if !t.hasBlockAt(bx, by) {
		return nil, errors.New(fmt.Sprintf("Cannot get block at [%d %d]: out of range", bx, by))
	}

	pts := &message.BlockPoints{}

	for i := 0; i < message.BlockLen; i++ {
		for j := 0; j < message.BlockLen; j++ {
			pts[i][j] = t.points[int(bx)*message.BlockLen+i][int(by)*message.BlockLen+j]
		}
	}

	return &message.Block{
		X:      bx,
		Y:      by,
		Points: pts,
	}, nil
}

// Get a point at a specific absolute position.
func (t *Terrain) GetPointAt(x, y int) (message.PointType, error) {
	if !t.hasPointAt(x, y) {
		return 0, errors.New(fmt.Sprintf("Cannot get point at [%d %d]: out of range", x, y))
	}

	return t.points[x][y], nil
}

// Update the map with a new block.
func (t *Terrain) SetBlock(blk *message.Block, tick message.AbsoluteTick) error {
	if !t.hasBlockAt(blk.X, blk.Y) {
		return errors.New(fmt.Sprintf("Cannot set block at [%d %d]: out of range", blk.X, blk.Y))
	}

	for i := range blk.Points {
		for j := range blk.Points[i] {
			x := int(blk.X)*message.BlockLen + i
			y := int(blk.Y)*message.BlockLen + j

			t.setPointAt(x, y, blk.Points[i][j], tick)
		}
	}

	return nil
}

// Same as SetPointAt(), but without safety checks.
func (t *Terrain) setPointAt(x, y int, pt message.PointType, tick message.AbsoluteTick) {
	if t.tick < tick {
		t.tick = tick
	}

	// If tick is set to zero or if the change is too old, do not keep track of it
	minTick := t.tick - message.MaxRewind
	if int(tick) != 0 && tick > minTick && t.points[x][y] != pt {
		// Cleanup old deltas
		for e := t.deltas.Front(); e != nil; e = e.Next() {
			d := e.Value.(*delta)

			if d.Tick < minTick {
				t.deltas.Remove(e)
			}

			// Make sure to insert the new delta at the right position: we want to keep
			// the list ordered
			if d.Tick > tick {
				t.deltas.InsertBefore(&delta{
					X:    x,
					Y:    y,
					Tick: tick,
					From: t.points[x][y],
					To:   pt,
				}, e)
				break
			}
		}
	}

	t.points[x][y] = pt
}

// Set a point of the terrain.
func (t *Terrain) SetPointAt(x, y int, pt message.PointType, tick message.AbsoluteTick) error {
	if !t.hasPointAt(x, y) {
		return errors.New(fmt.Sprintf("Cannot set point at [%d %d]: out of range", x, y))
	}

	t.setPointAt(x, y, pt, tick)
	return nil
}

// Reset the terrain with a given number of blocks.
func (t *Terrain) Reset(blkNbr int) {
	t.tick = 0
	t.points = make([][]message.PointType, blkNbr*message.BlockLen)
	for i := range t.points {
		t.points[i] = make([]message.PointType, blkNbr*message.BlockLen)
	}
}

// Auto-generate a new terrain.
func (t *Terrain) Generate() {
	for i := 0; i < 20; i++ {
		t.points[10][10+i] = message.PointType(1)
	}
}

func (t *Terrain) Rewind(dt message.AbsoluteTick) error {
	if dt > t.tick {
		return errors.New(fmt.Sprintf("Cannot rewind by %d: negative tick", dt))
	}

	target := t.tick - dt

	// Browse deltas list backward (from newer to older)
	for e := t.deltas.Back(); e != nil; e = e.Prev() {
		d := e.Value.(*delta)

		if d.Tick > t.tick {
			// This delta is in the future, ignore it
			continue
		}
		if d.Tick < target {
			// Reached target, stop here
			break
		}

		t.points[d.X][d.Y] = d.From
	}

	t.tick = target
	return nil
}

func (t *Terrain) FastForward(dt message.AbsoluteTick) error {
	target := t.tick + dt

	for e := t.deltas.Front(); e != nil; e = e.Next() {
		d := e.Value.(*delta)

		if d.Tick < t.tick {
			// This delta is in the past, ignore it.
			continue
		}
		if d.Tick > target {
			// Reached target, stop here
			break
		}

		t.points[d.X][d.Y] = d.To
	}

	t.tick = target
	return nil
}

// Create a new terrain.
func New() *Terrain {
	t := &Terrain{
		deltas: list.New(),
	}
	t.Reset(2)
	return t
}
