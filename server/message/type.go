package message

type Type uint8

const (
	Ping Type = 0
	Exit Type = 1
	Login Type = 2
)

func GetRespType(t Type) Type {
	return Type(t + 128)
}
