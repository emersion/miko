package message

type Session struct {
	Id int
	Username string
}

type AuthService interface {
	GetSession(id int) *Session
	Login(io *IO, username string, password string) LoginResponseCode
	Logout(io *IO)
	Register(io *IO, username string, password string) RegisterResponseCode
}