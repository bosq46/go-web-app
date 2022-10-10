package interactor

import (
	"go-web-app/domain"
	"go-web-app/usecase"

	"github.com/pkg/errors"
)

// usecace の ICreateUser インターフェイスの実装
type UserCreator struct {
	UserRepository domain.UserRepository
}

func NewCreateUser(repos domain.UserRepository) *UserCreator {
	return &UserCreator{
		UserRepository: repos,
	}
}

func (u *UserCreator) Execute(req *usecase.CreateUserRequest) (*usecase.CreateUserResponse, error) {
	user, err := u.UserRepository.CreateUser(req.ToUserModel())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &usecase.CreateUserResponse{User: user}, nil
}
