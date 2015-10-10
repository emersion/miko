package terrain

import (
	"../message"
)

type Terrain struct {
	points [][]message.PointType
}

func (t *Terrain) GetBlockAt(x, y message.BlockCoord) *message.Block {
	pts := message.BlockPoints{}

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

func New() *Terrain {
	return &Terrain{}
}
