package terrain

import (
	"git.emersion.fr/saucisse-royale/miko/server/message"
)

// The terrain
// It is a database of the game map. It offers functions to retrieve specific
// parts of it.
type Terrain struct {
	Points [][]message.PointType
}

// Get the block at a specific position.
func (t *Terrain) GetBlockAt(x, y message.BlockCoord) *message.Block {
	pts := &message.BlockPoints{}

	for i := 0; i < message.BLOCK_LEN; i++ {
		for j := 0; j < message.BLOCK_LEN; j++ {
			pts[i][j] = t.Points[int(x) + i][int(y) + j]
		}
	}

	return &message.Block{
		X: x,
		Y: y,
		Points: pts,
	}
}

// Update the map with a new block.
func (t *Terrain) SetBlock(blk *message.Block) {
	for i := range blk.Points {
		for j := range blk.Points[i] {
			t.Points[int(blk.X) + i][int(blk.Y) + j] = blk.Points[i][j]
		}
	}
}

func (t *Terrain) Reset(blkNbr int) {
	t.Points = make([][]message.PointType, blkNbr * message.BLOCK_LEN)
	for i := range t.Points {
		t.Points[i] = make([]message.PointType, blkNbr * message.BLOCK_LEN)
	}
}

// Check if it's possible for an entity to move from its current position to
// another one.
func (t *Terrain) CanMove(entity *message.Entity, to *message.Position) bool {
	return true // TODO
}

func (t *Terrain) Generate() {
	t.Points[10][10] = message.PointType(1)
	t.Points[20][10] = message.PointType(1)
}

func New() *Terrain {
	t := &Terrain{}
	t.Reset(2)
	return t
}
