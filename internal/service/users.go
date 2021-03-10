package service

import (
	"github.com/PolyProjectOPD/Backend/internal/entity"
	"github.com/PolyProjectOPD/Backend/internal/repository"
	"github.com/PolyProjectOPD/Backend/pkg/auth"
	"github.com/PolyProjectOPD/Backend/pkg/hash"
	"strconv"
	"time"
)

type UserSignUpInput struct {
	Name     string
	Email    string
	Password string
}

type UserSignInInput struct {
	Email    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type UsersService struct {
	repos        repository.Users
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUsersService(repos repository.Users, hasher hash.PasswordHasher, tokenManager auth.TokenManager, accessTokenTTL, refreshTokenTTL time.Duration) *UsersService {
	return &UsersService{
		repos:           repos,
		hasher:          hasher,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (u *UsersService) SignUp(input UserSignUpInput) (int, error) {
	user := entity.User{
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

func (u *UsersService) SignIn(input UserSignInInput) (Tokens, error) {
	user, err := u.repos.GetByCredentials(input.Email, u.hasher.Hash(input.Password))

	if err != nil {
		return Tokens{}, err
	}

	return u.createSession(user.ID)
}

func (u *UsersService) RefreshTokens(refreshToken string) (Tokens, error) {
	userID, err := u.repos.GetIDByRefreshToken(refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return u.createSession(userID)
}

func (u *UsersService) createSession(id int) (Tokens, error) {
	var (
		tokens Tokens
		err    error
	)

	tokens.AccessToken, err = u.tokenManager.NewJWT(strconv.FormatInt(int64(id), 16), u.accessTokenTTL)
	if err != nil {
		return tokens, err
	}

	tokens.RefreshToken, err = u.tokenManager.NewRefreshToken()
	if err != nil {
		return tokens, err
	}

	err = u.repos.CreateSession(id, entity.Session{
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    time.Now().Add(u.refreshTokenTTL),
	})

	return tokens, err
}
