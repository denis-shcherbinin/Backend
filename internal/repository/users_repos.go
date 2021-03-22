package repository

import (
	"fmt"
	"github.com/PolyProjectOPD/Backend/internal/entity"
	"github.com/PolyProjectOPD/Backend/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
	"time"
)

type UsersRepos struct {
	db *sqlx.DB
}

func NewUsersRepos(db *sqlx.DB) *UsersRepos {
	return &UsersRepos{
		db: db,
	}
}

// Create adds a new user to the users table.
// It returns new user id and error.
func (u *UsersRepos) Create(user entity.User) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (name, email, password_hash, registered_at) values ($1, $2, $3, $4) RETURNING id",
		postgres.UsersTable)
	row := u.db.QueryRow(query, user.Name, user.Email, user.Password, user.RegisteredAt)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// GetByCredentials looks in the users table for the presence of a user with passed credentials(email, password).
// It returns user and error.
func (u *UsersRepos) GetByCredentials(email, password string) (entity.User, error) {
	var user entity.User

	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", postgres.UsersTable)
	err := u.db.Get(&user, query, email, password)

	return user, err
}

// GetIDByRefreshToken looks in the user sessions table for the presence of a user with passed refresh token.
// It returns user id and error.
func (u *UsersRepos) GetIDByRefreshToken(refreshToken string) (int, error) {

	var userID int

	query := fmt.Sprintf("SELECT user_id FROM %s WHERE refresh_token=$1 AND expires_at > $2",
		postgres.UsersSessionsTable)
	err := u.db.Get(&userID, query, refreshToken, time.Now())

	return userID, err
}

// CreateSession adds a new session for the user with the passed id.
// It returns an error.
func (u *UsersRepos) CreateSession(id int, session entity.Session) error {
	// todo: [SCN-41]: Проверка на допустимое количество сессий, удаление прошлых, сохранение последней

	query := fmt.Sprintf("INSERT INTO %s (user_id, refresh_token, expires_at) values ($1, $2, $3)",
		postgres.UsersSessionsTable)
	_, err := u.db.Exec(query, id, session.RefreshToken, session.ExpiresAt)
	return err
}

// UpdateSession updates an existing session for a user with a passed id and refresh token.
// It returns an error.
func (u *UsersRepos) UpdateSession(refreshToken string, session entity.Session) error {
	query := fmt.Sprintf("UPDATE %s SET refresh_token=$1, expires_at=$2 WHERE refresh_token=$3",
		postgres.UsersSessionsTable)
	_, err := u.db.Exec(query, session.RefreshToken, session.ExpiresAt, refreshToken)

	return err
}
