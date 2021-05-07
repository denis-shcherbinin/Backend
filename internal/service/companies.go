package service

import (
	"github.com/PolyProjectOPD/Backend/internal/entity"
	"github.com/PolyProjectOPD/Backend/internal/repository"
	"github.com/PolyProjectOPD/Backend/internal/storage"
)

type CompaniesService struct {
	repos   repository.Companies
	storage *storage.Storage
}

func NewCompaniesService(repos repository.Companies, storage *storage.Storage) *CompaniesService {
	return &CompaniesService{
		repos:   repos,
		storage: storage,
	}
}

func (c *CompaniesService) Create(userID int, input entity.CompanyInput, fileBody, fileType string) (int, error) {
	imageURL, err := c.storage.Upload(storage.UploadInput{
		Body:        fileBody,
		ContentType: fileType,
	})
	if err != nil {
		return 0, err
	}

	company := entity.Company{
		Name:              input.Name,
		Location:          input.Location,
		ShortDescription:  input.ShortDescription,
		FullDescription:   input.FullDescription,
		ImageURL:          imageURL,
	}

	return c.repos.Create(userID, company)
}

func (c *CompaniesService) Profile(userID int) (entity.CompanyProfile, error) {
	return c.repos.Profile(userID)
}
