package repository

import (
	"errors"
	"fmt"
	"github.com/PolyProjectOPD/Backend/internal/entity"
	"github.com/PolyProjectOPD/Backend/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
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

// Create adds a new user to the users table, fills tables of user spheres and skills.
// It returns new user id and error.
func (u *UsersRepos) Create(user entity.User, spheres []entity.Sphere, skills []entity.Skill) (int, error) {
	var userID int

	query := fmt.Sprintf("INSERT INTO %s "+
		"(first_name, last_name, birth_date, email, password_hash, in_search, registered_at) "+
		"values ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		postgres.UsersTable)

	row := u.db.QueryRow(query, user.FirstName, user.LastName, user.BirthDate, user.Email, user.Password, user.InSearch, user.RegisteredAt)

	if err := row.Scan(&userID); err != nil {
		return 0, errors.New("user already exists")
	}

	for _, sphere := range spheres {
		query = fmt.Sprintf("INSERT INTO %s (user_id, sphere_id) values ($1, $2)", postgres.UsersSpheresTable)
		_, err := u.db.Exec(query, userID, sphere.ID)
		if err != nil {
			logrus.Error(err)
			continue
		}
	}

	for _, skill := range skills {
		query = fmt.Sprintf("INSERT INTO %s (user_id, skill_id) values ($1, $2)", postgres.UsersSkillsTable)
		_, err := u.db.Exec(query, userID, skill.ID)
		if err != nil {
			logrus.Error(err)
			continue
		}
	}

	return userID, nil
}

// GetByCredentials looks in the users table for the presence of a user with passed credentials(email, password).
// It returns user and error.
func (u *UsersRepos) GetByCredentials(email, password string) (entity.User, error) {
	var user entity.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1 AND password_hash=$2", postgres.UsersTable)

	raw, err := u.db.Query(query, email, password)
	if err != nil {
		return user, errors.New("invalid email or password")
	}

	raw.Next()
	err = raw.Scan(&user.ID, &user.FirstName, &user.LastName, &user.BirthDate,
		&user.Email, &user.Password, &user.InSearch, &user.RegisteredAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

// GetIDByRefreshToken looks in the user sessions table for the presence of a user with passed refresh token.
// It returns user id and error.
func (u *UsersRepos) GetIDByRefreshToken(refreshToken string) (int, error) {

	var userID int

	query := fmt.Sprintf("SELECT user_id FROM %s WHERE refresh_token=$1 AND expires_at > $2",
		postgres.UsersSessionsTable)
	if err := u.db.Get(&userID, query, refreshToken, time.Now()); err != nil {
		return 0, errors.New("invalid refresh token")
	}

	return userID, nil
}

// DeleteSessions removes all user sessions from user sessions table with passed user id.
// It returns an error.
func (u *UsersRepos) DeleteSessions(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1", postgres.UsersSessionsTable)
	_, err := u.db.Exec(query, id)
	return err
}

// CreateSession adds a new session for the user with the passed id.
// If the number of existing sessions is exceeded, then all of them are deleted, the last one is saved.
// It returns an error.
func (u *UsersRepos) CreateSession(id int, session entity.Session) error {
	var sessionsAmount []int
	query := fmt.Sprintf("SELECT COUNT(user_id) as amount FROM %s WHERE user_id=$1", postgres.UsersSessionsTable)
	if err := u.db.Select(&sessionsAmount, query, id); err != nil {
		return err
	}
	if sessionsAmount[0] == postgres.UsersMaxSessionsAmount {
		if err := u.DeleteSessions(id); err != nil {
			return err
		}
	}

	query = fmt.Sprintf("INSERT INTO %s (user_id, refresh_token, expires_at) values ($1, $2, $3)",
		postgres.UsersSessionsTable)
	_, err := u.db.Exec(query, id, session.RefreshToken, session.ExpiresAt)
	return err
}

// UpdateSession updates an existing session for a user with the passed refresh token.
// It returns an error.
func (u *UsersRepos) UpdateSession(refreshToken string, session entity.Session) error {
	query := fmt.Sprintf("UPDATE %s SET refresh_token=$1, expires_at=$2 WHERE refresh_token=$3",
		postgres.UsersSessionsTable)
	_, err := u.db.Exec(query, session.RefreshToken, session.ExpiresAt, refreshToken)

	return err
}

// Existence checks for the existence of a user with passed email.
// It returns true if user exists otherwise false.
func (u *UsersRepos) Existence(email string) bool {
	var amount []int
	query := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE email=$1", postgres.UsersTable)
	if err := u.db.Select(&amount, query, email); err != nil {
		logrus.Error(err)
	}

	return amount[0] == 1
}
