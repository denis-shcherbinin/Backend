package repository

import (
	"fmt"
	"github.com/PolyProjectOPD/Backend/internal/entity"
	"github.com/PolyProjectOPD/Backend/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

type UsersRepos struct {
	db *sqlx.DB
}

func NewUsersRepos(db *sqlx.DB) *UsersRepos {
	return &UsersRepos{
		db: db,
	}
}

func (u *UsersRepos) Create(user entity.User) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (name, email, password_hash, registered_at, last_visit_at) values ($1, $2, $3, $4, $5) RETURNING id", postgres.UsersTable)
	row := u.db.QueryRow(query, user.Name, user.Email, user.Password, user.RegisteredAt, user.LastVisitAt)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
