package v1

import (
	"encoding/json"
	"errors"
	"github.com/PolyProjectOPD/Backend/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userID"
)

var imageTypes = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
}

func (h *Handler) userIdentify(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(userCtx, id)
}

func (h *Handler) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return "", errors.New("empty authorization header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	return h.tokenManager.Parse(headerParts[1])
}

func (h *Handler) getUserID(c *gin.Context) (int, error) {
	idFromContext, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("userCtx not found")
	}

	idStr, ok := idFromContext.(string)
	if !ok {
		return 0, errors.New("userCtx is of invalid type")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (h *Handler) getUserAgent(c *gin.Context) (string, error) {
	userAgent := c.GetHeader("User-Agent")
	if userAgent == "" {
		return "", errors.New("empty user-agent header")
	}

	return userAgent, nil
}

func (h *Handler) getImageFromMultipartFormData(c *gin.Context) (string, string, error) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	fileBody := make([]byte, fileHeader.Size)
	_, err = file.Read(fileBody)
	if err != nil {
		return "", "", err
	}
	fileType := http.DetectContentType(fileBody)

	// Validate File Type
	if _, ok := imageTypes[fileType]; !ok {
		return "", "", errors.New("file type isn't supported")
	}

	return string(fileBody), fileType, nil
}

func (h *Handler) getUserSignUpInputFromMultipartFormData(c *gin.Context) (entity.UserSignUpInput, error) {
	formValue := c.Request.PostFormValue("user")
	var input entity.UserSignUpInput
	err := json.Unmarshal([]byte(formValue), &input)
	if err != nil {
		return input, err
	}

	if len(input.FirstName) < 2 || len(input.FirstName) > 64 {
		return input, errors.New("invalid firstName")
	}

	if len(input.LastName) < 2 || len(input.LastName) > 64 {
		return input, errors.New("invalid lastName")
	}

	if len(input.BirthDate) != 10 {
		return input, errors.New("invalid birthDate")
	}

	if len(input.Email) == 0 {
		return input, errors.New("invalid email")
	}

	if len(input.Password) < 8 {
		return input, errors.New("invalid password")
	}

	return input, nil
}

func (h *Handler) getProfileInputFromMultipartFormData(c *gin.Context) (entity.ProfileInput, error) {
	formValue := c.Request.PostFormValue("profile")

	var input entity.ProfileInput

	if err := json.Unmarshal([]byte(formValue), &input); err != nil {
		return input, err
	}

	if len(input.FirstName) < 2 || len(input.FirstName) > 64 {
		return input, errors.New("invalid firstName")
	}

	if len(input.LastName) < 2 || len(input.LastName) > 64 {
		return input, errors.New("invalid lastName")
	}

	if len(input.BirthDate) != 10 {
		return input, errors.New("invalid birthDate")
	}

	if len(input.Email) == 0 {
		return input, errors.New("invalid email")
	}

	return input, nil
}

func (h *Handler) getCompanyInputFromMultipartFormData(c *gin.Context) (entity.CompanyInput, error) {
	formValue := c.Request.PostFormValue("company")

	var input entity.CompanyInput

	if err := json.Unmarshal([]byte(formValue), &input); err != nil {
		return input, err
	}

	return input, nil
}

func (h *Handler) getCompanyProfileFromMultipartFormData(c *gin.Context) (entity.CompanyProfile, error) {
	formValue := c.Request.PostFormValue("companyProfile")

	var input entity.CompanyProfile

	if err := json.Unmarshal([]byte(formValue), &input); err != nil {
		return input, err
	}

	return input, nil
}
