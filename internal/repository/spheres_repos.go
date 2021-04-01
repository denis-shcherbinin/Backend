package repository

import (
	"fmt"
	"github.com/PolyProjectOPD/Backend/internal/entity"
	"github.com/PolyProjectOPD/Backend/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type SpheresRepos struct {
	db *sqlx.DB
}

func NewSpheresRepos(db *sqlx.DB) *SpheresRepos {
	return &SpheresRepos{
		db: db,
	}
}

// Get forms a slice of all spheres from DB.
// It returns all spheres and error.
func (s *SpheresRepos) Get() ([]entity.Sphere, error) {
	var spheres []entity.Sphere

	query := fmt.Sprintf("SELECT * FROM %s", postgres.SpheresTable)

	rows, err := s.db.Query(query)
	if err != nil {
		return spheres, err
	}

	for rows.Next() {
		var sphere entity.Sphere

		if err = rows.Scan(&sphere.ID, &sphere.Name); err != nil {
			logrus.Error(err)
			continue
		}

		spheres = append(spheres, sphere)
	}

	return spheres, err
}

// GetSkills forms a slice of all skills according to the passed sphere.
// It returns all skills of sphere and error.
func (s *SpheresRepos) GetSkills(sphere entity.Sphere) ([]entity.Skill, error) {
	var skills []entity.Skill

	query := fmt.Sprintf("SELECT skill_id FROM %s WHERE sphere_id=$1", postgres.SpheresSkillsTable)
	rows, err := s.db.Query(query, sphere.ID)
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
		query = fmt.Sprintf("SELECT * FROM %s WHERE id=$1", postgres.SkillsTable)

		skillRow, err := s.db.Query(query, skillID)
		if err != nil {
			logrus.Error(err)
			continue
		}
		if err = skillRow.Scan(&skill.ID, &skill.Name); err != nil {
			logrus.Error(err)
			continue
		}

		skills = append(skills, skill)
	}

	return skills, nil
}
