package auth

import(
	"../message"
)

type AuthService struct {
	sessions []*message.Session
}

func (a *AuthService) GetSession(id int) *message.Session {
	// TODO: make sure session ID is io ID
	return a.sessions[id]
}

func (a *AuthService) Login(io *message.IO, username string, password string) message.LoginResponseCode {
	var code string
	if username == "root" && password == "root" {
		code = "ok"

		session := &message.Session{
			Id: io.Id,
			Username: username,
		}
		a.sessions = append(a.sessions, session)
	} else {
		code = "unknownpseudo"
	}
	return message.LoginResponseCodes[code]
}

func (a *AuthService) Logout(io *message.IO) {}

func (a *AuthService) Register(io *message.IO, username string, password string) message.RegisterResponseCode {
	return message.RegisterResponseCodes["registerdisabled"]
}

func NewService() *AuthService {
	return &AuthService{}
}
