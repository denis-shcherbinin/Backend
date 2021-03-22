package service

import (
	"github.com/PolyProjectOPD/Backend/internal/repository"
	"github.com/PolyProjectOPD/Backend/pkg/auth"
	"github.com/PolyProjectOPD/Backend/pkg/hash"
	"time"
)

type Users interface {
	SignUp(input UserSignUpInput) (int, error)
	SignIn(input UserSignInInput) (Tokens, error)

	RefreshTokens(refreshToken string) (Tokens, error)
}

type Services struct {
	Users Users
}

type Deps struct {
	Repos           *repository.Repositories
	Hasher          hash.PasswordHasher
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewServices(deps Deps) *Services {
	usersService := NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL)

	return &Services{
		Users: usersService,
	}
}
