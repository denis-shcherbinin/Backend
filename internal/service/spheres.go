package service

import (
	"github.com/PolyProjectOPD/Backend/internal/entity"
	"github.com/PolyProjectOPD/Backend/internal/repository"
)

type SpheresService struct {
	repos repository.Spheres
}

func NewSpheresService(repos repository.Spheres) *SpheresService {
	return &SpheresService{
		repos: repos,
	}
}

// GetAll forms a slice of all spheres.
// It returns slice of all spheres and error.
func (s *SpheresService) GetAll() ([]entity.Sphere, error) {
	return s.repos.Get()
}

// GetSkills forms slice of skills according to the passed slice of spheres.
// It returns slice of skills and error.
func (s *SpheresService) GetSkills(spheres []entity.Sphere) ([]entity.Skill, error) {
	var skills []entity.Skill

	for _, sphere := range spheres {
		sphereSkills, err := s.repos.GetSkills(sphere)
		if err != nil {
			continue
		}

		skills = append(skills, sphereSkills...)
	}

	return skills, nil
}
