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
		t.FailNow()
	}

	registerCode := s.Register(0, username, password)
	if registerCode != message.RegisterResponseCodes["ok"] {
		t.FailNow()
	}

	registerCode = s.Register(0, username, "otherpass")
	if registerCode != message.RegisterResponseCodes["used_pseudo"] {
		t.FailNow()
	}

	loginCode = s.Login(0, username, "otherpass")
	if loginCode != message.LoginResponseCodes["wrong_password"] {
		t.FailNow()
	}

	loginCode = s.Login(0, username, password)
	if loginCode != message.LoginResponseCodes["ok"] {
		t.FailNow()
	}

	session := s.GetSession(0)
	if session == nil {
		t.FailNow()
	}

	s.Logout(0)

	session = s.GetSession(0)
	if session != nil {
		t.FailNow()
	}
}
