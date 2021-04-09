package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userID"
)

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
