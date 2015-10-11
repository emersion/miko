package message

type Session struct {
	Id int
	Username string
}

type AuthService interface {
	GetSession(id int) *Session
	HasSession(id int) bool
	Login(id int, username string, password string) LoginResponseCode
	Logout(id int)
	Register(id int, username string, password string) RegisterResponseCode
}
