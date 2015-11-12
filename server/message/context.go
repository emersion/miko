package message

type contextType int

const (
	serverContext contextType = iota
	clientContext
)

// A context contains all services that handles backend functionalities
// TODO: properly separate ClientContext and ServerContext
type Context struct {
	contextType

	Config *Config

	Entity  EntityService
	Terrain Terrain
	Clock   ClockService

	// Server
	Auth AuthService

	// Client
	Me *Session
}

func (c *Context) IsServer() bool {
	return c.contextType == serverContext
}
func (c *Context) IsClient() bool {
	return c.contextType == clientContext
}

func NewServerContext() *Context {
	return &Context{contextType: serverContext}
}
func NewClientContext() *Context {
	return &Context{contextType: clientContext}
}
