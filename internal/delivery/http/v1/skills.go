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

// @Summary Skills
// @Tags Skills
// @Description Get all skills
// @ModuleID getAllSkills
// @Produce  json
// @Success 200 {object} skillsResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /skills/all [get]
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
