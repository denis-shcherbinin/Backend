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

// SignUp registers a new user
// It returns new user id and error
func (u *UsersService) SignUp(input UserSignUpInput) (int, error) {
	user := entity.User{
		Name:         input.Name,
		Email:        input.Email,
		Password:     u.hasher.Hash(input.Password),
		RegisteredAt: time.Now(),
	}

	id, err := u.repos.Create(user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SignIn authenticates the user
// It returns tokens(access and refresh) and error
func (u *UsersService) SignIn(input UserSignInInput) (Tokens, error) {
	user, err := u.repos.GetByCredentials(input.Email, u.hasher.Hash(input.Password))

	if err != nil {
		return Tokens{}, err
	}

	return u.createSession(user.ID)
}

// RefreshTokens refreshes tokens for a user with passed refresh token
// It returns tokens(access and refresh) and error
func (u *UsersService) RefreshTokens(refreshToken string) (Tokens, error) {
	userID, err := u.repos.GetIDByRefreshToken(refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return u.updateSession(userID, refreshToken)
}

// generateTokens generates a new pair of tokens
// It returns tokens(access and refresh) and error
func (u *UsersService) generateTokens(id int) (Tokens, error) {
	var (
		tokens Tokens
		err    error
	)

	tokens.AccessToken, err = u.tokenManager.NewJWT(strconv.FormatInt(int64(id), 16), u.accessTokenTTL)
	if err != nil {
		return tokens, err
	}

	tokens.RefreshToken, err = u.tokenManager.NewRefreshToken()

	return tokens, err
}

// createSession creates a new session for the user with passed id
// It returns tokens(access and refresh) and error
func (u *UsersService) createSession(id int) (Tokens, error) {
	var (
		tokens Tokens
		err    error
	)

	tokens, err = u.generateTokens(id)
	if err != nil {
		return tokens, err
	}

	err = u.repos.CreateSession(id, entity.Session{
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    time.Now().Add(u.refreshTokenTTL),
	})

	return tokens, err
}

// updateSession updates an existing user session with the passed id and refresh token
// It returns tokens(access and refresh) and error
func (u *UsersService) updateSession(id int, refreshToken string) (Tokens, error) {
	var (
		tokens Tokens
		err    error
	)

	tokens, err = u.generateTokens(id)
	if err != nil {
		return tokens, err
	}

	err = u.repos.UpdateSession(id, refreshToken, entity.Session{
		RefreshToken: tokens.RefreshToken,
		ExpiresAt: time.Now().Add(u.refreshTokenTTL),
	})

	return tokens, err
}
