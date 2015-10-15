package message

// A session
// This maps connected clients to usernames and entities.
type Session struct {
	Id int
	Username string
	Entity *Entity
}

// An authentication service
type AuthService interface {
	GetSession(id int) *Session
	HasSession(id int) bool
	Login(id int, username string, password string) LoginResponseCode
	Logout(id int)
	Register(id int, username string, password string) RegisterResponseCode
}
