package message

type Type uint8
type LoginResponseCode uint8
type RegisterResponseCode uint8
type MetaActionCode uint8

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
	"actions": 9,
	"action": 10,
	"entitycreate": 11,
	"entitiydestroy": 12,
	"chatsend": 13,
	"chatreceive": 14,
}

var LoginResponseCodes = map[string]LoginResponseCode{
	"ok": 0,
	"unknownpseudo": 1,
	"wrongpassword": 2,
	"toomanytries": 3,
	"alreadyconnected": 4,
	"playerlimitreached": 5,
}

var RegisterResponseCodes = map[string]RegisterResponseCode{
	"ok": 0,
	"usedpseudo": 1,
	"invalidpseudo": 2,
	"invalidpassword": 2,
	"toomanytries": 3,
	"registerdisabled": 4,
}

var MetaActionCodes = map[string]MetaActionCode{
 	"playerjoined": 0,
	"playerleft": 1,
}

func GetRespType(t Type) Type {
	return Type(t + 128)
}
