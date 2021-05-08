package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initCompaniesRoutes(api *gin.RouterGroup) {
	company := api.Group("/company")
	{
		authenticated := company.Group("/", h.userIdentify)
		{
			authenticated.POST("/create", h.createCompany)
			authenticated.GET("/profile", h.companyProfile)
		}
	}
}

// @Summary Company Creation
// @Security UserAuth
// @Tags Company
// @Description company creation
// @ModuleID createCompany
// @Accept mpfd
// @Produce json
// @Param file formData file true "Image [jpeg/png]"
// @Param company formData string true "Look at the companyStringTemplate or entity.CompanyInput in Models"
// @Param companyStringTemplate body entity.CompanyInput false "Company creation template"
// @Success 201 {object} companyCreateResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /company/create [post]
func (h *Handler) createCompany(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	fileBody, fileType, err := h.getImageFromMultipartFormData(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	companyInput, err := h.getCompanyInputFromMultipartFormData(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	companyID, err := h.services.Companies.Create(userID, companyInput, fileBody, fileType)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, companyCreateResponse{
		CompanyID: companyID,
	})
}

// @Summary Company profile
// @Security UserAuth
// @Tags Company
// @Description company profile
// @ModuleID companyProfile
// @Accept json
// @Produce json
// @Success 200 {object} companyProfileResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /company/profile [get]
func (h *Handler) companyProfile(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	companyProfile, err := h.services.Companies.Profile(userID)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, companyProfileResponse{
		CompanyProfile: companyProfile,
	})
}
