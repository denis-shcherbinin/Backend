package entity

type Company struct {
	ID               int    `json:"id"`
	Name             string `json:"name" binding:"required"`
	Location         string `json:"location" binding:"required"`
	ShortDescription string `json:"shortDescription" binding:"required"`
	FullDescription  string `json:"fullDescription" binding:"required"`
	ImageURL         string `json:"imageURL" binding:"required"`
}

type CompanyProfile struct {
	Company   Company   `json:"company"`
	Vacancies []Vacancy `json:"vacancies"`
}
