package auth

type Store interface {
	CreateUser(body CreateUserBody) error
}
