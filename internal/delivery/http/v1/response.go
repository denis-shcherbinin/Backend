package v1

import (
	"github.com/PolyProjectOPD/Backend/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type response struct {
	Message string `json:"message"`
}

type userExistenceResponse struct {
	Exists bool `json:"exists"`
}

type signUpResponse struct {
	ID       int    `json:"id"`
	ImageURL string `json:"imageURL"`
}

type tokensResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type spheresResponse struct {
	Spheres []entity.Sphere `json:"spheres"`
}

type skillsResponse struct {
	Skills []entity.Skill `json:"skills"`
}

type userProfileResponse struct {
	Profile entity.UserProfile `json:"profile"`
}

type companyCreateResponse struct {
	CompanyID int `json:"companyID"`
}

type companyProfileResponse struct {
	CompanyProfile entity.CompanyProfile `json:"companyProfile"`
}

func newResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, response{
		Message: message,
	})
}
