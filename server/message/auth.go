package message

type AuthService interface {
	Login(io *IO, username string, password string) LoginResponseCode
	Logout(io *IO)
	Register(io *IO, username string, password string) RegisterResponseCode
}
