package message

type Type uint8

const (
	Ping Type = 0
	Exit Type = 1
	Login Type = 2
)

var LoginResponseCode = map[string]uint8{
	"ok": 0,
	"unknownpseudo": 1,
	"wrongpassword": 2,
	"toomanytries": 3,
	"alreadyconnected": 4,
	"playerlimitreached": 5,
}

func GetRespType(t Type) Type {
	return Type(t + 128)
}
