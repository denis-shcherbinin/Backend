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

	// User spheres
	for _, sphere := range spheres {
		query = fmt.Sprintf("INSERT INTO %s (user_id, sphere_id) VALUES ($1, $2)", postgres.UsersSpheresTable)
		_, err := u.db.Exec(query, userID, sphere.ID)
		if err != nil {
			logrus.Error(err)
			continue
		}
	}

	// User skills
	for _, skill := range skills {
		query = fmt.Sprintf("INSERT INTO %s (user_id, skill_id) VALUES ($1, $2)", postgres.UsersSkillsTable)
		_, err := u.db.Exec(query, userID, skill.ID)
		if err != nil {
			logrus.Error(err)
			continue
		}
	}

	// Profile create
	var profileID int
	query = fmt.Sprintf("INSERT INTO %s DEFAULT VALUES RETURNING id", postgres.ProfilesTable)
	row = u.db.QueryRow(query)
	if err := row.Scan(&profileID); err != nil {
		return userID, err
	}

	query = fmt.Sprintf("INSERT INTO %s (user_id, profile_id) VALUES ($1, $2)", postgres.UsersProfilesTable)
	_, err := u.db.Exec(query, userID, profileID)
	if err != nil {
		return userID, err
	}

	return userID, nil
}

// GetByCredentials looks in the users table for the presence of a user with passed credentials(email, password).
// It returns user and error.
func (u *UsersRepos) GetByCredentials(email, password string) (entity.User, error) {
	var user entity.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1 AND password_hash=$2", postgres.UsersTable)

	raw := u.db.QueryRow(query, email, password)

	err := raw.Scan(&user.ID, &user.FirstName, &user.LastName, &user.BirthDate,
		&user.Email, &user.Password, &user.InSearch, &user.RegisteredAt, &user.ImageURL)

	return user, err
}

func (u *UsersRepos) GetByID(id int) (entity.User, error) {
	var user entity.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", postgres.UsersTable)
	raw := u.db.QueryRow(query, id)

	err := raw.Scan(&user.ID, &user.FirstName, &user.LastName, &user.BirthDate,
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

	query := fmt.Sprintf("SELECT p.comment, p.experience, p.skill_level, p.min_salary, p.max_salary, p.about FROM %s "+
		"p INNER JOIN %s up on p.id=up.profile_id WHERE up.user_id=$1", postgres.ProfilesTable, postgres.UsersProfilesTable)
	row := u.db.QueryRow(query, id)

	err := row.Scan(&info[0], &info[1], &info[2], &info[3], &info[4], &info[5])

	return info, err
}

// GetSkills gather user skills with passed user id.
// It returns all user skills.
func (u *UsersRepos) GetSkills(id int) ([]entity.Skill, error) {
	var skills []entity.Skill

	skillsQuery := fmt.Sprintf("SELECT s.id, s.name FROM %s s INNER JOIN %s us on s.id=us.skill_id WHERE us.user_id=$1",
		postgres.SkillsTable, postgres.UsersSkillsTable)
	skillsRows, err := u.db.Query(skillsQuery, id)
	if err != nil {
		return skills, err
	}

	for skillsRows.Next() {
		var skill entity.Skill

		if err = skillsRows.Scan(&skill.ID, &skill.Name); err != nil {
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

	jobsQuery := fmt.Sprintf("SELECT j.id, j.company_name, j.position, j.from, j.to, j.responsibilities FROM %s "+
		"j INNER JOIN %s uj on j.id=uj.job_id WHERE uj.user_id=$1", postgres.JobsTable, postgres.UsersJobsTable)
	jobsRows, err := u.db.Query(jobsQuery, userID)
	if err != nil {
		return jobs, err
	}

	var jobsID []int
	for jobsRows.Next() {
		var jobID int
		var job entity.Job
		if err = jobsRows.Scan(&jobID, &job.CompanyName, &job.Position, &job.WorkFrom, &job.WorkTo, &job.Responsibilities); err != nil {
			logrus.Error(err)
			continue
		}

		jobsID = append(jobsID, jobID)
		jobs = append(jobs, job)
	}

	// Так как в пакете для работы с БД нет нормального решения для обработки двух query в цикле(либо надо искать лучше),
	// используется второй цикл и доп. массив с jobs id
	for i, id := range jobsID {
		jobs[i].Skills, err = u.getJobsSkills(id)
		if err != nil {
			logrus.Error(err)
		}
	}

	return jobs, nil
}

func (u *UsersRepos) GetImageURL(id int) (string, error) {
	query := fmt.Sprintf("SELECT image_url FROM %s WHERE id=$1", postgres.UsersTable)
	row := u.db.QueryRow(query, id)

	var imageURL string
	if err := row.Scan(&imageURL); err != nil {
		return "", err
	}

	return imageURL, nil
}

func (u *UsersRepos) DeleteImage(id int) error {
	query := fmt.Sprintf("UPDATE %s SET image_url=$1 WHERE id=$2", postgres.UsersTable)
	_, err := u.db.Exec(query, "", id)
	if err != nil {
		return err
	}

	return nil
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
	query := fmt.Sprintf("SELECT COUNT(user_id) as amount FROM %s WHERE user_id=$1 AND user_agent=$2",
		postgres.UsersSessionsTable)
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

func (u *UsersRepos) getJobsSkills(jobID int) ([]entity.Skill, error) {
	var skills []entity.Skill

	skillsQuery := fmt.Sprintf("SELECT s.id, s.name FROM %s s "+
		"INNER JOIN %s js on s.id=js.skill_id WHERE js.job_id=$1", postgres.SkillsTable, postgres.JobsSkillsTable)
	skillsRows, err := u.db.Query(skillsQuery, jobID)
	if err != nil {
		return skills, err
	}

	for skillsRows.Next() {
		var skill entity.Skill

		if err = skillsRows.Scan(&skill.ID, &skill.Name); err != nil {
			logrus.Error(err)
			continue
		}

		skills = append(skills, skill)
	}

	return skills, nil
}
