// Provides functions to manage users authentication.
package auth

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

// The authentication service
// It aims to manage users: login, register, sessions
type Service struct {
	sessions map[int]*message.Session
	users    []*User
}

func (a *Service) HasSession(id int) bool {
	_, ok := a.sessions[id]
	return ok
}

func (a *Service) GetSession(id int) *message.Session {
	if session, ok := a.sessions[id]; ok {
		return session
	}
	return nil
}

func (a *Service) getSessionByUsername(username string) *message.Session {
	for _, session := range a.sessions {
		if session != nil && session.Username == username {
			return session
		}
	}
	return nil
}

func (a *Service) GetSessionByEntity(id message.EntityId) *message.Session {
	for _, session := range a.sessions {
		if session != nil && session.Entity != nil && session.Entity.Id == id {
			return session
		}
	}
	return nil
}

func (a *Service) Login(id int, username string, password string) message.LoginResponseCode {
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
		a.sessions[id] = &message.Session{
			Id:       id,
			Username: username,
			Entity:   entity,
		}
	}

	return message.LoginResponseCodes[code]
}

func (a *Service) Logout(id int) {
	if !a.HasSession(id) {
		return
	}

	delete(a.sessions, id)
}

func (a *Service) Register(id int, username string, password string) message.RegisterResponseCode {
	for _, user := range a.users {
		if username == user.Username {
			return message.RegisterResponseCodes["used_pseudo"]
		}
	}

	hash, _ := hashPassword(password)
	a.users = append(a.users, &User{username, hash})

	return message.RegisterResponseCodes["ok"]
}

func NewService() *Service {
	return &Service{
		sessions: map[int]*message.Session{},
		users:    LoadUserDb(),
	}
}
