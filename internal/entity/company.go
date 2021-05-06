package entity

type Company struct {
	Name              string `json:"name" binding:"required"`
	Location          string `json:"location" binding:"required"`
	FoundationDate    string `json:"foundationDate" binding:"required"`
	NumberOfEmployees int    `json:"numberOfEmployees" binding:"required"`
	ShortDescription  string `json:"shortDescription" binding:"required"`
	FullDescription   string `json:"fullDescription" binding:"required"`
	ImageURL          string `json:"imageURL" binding:"required"`
}
