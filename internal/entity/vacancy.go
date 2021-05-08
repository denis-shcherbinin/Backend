package entity

type Vacancy struct {
	ID               int                     `json:"id"`
	Position         string                  `json:"position"`
	Description      string                  `json:"description"`
	IsFullTime       bool                    `json:"isFullTime"`
	MinSalary        string                  `json:"minSalary"`
	MaxSalary        string                  `json:"maxSalary"`
	SkillLevel       string                  `json:"skillLevel"`
	Responsibilities []VacancyResponsibility `json:"responsibilities"`
	Requirements     []VacancyRequirement    `json:"requirements"`
	Conditions       []VacancyCondition      `json:"conditions"`
	Skills           []Skill                 `json:"skills"`
}

type VacancyResponsibility struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type VacancyRequirement struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type VacancyCondition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
