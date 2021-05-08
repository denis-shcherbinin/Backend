package repository

import (
	"fmt"
	"github.com/PolyProjectOPD/Backend/internal/entity"
	"github.com/PolyProjectOPD/Backend/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type CompaniesRepos struct {
	db *sqlx.DB
}

func NewCompaniesRepos(db *sqlx.DB) *CompaniesRepos {
	return &CompaniesRepos{
		db: db,
	}
}

func (c *CompaniesRepos) Create(userID int, company entity.Company) (int, error) {
	var companyID int
	query := fmt.Sprintf("INSERT INTO %s "+
		"(name, location, short_description, full_description, image_url) "+
		"VALUES($1, $2, $3, $4, $5) RETURNING ID", postgres.CompaniesTable)
	row := c.db.QueryRow(query, company.Name, company.Location, company.ShortDescription, company.FullDescription,
		company.ImageURL)
	if err := row.Scan(&companyID); err != nil {
		return 0, err
	}

	query = fmt.Sprintf("INSERT INTO %s (user_id, company_id) VALUES ($1, $2)", postgres.UsersCompaniesTable)
	_, err := c.db.Exec(query, userID, companyID)
	if err != nil {
		return companyID, err
	}

	return companyID, nil
}

func (c *CompaniesRepos) Profile(userID int) (entity.CompanyProfile, error) {
	company, err := c.getCompany(userID)
	if err != nil {
		return entity.CompanyProfile{}, err
	}

	vacancies, err := c.getCompanyVacancies(company.ID)
	if err != nil {
		return entity.CompanyProfile{
			Company: company,
		}, err
	}

	return entity.CompanyProfile{
		Company:   company,
		Vacancies: vacancies,
	}, nil
}

func (c *CompaniesRepos) GetIDByUserID(userID int) (int, error) {
	var companyID int
	query := fmt.Sprintf("SELECT company_id FROM %s WHERE user_id=$1", postgres.UsersCompaniesTable)
	row := c.db.QueryRow(query, userID)
	if err := row.Scan(&companyID); err != nil {
		return 0, err
	}

	return companyID, nil
}

