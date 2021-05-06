package repository

import (
	"fmt"
	"github.com/PolyProjectOPD/Backend/internal/entity"
	"github.com/PolyProjectOPD/Backend/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
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
		"(name, location, foundation_date, number_of_employees, short_description, full_description, image_url) "+
		"VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING ID", postgres.CompaniesTable)
	row := c.db.QueryRow(query, company.Name, company.Location, company.FoundationDate,
		company.NumberOfEmployees, company.ShortDescription, company.FullDescription, company.ImageURL)
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
