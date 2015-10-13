package message

const BLOCK_LEN = 256

type BlockCoord int16
type PointCoord uint16
type BlockPoints [BLOCK_LEN][BLOCK_LEN]PointType

type Terrain interface {
	GetBlockAt(x, y BlockCoord) *Block
}

type Block struct {
	X BlockCoord
	Y BlockCoord
	Points *BlockPoints
}

func (b *Block) Size() int {
	return BLOCK_LEN*BLOCK_LEN
}

func (b *Block) Fill(t PointType) {
	for i := range b.Points {
		for j := range b.Points[i] {
			b.Points[i][j] = t
		}
	}
}
