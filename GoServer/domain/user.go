package domain

type User struct {
	ID             uint64
	Name           string
	HashedPassword []byte
}

func NewUser(name string, hashedPassword []byte) *User {
	return &User{Name: name, HashedPassword: hashedPassword}
}

type UserRepository interface {
	GetUsers() ([]*User, error)
	GetUserByID(id uint64) (*User, error)
	CreateUser(newUser *User) (*User, error)
	UpdateUser(newUser *User) error
	DeleteUser(targetUser *User) error
}
