package service

import (
	"github.com/PolyProjectOPD/Backend/internal/entity"
	"github.com/PolyProjectOPD/Backend/internal/repository"
)

type SkillsService struct {
	repos repository.Skills
}

func NewSkillsService(repos repository.Skills) *SkillsService {
	return &SkillsService{
		repos: repos,
	}
}

func (s *SkillsService) GetAll() ([]entity.Skill, error) {
	return s.repos.GetAll()
}
