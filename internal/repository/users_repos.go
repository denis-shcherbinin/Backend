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
func (u *UsersRepos) Create(user entity.User, spheres []entity.Sphere, skills []entity.Skill, imageURL string) (int, error) {
	var userID int

	query := fmt.Sprintf("INSERT INTO %s "+
		"(first_name, last_name, birth_date, email, password_hash, in_search, registered_at, image_url) "+
		"values ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		postgres.UsersTable)

	row := u.db.QueryRow(query, user.FirstName, user.LastName, user.BirthDate, user.Email, user.Password, user.InSearch, user.RegisteredAt, imageURL)

	if err := row.Scan(&userID); err != nil {
		return 0, errors.New("user already exists")
	}

	for _, sphere := range spheres {
		query = fmt.Sprintf("INSERT INTO %s (user_id, sphere_id) VALUES ($1, $2)", postgres.UsersSpheresTable)
		_, err := u.db.Exec(query, userID, sphere.ID)
		if err != nil {
			logrus.Error(err)
			continue
		}
	}

	for _, skill := range skills {
		query = fmt.Sprintf("INSERT INTO %s (user_id, skill_id) VALUES ($1, $2)", postgres.UsersSkillsTable)
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
		&user.Email, &user.Password, &user.InSearch, &user.RegisteredAt, &user.ImageURL)

	return user, err
}

func (u *UsersRepos) GetByID(id int) (entity.User, error) {
	var user entity.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", postgres.UsersTable)
	raw, err := u.db.Query(query, id)
	if err != nil {
		return user, errors.New("invalid user id")
	}

	raw.Next()
	err = raw.Scan(&user.ID, &user.FirstName, &user.LastName, &user.BirthDate,
		&user.Email, &user.Password, &user.InSearch, &user.RegisteredAt, &user.ImageURL)

	return user, err
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

// GetProfileInfo gather additional user information with passed user id.
// It returns all additional user profile info.
func (u *UsersRepos) GetProfileInfo(id int) ([]string, error) {
	info := make([]string, 6)

	var profileID int
	query := fmt.Sprintf("SELECT profile_id FROM %s WHERE user_id=$1", postgres.UsersProfilesTable)
	row, err := u.db.Query(query, id)
	if err != nil {
		return info, nil
	}
	row.Next()
	if err = row.Scan(&profileID); err != nil {
		return info, err
	}

	query = fmt.Sprintf("SELECT * FROM %s WHERE id=$1", postgres.ProfilesTable)
	row, err = u.db.Query(query, profileID)
	if err != nil {
		return info, err
	}

	row.Next()
	err = row.Scan(&profileID, &info[0], &info[1], &info[2], &info[3], &info[4], &info[5])

	return info, err
}

// GetSkills gather user skills with passed user id.
// It returns all user skills.
func (u *UsersRepos) GetSkills(id int) ([]entity.Skill, error) {
	var skills []entity.Skill

	query := fmt.Sprintf("SELECT skill_id FROM %s WHERE user_id=$1", postgres.UsersSkillsTable)
	rows, err := u.db.Query(query, id)
	if err != nil {
		return skills, err
	}

	for rows.Next() {
		var skillID int
		if err = rows.Scan(&skillID); err != nil {
			logrus.Error(err)
			continue
		}

		var skill entity.Skill
		query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", postgres.SkillsTable)
		skillRow, err := u.db.Query(query, skillID)
		if err != nil {
			logrus.Error(err)
			continue
		}

		skillRow.Next()
		if err = skillRow.Scan(&skill.ID, &skill.Name); err != nil {
			logrus.Error(err)
			continue
		}

		skills = append(skills, skill)
	}

	return skills, nil
}

// GetJobs gather user jobs with passed user id.
// It returns all user jobs.
func (u *UsersRepos) GetJobs(userID int) ([]entity.Job, error) {
	var jobs []entity.Job

	query := fmt.Sprintf("SELECT job_id FROM %s WHERE user_id=$1", postgres.UsersJobsTable)
	rows, err := u.db.Query(query, userID)
	if err != nil {
		return jobs, err
	}

	for rows.Next() {
		var jobID int
		if err = rows.Scan(&jobID); err != nil {
			logrus.Error(err)
			continue
		}

		// Getting jobs info
		var job entity.Job
		query = fmt.Sprintf("SELECT * FROM %s WHERE id=$1", postgres.JobsTable)
		jobRow, err := u.db.Query(query, jobID)
		if err != nil {
			logrus.Error(err)
			continue
		}
		jobRow.Next()
		if err = jobRow.Scan(&jobID, &job.CompanyName, &job.Position, &job.WorkFrom, &job.WorkTo, &job.Responsibilities); err != nil {
			logrus.Error(err)
			continue
		}

		// Getting job skills
		var skills []entity.Skill
		query = fmt.Sprintf("SELECT skill_id FROM %s WHERE job_id=$1", postgres.JobsSkillsTable)
		skillIDRows, err := u.db.Query(query, jobID)
		if err != nil {
			logrus.Error(err)
			continue
		}

		for skillIDRows.Next() {
			var skillID int
			if err = skillIDRows.Scan(&skillID); err != nil {
				logrus.Error(err)
				continue
			}

			query = fmt.Sprintf("SELECT * FROM %s WHERE id=$1", postgres.SkillsTable)
			skillRow, err := u.db.Query(query, skillID)
			if err != nil {
				logrus.Error(err)
				continue
			}

			var skill entity.Skill
			skillRow.Next()
			if err = skillRow.Scan(&skill.ID, &skill.Name); err != nil {
				logrus.Error(err)
				continue
			}

			skills = append(skills, skill)
		}
		job.Skills = skills

		jobs = append(jobs, job)
	}

	return jobs, nil
}

// DeleteAllSessions removes all user sessions from user sessions table with passed user id.
// It returns the error.
func (u *UsersRepos) DeleteAllSessions(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1", postgres.UsersSessionsTable)
	_, err := u.db.Exec(query, id)
	return err
}

// DeleteAllAgentSessions removes all user sessions from users sessions table with passed user-agent.
// It returns the error.
func (u *UsersRepos) DeleteAllAgentSessions(id int, userAgent string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1 AND user_agent=$2", postgres.UsersSessionsTable)
	_, err := u.db.Exec(query, id, userAgent)
	return err
}

// CreateSession adds a new session for the user with the passed id.
// If the number of existing sessions is exceeded, then all of them are deleted, the last one is saved.
// It returns the error.
func (u *UsersRepos) CreateSession(id int, session entity.Session) error {
	var sessionsAmount []int
	query := fmt.Sprintf("SELECT COUNT(user_id) as amount FROM %s WHERE user_id=$1 AND user_agent=$2", postgres.UsersSessionsTable)
	if err := u.db.Select(&sessionsAmount, query, id, session.UserAgent); err != nil {
		return err
	}
	if sessionsAmount[0] == postgres.UsersMaxSessionsAmount {
		if err := u.DeleteAllAgentSessions(id, session.UserAgent); err != nil {
			return err
		}
	}

	query = fmt.Sprintf("INSERT INTO %s (user_id, user_agent, refresh_token, expires_at) VALUES ($1, $2, $3, $4)",
		postgres.UsersSessionsTable)
	_, err := u.db.Exec(query, id, session.UserAgent, session.RefreshToken, session.ExpiresAt)
	return err
}

// UpdateSession updates an existing session for a user with the passed refresh token.
// It returns the error.
func (u *UsersRepos) UpdateSession(id int, refreshToken string, session entity.Session) error {
	query := fmt.Sprintf("UPDATE %s SET refresh_token=$1, expires_at=$2 WHERE user_id=$3 AND user_agent=$4 AND refresh_token=$5",
		postgres.UsersSessionsTable)
	_, err := u.db.Exec(query, session.RefreshToken, session.ExpiresAt, id, session.UserAgent, refreshToken)

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
