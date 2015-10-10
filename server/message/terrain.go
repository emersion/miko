package message

type Block struct {
	X int16
	Y int16
	Points [256][256]PointType
}

func (b *Block) Size() int {
	return 256*256
}
