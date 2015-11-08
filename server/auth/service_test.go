package auth_test

import (
	"git.emersion.fr/saucisse-royale/miko.git/server/auth"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"testing"
)

const (
	username = "username"
	password = "password"
)

func TestService(t *testing.T) {
	s := auth.NewService()

	loginCode := s.Login(0, username, password)
	if loginCode != message.LoginResponseCodes["unknown_pseudo"] {
		t.Error("Bad response code after wrong login attempt:", loginCode)
	}

	registerCode := s.Register(0, username, password)
	if registerCode != message.RegisterResponseCodes["ok"] {
		t.Error("Bad register response code:", registerCode)
	}

	registerCode = s.Register(0, username, "otherpass")
	if registerCode != message.RegisterResponseCodes["used_pseudo"] {
		t.Error("Bad response code after wrong register attempt:", registerCode)
	}

	loginCode = s.Login(0, username, "otherpass")
	if loginCode != message.LoginResponseCodes["wrong_password"] {
		t.Error("Bad response code after wrong login attempt:", loginCode)
	}

	loginCode = s.Login(0, username, password)
	if loginCode != message.LoginResponseCodes["ok"] {
		t.Error("Bad login response code:", loginCode)
	}

	session := s.GetSession(0)
	if session == nil {
		t.Error("Could not get session after login")
	}

	s.Logout(0)

	session = s.GetSession(0)
	if session != nil {
		t.Error("Session not flushed after logout")
	}
}
