package usecase

import "go-web-app/domain"

type ICreateUser interface {
	Execute(req *CreateUserRequest) (*CreateUserResponse, error)
}

type CreateUserRequest struct {
	Name           string
	HashedPassword []byte
}

func (u *CreateUserRequest) ToUserModel() *domain.User {
	return domain.NewUser(u.Name, u.HashedPassword)
}

type CreateUserResponse struct {
	User *domain.User
}
