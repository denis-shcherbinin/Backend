package v1

import (
	"github.com/PolyProjectOPD/Backend/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	user := api.Group("/user")
	{
		auth := user.Group("/auth")
		{
			auth.POST("/sign-up", h.signUp)
			auth.POST("/sign-in", h.signIn)
			auth.POST("/refresh", h.refresh)
		}

		authenticated := user.Group("/", h.userIdentify)
		{
			authenticated.GET("/logout", h.logout)
			authenticated.GET("/sign-out", h.signOut)
		}

		user.POST("/existence", h.isExists)
	}
}

// @Summary User SignUp
// @Tags User auth
// @Description User sign-up
// @ModuleID signUp
// @Accept mpfd
// @Produce json
// @Param input body entity.UserSignUpInput true "Sign-up info"
// @Success 201 {object} signUpResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /user/auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	// Image Upload
	fileBody, fileType, err := h.getImageFromMultipartFormData(c)
	if err != nil {
		if err.Error() != "http: no such file" {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	signUpInput, err := h.getUserSignUpInputFromMultipartFormData(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, imageURL, err := h.services.Users.SignUp(signUpInput, fileBody, fileType)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, signUpResponse{
		ID:       id,
		ImageURL: imageURL,
	})
}

// @Summary User SignIn
// @Tags User auth
// @Description User sign-in
// @ModuleID signIn
// @Accept json
// @Produce json
// @Param input body entity.UserSignInInput true "Sign-in info"
// @Success 200 {object} tokensResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /user/auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input entity.UserSignInInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userAgent, err := h.getUserAgent(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.services.Users.SignIn(input, userAgent)

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
// @Tags User auth
// @Description User refresh tokens
// @Accept json
// @Produce json
// @Param input body entity.UserRefreshInput true "Refresh tokens info"
// @Success 200 {object} tokensResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /user/auth/refresh [post]
func (h *Handler) refresh(c *gin.Context) {
	var input entity.UserRefreshInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userAgent, err := h.getUserAgent(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.services.Users.RefreshTokens(input, userAgent)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// @Summary User logout
// @Security UserAuth
// @Tags User
// @Description user sign out from all devices
// @ModuleID logout
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /user/logout [get]
func (h *Handler) logout(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.Users.Logout(userID); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

// @Summary User signOut
// @Security UserAuth
// @Tags User
// @Description user sign out from current device
// @ModuleID signOut
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /user/sign-out [get]
func (h *Handler) signOut(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userAgent, err := h.getUserAgent(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Users.SignOut(userID, userAgent)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.Status(http.StatusOK)
}

// @Summary User existence
// @Tags User
// @Description User existence
// @Accept json
// @Produce json
// @Param input body entity.UserExistenceInput true "User existence info"
// @Success 200 {object} userExistenceResponse
// @Failure 400,404 {object} response
// @Router /user/existence [post]
func (h *Handler) isExists(c *gin.Context) {
	var input entity.UserExistenceInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, userExistenceResponse{
		Exists: h.services.Users.Existence(input),
	})
}
