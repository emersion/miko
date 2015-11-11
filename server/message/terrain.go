package message

// The number of rows/columns in a block
const BlockLen = 256

type BlockCoord int16
type PointCoord uint8
type BlockPoints [BlockLen][BlockLen]PointType

// A terrain
type Terrain interface {
	GetBlockAt(bx, by BlockCoord) (*Block, error)
	GetPointAt(x, y int) (PointType, error)
	SetBlock(blk *Block, t AbsoluteTick) error
	SetPointAt(x, y int, pt PointType, t AbsoluteTick) error
}

// A block
type Block struct {
	X      BlockCoord
	Y      BlockCoord
	Points *BlockPoints
}

// Get the number of points in this block
func (b *Block) Size() int {
	return BlockLen * BlockLen
}

// Change all points in this block to a given type
func (b *Block) Fill(t PointType) {
	for i := range b.Points {
		for j := range b.Points[i] {
			b.Points[i][j] = t
		}
	}
}
