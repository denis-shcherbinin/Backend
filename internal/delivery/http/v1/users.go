package v1

import (
	"github.com/PolyProjectOPD/Backend/internal/entity"
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

// @Summary User SignUp
// @Tags User Auth
// @Description User sign-up
// @ModuleID signUp
// @Accept  json
// @Produce  json
// @Param input body entity.UserSignUpInput true "Sign-up info"
// @Success 201 {object} signUpResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input entity.UserSignUpInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Users.SignUp(input)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, signUpResponse{
		ID: id,
	})
}

// @Summary User SignIn
// @Tags User Auth
// @Description User sign-in
// @ModuleID signIn
// @Accept  json
// @Produce  json
// @Param input body entity.UserSignInInput true "Sign-in info"
// @Success 200 {object} tokensResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {

	var input entity.UserSignInInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.services.Users.SignIn(input)

	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// @Summary User Refresh Tokens
// @Tags User Auth
// @Description User refresh tokens
// @Accept  json
// @Produce  json
// @Param input body entity.UserRefreshInput true "Refresh tokens info"
// @Success 200 {object} tokensResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/refresh [post]
func (h *Handler) refresh(c *gin.Context) {
	var input entity.UserRefreshInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.services.Users.RefreshTokens(input)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
