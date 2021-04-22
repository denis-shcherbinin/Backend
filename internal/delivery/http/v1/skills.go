package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initSkillsRoutes(api *gin.RouterGroup) {
	skills := api.Group("/skills")
	{
		skills.GET("/all", h.getAllSkills)
	}
}

func (h *Handler) getAllSkills(c *gin.Context) {
	skills, err := h.services.Skills.GetAll()
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, skillsResponse{
		Skills: skills,
	})
}
