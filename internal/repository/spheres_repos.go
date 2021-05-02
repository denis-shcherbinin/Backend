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

// GetAll forms a slice of all spheres from DB.
// It returns all spheres and error.
func (s *SpheresRepos) GetAll() ([]entity.Sphere, error) {
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

	return spheres, nil
}

// GetSkills forms a slice of all skills according to the passed sphere.
// It returns all skills of sphere and error.
func (s *SpheresRepos) GetSkills(sphere entity.Sphere) ([]entity.Skill, error) {
	var skills []entity.Skill

	skillsQuery := fmt.Sprintf("SELECT s.id, s.name FROM %s s INNER JOIN %s ss on s.id=ss.skill_id WHERE ss.sphere_id=$1",
		postgres.SkillsTable, postgres.SpheresSkillsTable)
	skillsRows, err := s.db.Query(skillsQuery, sphere.ID)
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
