package v1

import (
	"github.com/PolyProjectOPD/Backend/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initSpheresRoutes(api *gin.RouterGroup) {
	spheres := api.Group("/spheres")
	{
		spheres.GET("/all", h.getAllSpheres)
		spheres.POST("/skills", h.getSkills)
	}
}

// @Summary Spheres
// @Tags Spheres
// @Description Get all spheres
// @ModuleID getAll
// @Produce  json
// @Success 200 {object} spheresResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /spheres/all [get]
func (h *Handler) getAllSpheres(c *gin.Context) {
	spheres, err := h.services.Spheres.GetAll()
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, spheresResponse{
		Spheres: spheres,
	})
}

// @Summary Spheres
// @Tags Spheres
// @Description Get all skills from sphere
// @ModuleID getSkills
// @Accept json
// @Produce  json
// @Param input body entity.SpheresInput true "Spheres input info"
// @Success 200 {object} skillsResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /spheres/skills [post]
func (h *Handler) getSkills(c *gin.Context) {
	var input entity.SpheresInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	skills, err := h.services.Spheres.GetSkills(input.Spheres)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, skillsResponse{
		Skills: skills,
	})
}
