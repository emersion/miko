package message

// A context type is either client or server
type ContextType int

const (
	ServerContext ContextType = 1
	ClientContext ContextType = 2
)

// A context contains all services that handles backend functionalities
type Context struct {
	Type ContextType
	Auth AuthService
	Entity EntityService
	Terrain Terrain
}
