package v1

import (
	"github.com/PolyProjectOPD/Backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.refresh)
	}
}

type signUpInput struct {
	Name     string `json:"name" binding:"required,min=2,max=64"`
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

func (h *Handler) signUp(c *gin.Context) {
	var input signUpInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Users.SignUp(service.UserSignUpInput{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.services.Users.SignIn(service.UserSignInInput{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
	})
}

type refreshInput struct {
	Token string `json:"token" binding:"required"`
}

func (h *Handler) refresh(c *gin.Context) {
	var input refreshInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.services.Users.RefreshTokens(input.Token)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
	})
}
