package auth

type Store interface {
	CreateUser(body CreateUserBody) error
	FindUserByEmail(email string) (User, error)
	ComparePasswords(password, candidatePassword string) bool
}
