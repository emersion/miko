package message

type ContextType int

const (
	ServerContext ContextType = 1
	ClientContext ContextType = 2
)

type Context struct {
	Type ContextType
	Auth AuthService
	Terrain Terrain
}
