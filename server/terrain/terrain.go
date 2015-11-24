// Provides functions to manage the game terrain.
package terrain

import (
	"errors"
	"fmt"
	"git.emersion.fr/saucisse-royale/miko.git/server/delta"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

// A terrain's point change.
type Delta struct {
	tick message.AbsoluteTick
	X    int
	Y    int
	From message.PointType
	To   message.PointType
}

func (d *Delta) GetTick() message.AbsoluteTick {
	return d.tick
}

func (d *Delta) Inverse() delta.Delta {
	return &Delta{
		tick: d.tick,
		X:    d.X,
		Y:    d.Y,
		From: d.To,
		To:   d.From,
	}
}

// The terrain.
// It is a database of the game map. It offers functions to retrieve specific
// parts of it.
type Terrain struct {
	points [][]message.PointType
	tick   message.AbsoluteTick
	deltas *delta.List
}

// Get the last terrain update time.
func (t *Terrain) GetTick() message.AbsoluteTick {
	return t.tick
}

func (t *Terrain) GetBounds() (int, int, int, int) {
	return 0, 0, len(t.points), len(t.points[0])
}

func (t *Terrain) GetCenter() (int, int) {
	x1, y1, x2, y2 := t.GetBounds()
	return int((x2 - x1) / 2), int((y2 - y1) / 2)
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

	// If tick is set to zero, do not keep track of it
	if int(tick) != 0 && t.points[x][y] != pt {
		t.deltas.Insert(&Delta{
			X:    x,
			Y:    y,
			tick: tick,
			From: t.points[x][y],
			To:   pt,
		})

		t.deltas.Cleanup(t.tick)
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
	t.deltas = delta.NewList()
	t.points = make([][]message.PointType, blkNbr*message.BlockLen)

	for i := range t.points {
		t.points[i] = make([]message.PointType, blkNbr*message.BlockLen)
	}
}

// Auto-generate a new terrain.
func (t *Terrain) Generate() {
	cx, cy := t.GetCenter()
	radius := 5 * message.BlockLen
	pt := message.PointType(1)

	for i := cx - radius; i < cx+radius; i++ {
		t.points[i][cy-radius] = pt
	}
	for i := cx - radius; i < cx+radius; i++ {
		t.points[i][cy+radius] = pt
	}
	for j := cx - radius; j < cx+radius; j++ {
		t.points[cx-radius][j] = pt
	}
	for j := cx - radius; j < cx+radius; j++ {
		t.points[cx+radius][j] = pt
	}
}

func (t *Terrain) Rewind(dt message.AbsoluteTick) error {
	if dt > t.tick {
		return errors.New(fmt.Sprintf("Cannot rewind by %d: negative tick", dt))
	} else if dt < 0 {
		return errors.New(fmt.Sprintf("Cannot rewind by %d: negative rewind", dt))
	} else if dt == 0 {
		return nil
	}

	target := t.tick - dt

	// Browse deltas list backward (from newer to older)
	for e := t.deltas.LastBefore(t.tick); e != nil; e = e.Prev() {
		d := e.Value.(*Delta)

		if d.tick < target {
			// Reached target, stop here
			break
		}

		t.points[d.X][d.Y] = d.From
	}

	t.tick = target
	return nil
}

func (t *Terrain) Deltas() *delta.List {
	return t.deltas
}

func (t *Terrain) Redo(d *Delta) error {
	t.points[d.X][d.Y] = d.To
	t.tick = d.tick

	return nil
}

// Create a new terrain.
func New() *Terrain {
	t := &Terrain{}
	t.Reset(40)
	return t
}
