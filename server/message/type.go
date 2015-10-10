package message

type Type uint8
type LoginResponseCode uint8
type RegisterResponseCode uint8
type MetaActionCode uint8
type PointType uint8

var Types = map[string]Type{
	"ping": 0,
	"pong": 1,
	"exit": 2,
	"login": 3,
	"loginresponse": 4,
	"register": 5,
	"registerresponse": 6,
	"playermeta": 7,
	"terrainupdate": 8,
	"terrainrequest": 9,
	"entitiesupdate": 10,
	"entityupdate": 11,
	"actions": 12,
	"action": 13,
	"entitycreate": 14,
	"entitiydestroy": 15,
	"chatsend": 16,
	"chatreceive": 17,
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
