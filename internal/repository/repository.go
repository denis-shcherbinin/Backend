package repository

import (
	"github.com/PolyProjectOPD/Backend/internal/entity"
	"github.com/jmoiron/sqlx"
)

type Users interface {
	Create(user entity.User, spheres []entity.Sphere, skills []entity.Skill, imageURL string) (int, error)

	GetByCredentials(email, password string) (entity.User, error)
	GetByID(id int) (entity.User, error)
	GetIDByRefreshToken(refreshToken string) (int, error)
	GetProfileInfo(id int) ([]string, error)
	GetSkills(id int) ([]entity.Skill, error)
	GetJobs(id int) ([]entity.Job, error)

	DeleteAllSessions(id int) error
	DeleteAllAgentSessions(id int, userAgent string) error

	CreateSession(id int, session entity.Session) error
	UpdateSession(id int, refreshToken string, session entity.Session) error

	Existence(email string) bool
}

type Spheres interface {
	GetAll() ([]entity.Sphere, error)
	GetSkills(sphere entity.Sphere) ([]entity.Skill, error)
}

type Skills interface {
	GetAll() ([]entity.Skill, error)
}

type Companies interface {
	Create(userID int, company entity.Company) (int, error)
}

type Repositories struct {
	Users     Users
	Spheres   Spheres
	Skills    Skills
	Companies Companies
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Users:     NewUsersRepos(db),
		Spheres:   NewSpheresRepos(db),
		Skills:    NewSkillsRepos(db),
		Companies: NewCompaniesRepos(db),
	}
}
