package service

import (
	"github.com/PolyProjectOPD/Backend/internal/domain"
	"github.com/PolyProjectOPD/Backend/internal/repository"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
