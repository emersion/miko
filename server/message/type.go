package message

type Type uint8

var Types = map[string]Type{
	"ping": 0,
	"exit": 1,
	"login": 2,
	"register": 3,
	"playermeta": 4,
	"terrainupdate": 5,
	"terrainrequest": 6,
	"entitiesupdate": 7,
	"entityupdate": 8,
}

var LoginResponseCode = map[string]uint8{
	"ok": 0,
	"unknownpseudo": 1,
	"wrongpassword": 2,
	"toomanytries": 3,
	"alreadyconnected": 4,
	"playerlimitreached": 5,
}

var MetaActionCode = map[string]uint8{
 	"playerjoined": 0,
	"playerleft": 1,
}

func GetRespType(t Type) Type {
	return Type(t + 128)
}
