package repository

import (
	"github.com/PolyProjectOPD/Backend/internal/entity"
	"github.com/jmoiron/sqlx"
)

type Users interface {
	Create(user entity.User, spheres []entity.Sphere, skills []entity.Skill) (int, error)
	GetByCredentials(email, password string) (entity.User, error)
	GetIDByRefreshToken(refreshToken string) (int, error)

	DeleteSessions(id int) error
	CreateSession(id int, session entity.Session) error
	UpdateSession(refreshToken string, session entity.Session) error
	Existence(email string) bool
}

type Spheres interface {
	Get() ([]entity.Sphere, error)
	GetSkills(sphere entity.Sphere) ([]entity.Skill, error)
}

type Repositories struct {
	Users   Users
	Spheres Spheres
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Users:   NewUsersRepos(db),
		Spheres: NewSpheresRepos(db),
	}
}
