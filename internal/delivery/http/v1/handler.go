package v1

import (
	"github.com/PolyProjectOPD/Backend/internal/service"
	"github.com/PolyProjectOPD/Backend/pkg/auth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUsersRoutes(v1)
		h.initSpheresRoutes(v1)
		h.initSkillsRoutes(v1)
		h.initCompaniesRoutes(v1)
	}
}
