package auth

// A user in the database
type User struct {
	Username string
	password string
}

func (u *User) VerifyPassword(password string) bool {
	if err := verifyPassword(u.password, password); err != nil {
		return false
	}
	return true
}

func LoadUserDb() []*User {
	hash, _ := hashPassword("root")
	return []*User{
		&User{"root", hash},
	}
}
