package handler

import (
	"github.com/PolyProjectOPD/Backend/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

// set of end-points
const (
	// POST
	auth   = "/auth"
	signUP = "/sign-up"
	signIn = "/signing"
)

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group(auth)
	{
		auth.POST(signUP, h.signUp)
		auth.POST(signIn, h.signIn)
	}

	return router
}
