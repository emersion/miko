// Provides functions to manage users authentication.
package auth

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

// The authentication service
// It aims to manage users: login, register, sessions
type Service struct {
	sessions      map[int]*message.Session
	users         []*User
	LoginCallback func(*message.Session)
}

func (s *Service) HasSession(id int) bool {
	_, ok := s.sessions[id]
	return ok
}

func (s *Service) GetSession(id int) *message.Session {
	if session, ok := s.sessions[id]; ok {
		return session
	}
	return nil
}

func (s *Service) getSessionByUsername(username string) *message.Session {
	for _, session := range s.sessions {
		if session != nil && session.Username == username {
			return session
		}
	}
	return nil
}

func (s *Service) GetSessionByEntity(id message.EntityId) *message.Session {
	for _, session := range s.sessions {
		if session != nil && session.Entity != nil && session.Entity.Id == id {
			return session
		}
	}
	return nil
}

func (s *Service) List() []*message.Session {
	l := []*message.Session{}

	for _, session := range s.sessions {
		if session != nil {
			l = append(l, session)
		}
	}

	return l
}

func (s *Service) Login(id int, username string, password string) message.LoginResponseCode {
	code := "unknown_pseudo"
	for _, user := range s.users {
		if username != user.Username {
			continue
		}

		session := s.getSessionByUsername(username)
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
		session := &message.Session{
			Id:       id,
			Username: username,
			Entity:   message.NewEntity(),
		}

		if s.LoginCallback != nil {
			s.LoginCallback(session)
		}

		s.sessions[id] = session
	}

	return message.LoginResponseCodes[code]
}

func (s *Service) Logout(id int) {
	if !s.HasSession(id) {
		return
	}

	delete(s.sessions, id)
}

func (s *Service) Register(id int, username string, password string) message.RegisterResponseCode {
	for _, user := range s.users {
		if username == user.Username {
			return message.RegisterResponseCodes["used_pseudo"]
		}
	}

	hash, _ := hashPassword(password)
	s.users = append(s.users, &User{username, hash})

	return message.RegisterResponseCodes["ok"]
}

func NewService() *Service {
	return &Service{
		sessions: map[int]*message.Session{},
		users:    LoadUserDb(),
	}
}
