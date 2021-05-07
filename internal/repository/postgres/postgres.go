package postgres

import (
	"fmt"
	"github.com/PolyProjectOPD/Backend/internal/config"
	"github.com/jmoiron/sqlx"
)

const (
	UsersTable                     = "users"
	UsersSessionsTable             = "users_sessions"
	UsersSpheresTable              = "users_spheres"
	UsersSkillsTable               = "users_skills"
	UsersProfilesTable             = "users_profiles"
	UsersJobsTable                 = "users_jobs"
	UsersCompaniesTable            = "users_companies"
	ProfilesTable                  = "profiles"
	JobsTable                      = "jobs"
	JobsSkillsTable                = "jobs_skills"
	SkillsTable                    = "skills"
	SpheresTable                   = "spheres"
	SpheresSkillsTable             = "spheres_skills"
	CompaniesTable                 = "companies"
	CompaniesVacanciesTable        = "companies_vacancies"
	VacanciesTable                 = "vacancies"
	VacanciesResponsibilitiesTable = "vacancies_responsibilities"
	VacanciesRequirementsTable     = "vacancies_requirements"
	VacanciesConditionsTable       = "vacancies_conditions"
	VacanciesSkillsTable           = "vacancies_skills"
	ResponsibilitiesTable          = "responsibilities"
	RequirementsTable              = "requirements"
	ConditionsTable                = "conditions"
	UsersMaxSessionsAmount         = 5
)

func NewPostgresDB(cfg *config.DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect(cfg.DriverName, fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Password, cfg.SSLMode))

	if err != nil {
		return nil, err
	}

	return db, nil
}
