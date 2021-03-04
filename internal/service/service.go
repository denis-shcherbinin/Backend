package service

import (
	"github.com/PolyProjectOPD/Backend/internal/repository"
	"github.com/PolyProjectOPD/Backend/pkg/hash"
)

type UserSignUpInput struct {
	Name     string
	Email    string
	Password string
}

type Users interface {
	SignUp(input UserSignUpInput) (int, error)
}

type Services struct {
	Users Users
}

type Deps struct {
	Repos  *repository.Repositories
	Hasher hash.PasswordHasher
}

func NewServices(deps Deps) *Services {
	return &Services{
		Users: NewUsersService(deps.Repos.Users, deps.Hasher),
	}
}
