package server

import(
	"../message"
)

type AuthService struct {}

func (a *AuthService) Login(io *message.IO, username string, password string) message.LoginResponseCode {
	var code string
	if username == "root" && password == "root" {
		code = "ok"
	} else {
		code = "unknownpseudo"
	}
	return message.LoginResponseCodes[code]
}

func (a *AuthService) Logout(io *message.IO) {}

func (a *AuthService) Register(io *message.IO, username string, password string) message.RegisterResponseCode {
	return message.RegisterResponseCodes["registerdisabled"]
}

func NewAuthService() *AuthService {
	return &AuthService{}
}
