package entity

type UserExistenceInput struct {
	Email string `json:"email" binding:"required,email,max=64"`
}

type UserCredentialsInput struct {
	FirstName string `json:"firstName" binding:"required,min=2,max=64"`
	LastName  string `json:"lastName" binding:"required,min=2,max=64"`
	BirthDate string `json:"birthDate" binding:"required"`
	Email     string `json:"email" binding:"required,email,max=64"`
	Password  string `json:"password" binding:"required,min=8,max=64"`
	InSearch  bool   `json:"inSearch" binding:"required"`
}

type UsersSpheresInput struct {
	Spheres []Sphere `json:"spheres" binding:"required"`
}

type UsersSkillsInput struct {
	Skills []Skill `json:"skills" binding:"required"`
}

type UserSignUpInput struct {
	UserCredentialsInput
	UsersSpheresInput
	UsersSkillsInput
}

type UserSignInInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type UserRefreshInput struct {
	Token string `json:"token" binding:"required"`
}

type SpheresInput struct {
	Spheres []Sphere `json:"spheres" binding:"required"`
}
