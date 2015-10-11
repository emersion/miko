package auth

import(
	"git.emersion.fr/saucisse-royale/miko/server/message"
)

type AuthService struct {
	sessions []*message.Session
	users []*User
}

func (a *AuthService) HasSession(id int) bool {
	session := a.GetSession(id)
	return (session != nil)
}

func (a *AuthService) GetSession(id int) *message.Session {
	// TODO: make sure session ID is io ID
	if id >= len(a.sessions) {
		return nil
	}
	return a.sessions[id]
}

func (a *AuthService) Login(io *message.IO, username string, password string) message.LoginResponseCode {
	var code string
	for _, user := range a.users {
		if username != user.Username {
			continue
		}

		if user.VerifyPassword(password) {
			code = "ok"
		} else {
			code = "wrong_password"
		}
	}

	if code == "ok" {
		session := &message.Session{
			Id: io.Id,
			Username: username,
		}
		a.sessions = append(a.sessions, session)
	} else if code == "" {
		code = "unknown_pseudo"
	}

	return message.LoginResponseCodes[code]
}

func (a *AuthService) Logout(io *message.IO) {
	if !a.HasSession(io.Id) {
		return
	}

	a.sessions[io.Id] = nil
}

func (a *AuthService) Register(io *message.IO, username string, password string) message.RegisterResponseCode {
	var code string
	for _, user := range a.users {
		if username == user.Username {
			code = "used_pseudo"
			break
		}
	}

	if code != "" {
		code = "ok"
		a.users = append(a.users, &User{username, password})
	}

	return message.RegisterResponseCodes[code]
}

func NewService() *AuthService {
	return &AuthService{
		users: LoadUserDb(),
	}
}
