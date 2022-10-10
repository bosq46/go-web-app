package registry

import (
	"go-web-app/domain"
	"go-web-app/interactor"
	"go-web-app/usecase"
)

var FactorySingleton *Factory

type Factory struct {
	cache map[string]interface{}
}

func ClearFactory() {
	FactorySingleton = nil
}

func GetFactory() *Factory {
	if FactorySingleton == nil {
		FactorySingleton = &Factory{}
	}
	return FactorySingleton
}

func (f *Factory) container(key string, builder func() interface{}) interface{} {
	if f.cache == nil {
		f.cache = map[string]interface{}{}
	}
	if f.cache[key] == nil {
		f.cache[key] = builder()
	}
	return f.cache[key]
}
func (f *Factory) BuildUserOperator() domain.UserRepository {
	return f.container("UserOperator", func() interface{} {
		return &adapter.UserOperator{
			Client:                 f.BuildResourceTableOperator(),
			Mapper:                 f.BuildDynamoModelMapper(),
			UserEmailUniqGenerator: f.BuildUserEmailUniqGenerator(),
		}
	}).(domain.UserRepository)
}

func (f *Factory) BuildCreateUser() usecase.ICreateUser {
	return f.container("CreateUser", func() interface{} {
		return interactor.NewCreateUser(
			f.BuildUserOperator(),
		)
	}).(usecase.ICreateUser)
}
