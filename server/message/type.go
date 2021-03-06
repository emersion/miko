// Provides the protocol core data types.
package message

// A protocol version.
type ProtocolVersion uint16

// The current protocol version.
const CurrentVersion ProtocolVersion = 10

type (
	Type uint8 // A message type.
	ExitCode uint8
	LoginResponseCode uint8
	RegisterResponseCode uint8
	MetaActionCode uint8
	VersionResponseCode uint8
)

type PointType uint8

// A map containing all message types.
var Types = map[string]Type{
	"ping":              0,
	"pong":              1,
	"exit":              2,
	"login":             3,
	"login_response":    4,
	"register":          5,
	"register_response": 6,
	"meta_action":       7,
	"chunk_update":      8,
	"terrain_request":   9,
	"entities_update":   10,
	"entity_update":     11,
	"actions_done":      12,
	"action_do":         13,
	"entity_create":     14,
	"entity_destroy":    15,
	"chat_send":         16,
	"chat_receive":      17,
	"version":           18,
	"config":            19,
	"entity_id_change":  20,
	"chunks_update":     21,
}

// Get a message type name from its type.
func GetTypeName(t Type) string {
	for name, val := range Types {
		if t == val {
			return name
		}
	}
	return ""
}

var ExitCodes = map[string]ExitCode{
	"client_quit":     0,
	"server_closed":   1,
	"network_error":   2,
	"ping_timeout":    3,
	"client_kicked":   4,
	"client_banned":   5,
	"client_outdated": 6,
	"server_outdated": 7,
}

var LoginResponseCodes = map[string]LoginResponseCode{
	"ok":                   0,
	"unknown_pseudo":       1,
	"wrong_password":       2,
	"too_many_tries":       3,
	"already_connected":    4,
	"player_limit_reached": 5,
}

var RegisterResponseCodes = map[string]RegisterResponseCode{
	"ok":                0,
	"used_pseudo":       1,
	"invalid_pseudo":    2,
	"invalid_password":  2,
	"too_many_tries":    3,
	"register_disabled": 4,
}

var MetaActionCodes = map[string]MetaActionCode{
	"player_joined": 0,
	"player_left":   1,
}
