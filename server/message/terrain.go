package message

const BLOCK_LEN = 256

type Block struct {
	X int16
	Y int16
	Points [BLOCK_LEN][BLOCK_LEN]PointType
}

func (b *Block) Size() int {
	return BLOCK_LEN*BLOCK_LEN
}
