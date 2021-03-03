package service

import (
	"crypto/sha1"
	"fmt"
	"github.com/PolyProjectOPD/Backend/internal/domain"
	"github.com/PolyProjectOPD/Backend/internal/repository"
)

const salt = "rer32r2frg1yu2kyd8fd"

type AuthService struct {
	repos repository.Authorization
}

func NewAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repos: repos}
}

func (a *AuthService) CreateUser(user domain.User) (int, error) {
	user.Password = a.generatePasswordHash(user.Password)
	return a.repos.CreateUser(user)
}

func (a *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
