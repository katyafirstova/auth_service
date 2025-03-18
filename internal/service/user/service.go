package user

import (
	"github.com/katyafirstova/auth_service/internal/repository"
	"github.com/katyafirstova/auth_service/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
}

func NewService(
	userRepository repository.UserRepository,
) *serv {
	return &serv{
		userRepository: userRepository,
	}
}

func NewMockService(deps ...interface{}) service.UserService {
	srv := serv{}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.UserRepository:
			srv.userRepository = s
		}
	}

	return &srv
}
