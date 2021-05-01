package entity

type Job struct {
	CompanyName      string  `json:"companyName"`
	Position         string  `json:"position"`
	WorkFrom         string  `json:"workFrom"`
	WorkTo           string  `json:"workTo"`
	Responsibilities string  `json:"responsibilities"`
	Skills           []Skill `json:"skills"`
}
