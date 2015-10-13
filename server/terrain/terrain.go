package terrain

import (
	"git.emersion.fr/saucisse-royale/miko/server/message"
)

const DEFAULT_LEN = 2 * message.BLOCK_LEN

type Terrain struct {
	points [][]message.PointType
}

func (t *Terrain) GetBlockAt(x, y message.BlockCoord) *message.Block {
	pts := &message.BlockPoints{}

	for i := 0; i < message.BLOCK_LEN; i++ {
		for j := 0; j < message.BLOCK_LEN; j++ {
			pts[i][j] = t.points[int(x) + i][int(y) + j]
		}
	}

	return &message.Block{
		X: x,
		Y: y,
		Points: pts,
	}
}

func (t *Terrain) SetBlock(blk *message.Block) {
	for i := range blk.Points {
		for j := range blk.Points[i] {
			t.points[int(blk.X) + i][int(blk.Y) + j] = blk.Points[i][j]
		}
	}
}

func (t *Terrain) Generate() {
	t.points[100][100] = message.PointType(1)
}

func New() *Terrain {
	t := &Terrain{}

	t.points = make([][]message.PointType, DEFAULT_LEN)
	for i := range t.points {
		t.points[i] = make([]message.PointType, DEFAULT_LEN)
	}

	return t
}
