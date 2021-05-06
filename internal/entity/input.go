package entity

type (
	UserExistenceInput struct {
		Email string `json:"email" binding:"required,email,max=64"`
	}

	UserCredentialsInput struct {
		FirstName string `json:"firstName" binding:"required,min=2,max=64"`
		LastName  string `json:"lastName" binding:"required,min=2,max=64"`
		BirthDate string `json:"birthDate" binding:"required"`
		Email     string `json:"email" binding:"required,email,max=64"`
		Password  string `json:"password" binding:"required,min=8,max=64"`
		InSearch  bool   `json:"inSearch" binding:"required"`
	}

	UsersSpheresInput struct {
		Spheres []Sphere `json:"spheres" binding:"required"`
	}

	UsersSkillsInput struct {
		Skills []Skill `json:"skills" binding:"required"`
	}

	UserSignUpInput struct {
		UserCredentialsInput
		UsersSpheresInput
		UsersSkillsInput
	}

	UserSignInInput struct {
		Email    string `json:"email" binding:"required,email,max=64"`
		Password string `json:"password" binding:"required,min=8,max=64"`
	}

	UserRefreshInput struct {
		Token string `json:"token" binding:"required"`
	}

	SpheresInput struct {
		Spheres []Sphere `json:"spheres" binding:"required"`
	}

	CompanyInput struct {
		Name              string `json:"name" binding:"required"`
		Location          string `json:"location" binding:"required"`
		FoundationDate    string `json:"foundationDate" binding:"required"`
		NumberOfEmployees int    `json:"numberOfEmployees" binding:"required"`
		ShortDescription  string `json:"shortDescription" binding:"required"`
		FullDescription   string `json:"fullDescription" binding:"required"`
  }
  
	ProfileInput struct {
		FirstName  string  `json:"firstName" binding:"required"`
		LastName   string  `json:"lastName" binding:"required"`
		Email      string  `json:"email" binding:"required"`
		BirthDate  string  `json:"birthDate" binding:"required"`
		Comment    string  `json:"comment" binding:"required"`
		Experience string  `json:"experience" binding:"required"`
		SkillLevel string  `json:"skillLevel" binding:"required"`
		MinSalary  string  `json:"minSalary" binding:"required"`
		MaxSalary  string  `json:"maxSalary" binding:"required"`
		About      string  `json:"about" binding:"required"`
		Skills     []Skill `json:"skills" binding:"required"`
		Jobs       []Job   `json:"jobs" binding:"required"`
	}
)
