package authorization

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	serverstatus "github.com/qwuiemme/ellipsespace-server/internal/server-status"
)

func ParseJWTFromHeader(c *gin.Context) SessionBase {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return SessionBase{}
	}

	headerParts := strings.Split(authHeader, " ")

	if len(headerParts) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return SessionBase{}
	}

	if headerParts[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return SessionBase{}
	}

	sb, err := ParseJWT(headerParts[1])

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, serverstatus.StatusJson{
			Message: err.Error(),
		})
		return SessionBase{}
	}

	return sb
}

func AuthorizationRequired(c *gin.Context) {
	ParseJWTFromHeader(c)
}

func AdminAccessLevelRequired(c *gin.Context) {
	data := ParseJWTFromHeader(c)

	if data.AccessLevel != 1 {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
}
