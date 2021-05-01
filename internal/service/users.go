package service

import (
	"github.com/PolyProjectOPD/Backend/internal/entity"
	"github.com/PolyProjectOPD/Backend/internal/repository"
	"github.com/PolyProjectOPD/Backend/internal/storage"
	"github.com/PolyProjectOPD/Backend/pkg/auth"
	"github.com/PolyProjectOPD/Backend/pkg/hash"
	"strconv"
	"time"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type UsersService struct {
	repos        repository.Users
	storage      *storage.Storage
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUsersService(repos repository.Users, storage *storage.Storage, hasher hash.PasswordHasher,
	tokenManager auth.TokenManager, accessTokenTTL, refreshTokenTTL time.Duration) *UsersService {
	return &UsersService{
		repos:           repos,
		storage:         storage,
		hasher:          hasher,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

// SignUp registers a new user and upload his image.
// It returns new user id and error.
func (u *UsersService) SignUp(input entity.UserSignUpInput, fileBody, fileType string) (int, string, error) {
	user := entity.User{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		BirthDate:    input.BirthDate,
		Email:        input.Email,
		Password:     u.hasher.Hash(input.Password),
		RegisteredAt: time.Now(),
		InSearch:     input.InSearch,
	}

	spheres := input.Spheres
	skills := input.Skills

	var (
		imageURL string
		err      error
	)

	if len(fileBody) != 0 && len(fileType) != 0 {
		imageURL, err = u.storage.Upload(storage.UploadInput{
			Body:        fileBody,
			ContentType: fileType,
		})
		if err != nil {
			return 0, "", err
		}
	}

	id, err := u.repos.Create(user, spheres, skills, imageURL)
	if err != nil {
		return 0, "", err
	}

	return id, imageURL, nil
}

// SignIn authenticates the user.
// It returns tokens(access and refresh) and error.
func (u *UsersService) SignIn(input entity.UserSignInInput, userAgent string) (Tokens, error) {
	user, err := u.repos.GetByCredentials(input.Email, u.hasher.Hash(input.Password))

	if err != nil {
		return Tokens{}, err
	}

	return u.createSession(user.ID, userAgent)
}

// RefreshTokens refreshes tokens for a user with passed refresh token.
// It returns tokens(access and refresh) and error.
func (u *UsersService) RefreshTokens(input entity.UserRefreshInput, userAgent string) (Tokens, error) {
	userID, err := u.repos.GetIDByRefreshToken(input.Token)
	if err != nil {
		return Tokens{}, err
	}

	return u.updateSession(userID, userAgent, input.Token)
}

// Profile gather user profile with passed user id.
// It returns user profile.
func (u *UsersService) Profile(userID int) (entity.UserProfile, error) {
	var userProfile entity.UserProfile

	// Getting user
	user, err := u.repos.GetByID(userID)
	if err != nil {
		return userProfile, err
	}
	userProfile.FirstName = user.FirstName
	userProfile.LastName = user.LastName
	userProfile.Email = user.Email
	userProfile.ImageURL = user.ImageURL

	// Age calculating
	l, _ := time.LoadLocation("Local")
	day, _ := strconv.Atoi(user.BirthDate[:2])
	month, _ := strconv.Atoi(user.BirthDate[3:5])
	year, _ := strconv.Atoi(user.BirthDate[6:10])
	userAge := strconv.Itoa(int(time.Now().Sub(time.Date(year, time.Month(month), day, 0, 0, 0, 0, l)).Hours() / 24 / 365))
	userProfile.Age = userAge

	// Getting user profile info
	userProfileInfo, err := u.repos.GetProfileInfo(userID)
	if err != nil {
		return userProfile, err
	}
	userProfile.Comment = userProfileInfo[0]
	userProfile.Experience = userProfileInfo[1]
	userProfile.SkillLevel = userProfileInfo[2]
	userProfile.MinSalary = userProfileInfo[3]
	userProfile.MaxSalary = userProfileInfo[4]
	userProfile.About = userProfileInfo[5]

	// Getting user skills
	userSkills, err := u.repos.GetSkills(userID)
	if err != nil {
		return userProfile, err
	}
	userProfile.Skills = userSkills

	// Getting user jobs
	userJobs, err := u.repos.GetJobs(userID)
	if err != nil {
		return userProfile, err
	}
	userProfile.Jobs = userJobs

	return userProfile, nil
}

// Logout deletes all active sessions the user with passer id.
// It returns an error.
func (u *UsersService) Logout(userID int) error {
	return u.repos.DeleteAllSessions(userID)
}

// SignOut deletes all sessions with passed userID and userAgent.
// It return an error.
func (u *UsersService) SignOut(userID int, userAgent string) error {
	return u.repos.DeleteAllAgentSessions(userID, userAgent)
}

// Existence checks for the existence of a user with passed email.
// It returns true if user exists otherwise false.
func (u *UsersService) Existence(input entity.UserExistenceInput) bool {
	return u.repos.Existence(input.Email)
}

// generateTokens generates a new pair of tokens.
// It returns tokens(access and refresh) and error.
func (u *UsersService) generateTokens(id int) (Tokens, error) {
	var (
		tokens Tokens
		err    error
	)

	tokens.AccessToken, err = u.tokenManager.NewJWT(strconv.Itoa(id), u.accessTokenTTL)
	if err != nil {
		return tokens, err
	}

	tokens.RefreshToken, err = u.tokenManager.NewRefreshToken()

	return tokens, err
}

// createSession creates a new session for the user with passed id and userAgent.
// It returns tokens(access and refresh) and error.
func (u *UsersService) createSession(id int, userAgent string) (Tokens, error) {
	var (
		tokens Tokens
		err    error
	)

	tokens, err = u.generateTokens(id)
	if err != nil {
		return tokens, err
	}

	err = u.repos.CreateSession(id, entity.Session{
		UserAgent:    userAgent,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    time.Now().Add(u.refreshTokenTTL),
	})

	return tokens, err
}

// updateSession updates an existing user session with the passed id and refresh token.
// It returns tokens(access and refresh) and error.
func (u *UsersService) updateSession(id int, userAgent string, refreshToken string) (Tokens, error) {
	var (
		tokens Tokens
		err    error
	)

	tokens, err = u.generateTokens(id)
	if err != nil {
		return tokens, err
	}

	err = u.repos.UpdateSession(id, refreshToken, entity.Session{
		UserAgent:    userAgent,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    time.Now().Add(u.refreshTokenTTL),
	})

	return tokens, err
}
