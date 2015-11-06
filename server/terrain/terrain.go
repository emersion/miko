package terrain

import (
	"errors"
	"fmt"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

// The terrain
// It is a database of the game map. It offers functions to retrieve specific
// parts of it.
type Terrain struct {
	Points [][]message.PointType
}

func (t *Terrain) hasPointAt(x, y int) bool {
	return x >= 0 && y >= 0 && x < len(t.Points) && y < len(t.Points[0])
}

func (t *Terrain) hasBlockAt(bx, by message.BlockCoord) bool {
	x := int(bx) * message.BLOCK_LEN
	y := int(by) * message.BLOCK_LEN

	return t.hasPointAt(x, y) && t.hasPointAt(x+message.BLOCK_LEN-1, y+message.BLOCK_LEN-1)
}

// Get the block at a specific position.
func (t *Terrain) GetBlockAt(bx, by message.BlockCoord) (*message.Block, error) {
	if !t.hasBlockAt(bx, by) {
		return nil, errors.New(fmt.Sprintf("Cannot get block at [%d %d]: out of range", bx, by))
	}

	pts := &message.BlockPoints{}

	for i := 0; i < message.BLOCK_LEN; i++ {
		for j := 0; j < message.BLOCK_LEN; j++ {
			pts[i][j] = t.Points[int(bx)*message.BLOCK_LEN+i][int(by)*message.BLOCK_LEN+j]
		}
	}

	return &message.Block{
		X:      bx,
		Y:      by,
		Points: pts,
	}, nil
}

func (t *Terrain) GetPointAt(x, y int) (message.PointType, error) {
	if !t.hasPointAt(x, y) {
		return 0, errors.New(fmt.Sprintf("Cannot get point at [%d %d]: out of range", x, y))
	}

	return t.Points[x][y], nil
}

// Update the map with a new block.
func (t *Terrain) SetBlock(blk *message.Block) error {
	if !t.hasBlockAt(blk.X, blk.Y) {
		return errors.New(fmt.Sprintf("Cannot set block at [%d %d]: out of range", blk.X, blk.Y))
	}

	for i := range blk.Points {
		for j := range blk.Points[i] {
			t.Points[int(blk.X)*message.BLOCK_LEN+i][int(blk.Y)*message.BLOCK_LEN+j] = blk.Points[i][j]
		}
	}

	return nil
}

func (t *Terrain) Reset(blkNbr int) {
	t.Points = make([][]message.PointType, blkNbr*message.BLOCK_LEN)
	for i := range t.Points {
		t.Points[i] = make([]message.PointType, blkNbr*message.BLOCK_LEN)
	}
}

func (t *Terrain) Generate() {
	for i := 0; i < 20; i++ {
		t.Points[10][10+i] = message.PointType(1)
	}
}

func New() *Terrain {
	t := &Terrain{}
	t.Reset(2)
	return t
}
