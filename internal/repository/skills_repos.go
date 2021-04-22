package repository

import (
	"fmt"
	"github.com/PolyProjectOPD/Backend/internal/entity"
	"github.com/PolyProjectOPD/Backend/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type SkillsRepos struct {
	db *sqlx.DB
}

func NewSkillsRepos(db *sqlx.DB) *SkillsRepos {
	return &SkillsRepos{
		db: db,
	}
}

// GetAll forms a slice of all skills.
// It returns slice of all skills.
func (s *SkillsRepos) GetAll() ([]entity.Skill, error) {
	var skills []entity.Skill

	query := fmt.Sprintf("SELECT * FROM %s", postgres.SkillsTable)

	rows, err := s.db.Query(query)
	if err != nil {
		return skills, err
	}

	for rows.Next() {
		var skill entity.Skill

		if err = rows.Scan(&skill.ID, &skill.Name); err != nil {
			logrus.Error(err)
			continue
		}

		skills = append(skills, skill)
	}

	return skills, nil
}
