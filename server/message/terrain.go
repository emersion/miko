package message

// The number of rows/columns in a block
const BLOCK_LEN = 256

type BlockCoord int16
type PointCoord uint8
type BlockPoints [BLOCK_LEN][BLOCK_LEN]PointType

// A terrain
type Terrain interface {
	GetBlockAt(x, y BlockCoord) *Block
	SetBlock(blk *Block)
}

// A block
type Block struct {
	X BlockCoord
	Y BlockCoord
	Points *BlockPoints
}

// Get the number of points in this block
func (b *Block) Size() int {
	return BLOCK_LEN*BLOCK_LEN
}

// Change all points in this block to a given type
func (b *Block) Fill(t PointType) {
	for i := range b.Points {
		for j := range b.Points[i] {
			b.Points[i][j] = t
		}
	}
}
