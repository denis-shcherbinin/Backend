package entity

type UserProfile struct {
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	Email      string  `json:"email"`
	Age        string  `json:"age"`
	ImageURL   string  `json:"imageURL,"`
	Comment    string  `json:"comment"`
	Experience string  `json:"experience"`
	SkillLevel string  `json:"skillLevel"`
	MinSalary  string  `json:"minSalary"`
	MaxSalary  string  `json:"maxSalary"`
	About      string  `json:"about"`
	Skills     []Skill `json:"skills"`
	Jobs       []Job   `json:"jobs"`
}
