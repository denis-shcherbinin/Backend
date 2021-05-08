package service

import (
	"errors"
	"github.com/PolyProjectOPD/Backend/internal/entity"
	"github.com/PolyProjectOPD/Backend/internal/repository"
	"github.com/PolyProjectOPD/Backend/internal/storage"
	"github.com/sirupsen/logrus"
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
		Name:             input.Name,
		Location:         input.Location,
		ShortDescription: input.ShortDescription,
		FullDescription:  input.FullDescription,
		ImageURL:         imageURL,
	}

	return c.repos.Create(userID, company)
}

func (c *CompaniesService) Profile(userID int) (entity.CompanyProfile, error) {
	return c.repos.Profile(userID)
}

func (c *CompaniesService) UpdateProfile(userID int, companyProfile entity.CompanyProfile, fileBody, fileType string) error {
	companyID, err := c.repos.GetIDByUserID(userID)
	if err != nil {
		return err
	}
	if companyID != companyProfile.Company.ID {
		return errors.New("invalid company id")
	}

	imageURL, err := c.storage.Upload(storage.UploadInput{
		Body:        fileBody,
		ContentType: fileType,
	})
	if err != nil {
		return err
	}
	companyProfile.Company.ImageURL = imageURL

	if err = c.DeleteImage(companyID); err != nil {
		logrus.Error(err)
	}

	return c.repos.UpdateProfile(companyProfile)
}

func (c *CompaniesService) DeleteImage(companyID int) error {
	imageURL, err := c.repos.GetImageURL(companyID)
	if err != nil {
		return err
	}

	if imageURL == "" {
		return nil
	}

	if err = c.storage.Delete(imageURL); err != nil {
		return err
	}

	return c.repos.DeleteImage(companyID)
}
