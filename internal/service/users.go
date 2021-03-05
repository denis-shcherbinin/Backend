package service

import (
	"github.com/PolyProjectOPD/Backend/internal/domain"
	"github.com/PolyProjectOPD/Backend/internal/repository"
	"github.com/PolyProjectOPD/Backend/pkg/hash"
	"time"
)

type UsersService struct {
	repos  repository.Users
	hasher hash.PasswordHasher
}

func NewUsersService(repos repository.Users, hasher hash.PasswordHasher) *UsersService {
	return &UsersService{
		repos:  repos,
		hasher: hasher,
	}
}

func (u *UsersService) SignUp(input UserSignUpInput) (int, error) {
	user := domain.User{
		Name:         input.Name,
		Email:        input.Email,
		Password:     u.hasher.Hash(input.Password),
		RegisteredAt: time.Now(),
		LastVisitAt:  time.Now(),
	}

	id, err := u.repos.Create(user)
	if err != nil {
		return 0, err
	}

	return id, nil
}
