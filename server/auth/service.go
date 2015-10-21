package auth

import(
	"git.emersion.fr/saucisse-royale/miko/server/message"
)

// The authentication service
// It aims to manage users: login, register, sessions
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

func (a *AuthService) getSessionByUsername(username string) *message.Session {
	for _, session := range a.sessions {
		if session != nil && session.Username == username {
			return session
		}
	}
	return nil
}

func (a *AuthService) Login(id int, username string, password string) message.LoginResponseCode {
	code := "unknown_pseudo"
	for _, user := range a.users {
		if username != user.Username {
			continue
		}

		session := a.getSessionByUsername(username)
		if session != nil {
			code = "already_connected"
			break
		}

		if user.VerifyPassword(password) {
			code = "ok"
		} else {
			code = "wrong_password"
		}
		break
	}

	if code == "ok" {
		entity := message.NewEntity()
		session := &message.Session{
			Id: id,
			Username: username,
			Entity: entity,
		}
		a.sessions = append(a.sessions, session)
	}

	return message.LoginResponseCodes[code]
}

func (a *AuthService) Logout(id int) {
	if !a.HasSession(id) {
		return
	}

	a.sessions[id] = nil
}

func (a *AuthService) Register(id int, username string, password string) message.RegisterResponseCode {
	for _, user := range a.users {
		if username == user.Username {
			return message.RegisterResponseCodes["used_pseudo"]
		}
	}

	hash, _ := hashPassword(password)
	a.users = append(a.users, &User{username, hash})

	return message.RegisterResponseCodes["ok"]
}

func NewService() *AuthService {
	return &AuthService{
		users: LoadUserDb(),
	}
}
