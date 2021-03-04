package repository

import (
	"github.com/PolyProjectOPD/Backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Users interface {
	Create(user domain.User) (int, error)
}

type Repositories struct {
	Users Users
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Users: NewUsersRepos(db),
	}
}