func (c *CompaniesRepos) UpdateProfile(companyProfile entity.CompanyProfile) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET name=$1, location=$2, short_description=$3, full_description=$4, "+
		"image_url=$5 WHERE id=$6", postgres.CompaniesTable)
	_, err = tx.Exec(query, companyProfile.Company.Name, companyProfile.Company.Location,
		companyProfile.Company.ShortDescription, companyProfile.Company.FullDescription,
		companyProfile.Company.ImageURL, companyProfile.Company.ID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Getting all company vacanciesID
	var vacanciesID []int
	query = fmt.Sprintf("SELECT vacancy_id FROM %s WHERE company_id=$1", postgres.CompaniesVacanciesTable)
	rows, err := tx.Query(query, companyProfile.Company.ID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	for rows.Next() {
		var vacancyID int
		if err = rows.Scan(&vacancyID); err != nil {
			_ = tx.Rollback()
			return err
		}
		vacanciesID = append(vacanciesID, vacancyID)
	}

	// Deleting all company vacancies
	for _, vacancyID := range vacanciesID {
		if err = c.deleteVacancyTx(vacancyID, tx); err != nil {
			return err
		}
	}

	// Inserting new company vacancies
	for _, vacancy := range companyProfile.Vacancies {
		if err = c.createVacancyTx(companyProfile.Company.ID, vacancy, tx); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (c *CompaniesRepos) GetImageURL(companyID int) (string, error) {
	var imageURL string
	query := fmt.Sprintf("SELECT image_url FROM %s WHERE id=$1", postgres.CompaniesTable)
	row := c.db.QueryRow(query, companyID)
	if err := row.Scan(&imageURL); err != nil {
		return "", err
	}
	return imageURL, nil
}

func (c *CompaniesRepos) DeleteImage(companyID int) error {
	query := fmt.Sprintf("UPDATE %s SET image_url=$1 WHERE id=$2", postgres.CompaniesTable)
	_, err := c.db.Exec(query, "", companyID)
	return err
}

func (c *CompaniesRepos) getCompany(userID int) (entity.Company, error) {
	var company entity.Company
	query := fmt.Sprintf("SELECT c.id, c.name, c.location, c.short_description, c.full_description, c.image_url FROM %s "+
		"c INNER JOIN %s uc on c.id=uc.company_id WHERE uc.user_id=$1",
		postgres.CompaniesTable, postgres.UsersCompaniesTable)
	row := c.db.QueryRow(query, userID)
	if err := row.Scan(&company.ID, &company.Name, &company.Location, &company.ShortDescription,
		&company.FullDescription, &company.ImageURL); err != nil {
		return company, err
	}

	return company, nil
}

func (c *CompaniesRepos) getCompanyVacancies(companyID int) ([]entity.Vacancy, error) {
	var vacancies []entity.Vacancy

	query := fmt.Sprintf("SELECT v.id, v.position, v.description, v.is_full_time, v.min_salary, v.max_salary, v.skill_level "+
		"FROM %s v INNER JOIN %s cv on v.id=cv.vacancy_id WHERE cv.company_id=$1",
		postgres.VacanciesTable, postgres.CompaniesVacanciesTable)
	rows, err := c.db.Query(query, companyID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var vacancy entity.Vacancy
		if err = rows.Scan(&vacancy.ID, &vacancy.Position, &vacancy.Description, &vacancy.IsFullTime,
			&vacancy.MinSalary, &vacancy.MaxSalary, &vacancy.SkillLevel); err != nil {
			logrus.Error(err)
			continue
		}

		vacancies = append(vacancies, vacancy)
	}

	for _, vacancy := range vacancies {
		vacancy.Responsibilities, err = c.getVacancyResponsibilities(vacancy.ID)
		if err != nil {
			logrus.Error(err)
		}

		vacancy.Requirements, err = c.getVacancyRequirements(vacancy.ID)
		if err != nil {
			logrus.Error(err)
		}

		vacancy.Conditions, err = c.getVacancyConditions(vacancy.ID)
		if err != nil {
			logrus.Error(err)
		}

		vacancy.Skills, err = c.getVacancySkills(vacancy.ID)
		if err != nil {
			logrus.Error(err)
		}
	}

	return vacancies, nil
}

func (c *CompaniesRepos) getVacancyResponsibilities(vacancyID int) ([]entity.VacancyResponsibility, error) {
	query := fmt.Sprintf("SELECT r.name, r.description FROM %s r INNER JOIN %s vr on r.id=vr.responsibility_id WHERE vr.vacancy_id=$1",
		postgres.ResponsibilitiesTable, postgres.VacanciesResponsibilitiesTable)
	rows, err := c.db.Query(query, vacancyID)
	if err != nil {
		return nil, err
	}

	var responsibilities []entity.VacancyResponsibility
	for rows.Next() {
		var responsibility entity.VacancyResponsibility
		if err = rows.Scan(&responsibility); err != nil {
			logrus.Error(err)
			continue
		}
		responsibilities = append(responsibilities, responsibility)
	}

	return responsibilities, nil
}

func (c *CompaniesRepos) getVacancyRequirements(vacancyID int) ([]entity.VacancyRequirement, error) {
	query := fmt.Sprintf("SELECT r.name, r.description FROM %s r INNER JOIN %s vr on r.id=vr.requirement_id WHERE vr.vacancy_id=$1",
		postgres.RequirementsTable, postgres.VacanciesRequirementsTable)
	rows, err := c.db.Query(query, vacancyID)
	if err != nil {
		return nil, err
	}

	var requirements []entity.VacancyRequirement
	for rows.Next() {
		var requirement entity.VacancyRequirement
		if err = rows.Scan(&requirement); err != nil {
			logrus.Error(err)
			continue
		}
		requirements = append(requirements, requirement)
	}

	return requirements, nil
}

func (c *CompaniesRepos) getVacancyConditions(vacancyID int) ([]entity.VacancyCondition, error) {
	query := fmt.Sprintf("SELECT c.name, c.description FROM %s c INNER JOIN %s vc on c.id=vc.requirement_id WHERE vc.vacancy_id=$1",
		postgres.ConditionsTable, postgres.VacanciesConditionsTable)
	rows, err := c.db.Query(query, vacancyID)
	if err != nil {
		return nil, err
	}

	var conditions []entity.VacancyCondition
	for rows.Next() {
		var condition entity.VacancyCondition
		if err = rows.Scan(&condition); err != nil {
			logrus.Error(err)
			continue
		}
		conditions = append(conditions, condition)
	}

	return conditions, nil
}

func (c *CompaniesRepos) getVacancySkills(vacancyID int) ([]entity.Skill, error) {
	query := fmt.Sprintf("SELECT s.id, s.name FROM %s s INNER JOIN %s vs on s.id=vs.skill_id WHERE vs.vacancy_id=$1",
		postgres.SkillsTable, postgres.VacanciesSkillsTable)
	rows, err := c.db.Query(query, vacancyID)
	if err != nil {
		return nil, err
	}

	var skills []entity.Skill
	for rows.Next() {
		var skill entity.Skill
		if err = rows.Scan(&skill); err != nil {
			logrus.Error(err)
			continue
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

func (c *CompaniesRepos) deleteVacancyTx(vacancyID int, tx *sqlx.Tx) error {
	// DELETE vacancy responsibilities
	query := fmt.Sprintf("SELECT r.id FROM %s r INNER JOIN %s vr on r.id=vr.responsibility_id WHERE vr.vacancy_id=$1",
		postgres.ResponsibilitiesTable, postgres.VacanciesResponsibilitiesTable)
	rows, err := tx.Query(query, vacancyID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	var responsibilitiesID []int
	for rows.Next() {
		var responsibilityID int
		if err = rows.Scan(&responsibilityID); err != nil {
			_ = tx.Rollback()
			return err
		}
		responsibilitiesID = append(responsibilitiesID, responsibilityID)
	}
	for id := range responsibilitiesID {
		query = fmt.Sprintf("DELETE FROM %s WHERE id=$1", postgres.ResponsibilitiesTable)
		_, err = tx.Exec(query, id)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// DELETE vacancy requirements
	query = fmt.Sprintf("SELECT r.id FROM %s r INNER JOIN %s vr on r.id=vr.requirement_id WHERE vr.vacancy_id=$1",
		postgres.RequirementsTable, postgres.VacanciesRequirementsTable)
	rows, err = tx.Query(query, vacancyID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	var requirementsID []int
	for rows.Next() {
		var requirementID int
		if err = rows.Scan(&requirementID); err != nil {
			_ = tx.Rollback()
			return err
		}
		requirementsID = append(requirementsID, requirementID)
	}
	for id := range requirementsID {
		query = fmt.Sprintf("DELETE FROM %s WHERE id=$1", postgres.RequirementsTable)
		_, err = tx.Exec(query, id)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// DELETE vacancy conditions
	query = fmt.Sprintf("SELECT c.id FROM %s c INNER JOIN %s vc on c.id=vc.condition_id WHERE vc.vacancy_id=$1",
		postgres.ConditionsTable, postgres.VacanciesConditionsTable)
	rows, err = tx.Query(query, vacancyID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	var conditionsID []int
	for rows.Next() {
		var conditionID int
		if err = rows.Scan(&conditionID); err != nil {
			_ = tx.Rollback()
			return err
		}
		conditionsID = append(conditionsID, conditionID)
	}
	for id := range conditionsID {
		query = fmt.Sprintf("DELETE FROM %s WHERE id=$1", postgres.ConditionsTable)
		_, err = tx.Exec(query, id)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// DELETE vacancy and vacancies_skills relations table
	query = fmt.Sprintf("DELETE FROM %s WHERE id=$1", postgres.VacanciesTable)
	_, err = tx.Exec(query, vacancyID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return nil
}

func (c *CompaniesRepos) createVacancyTx(companyID int, vacancy entity.Vacancy, tx *sqlx.Tx) error {
	var vacancyID int
	query := fmt.Sprintf("INSERT INTO %s (position, description, is_full_time, min_salary, max_salary, skill_level) "+
		"VALUES($1, $2, $3, $4, $5, $6) RETURNING id", postgres.VacanciesTable)
	row := tx.QueryRow(query, vacancy.Position, vacancy.Description, vacancy.IsFullTime, vacancy.MinSalary,
		vacancy.MaxSalary, vacancy.SkillLevel)
	if err := row.Scan(&vacancyID); err != nil {
		_ = tx.Rollback()
		return err
	}

	query = fmt.Sprintf("INSERT INTO %s (company_id, vacancy_id) VALUES ($1, $2)", postgres.CompaniesVacanciesTable)
	_, err := tx.Exec(query, companyID, vacancyID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	for _, responsibility := range vacancy.Responsibilities {
		var id int
		query = fmt.Sprintf("INSERT INTO %s (name, description) VALUES ($1, $2) RETURNING id", postgres.ResponsibilitiesTable)
		row = tx.QueryRow(query, responsibility.Name, responsibility.Description)
		if err = row.Scan(&id); err != nil {
			_ = tx.Rollback()
			return err
		}

		query = fmt.Sprintf("INSERT INTO %s (vacancy_id, responsibility_id) VALUES ($1, $2)", postgres.VacanciesResponsibilitiesTable)
		_, err = tx.Exec(query, vacancyID, id)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	for _, requirement := range vacancy.Requirements {
		var id int
		query = fmt.Sprintf("INSERT INTO %s (name, description) VALUES ($1, $2) RETURNING id", postgres.RequirementsTable)
		row = tx.QueryRow(query, requirement.Name, requirement.Description)
		if err = row.Scan(&id); err != nil {
			_ = tx.Rollback()
			return err
		}

		query = fmt.Sprintf("INSERT INTO %s (vacancy_id, requirement_id) VALUES ($1, $2)", postgres.VacanciesRequirementsTable)
		_, err = tx.Exec(query, vacancyID, id)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	for _, condition := range vacancy.Conditions {
		var id int
		query = fmt.Sprintf("INSERT INTO %s (name, description) VALUES ($1, $2) RETURNING id", postgres.ConditionsTable)
		row = tx.QueryRow(query, condition.Name, condition.Description)
		if err = row.Scan(&id); err != nil {
			_ = tx.Rollback()
			return err
		}

		query = fmt.Sprintf("INSERT INTO %s (vacancy_id, condition_id) VALUES ($1, $2)", postgres.VacanciesConditionsTable)
		_, err = tx.Exec(query, vacancyID, id)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	for _, skill := range vacancy.Skills {
		query = fmt.Sprintf("INSERT INTO %s (vacancy_id, skill_id) VALUES ($1, $2)", postgres.VacanciesSkillsTable)
		_, err = tx.Exec(query, vacancyID, skill.ID)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return nil
}
