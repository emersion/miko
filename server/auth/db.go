package auth

type User struct {
	Username string
	password string
}

func (u *User) CheckPassword(password string) bool {
	return (u.password == password)
}

func LoadUserDb() []*User {
	return []*User{
		&User{"root", "root"},
	}
}
