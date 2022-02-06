package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sampalm/projectapi/services"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		const Bearer_scheme = "Bearer "
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		token := header[len(Bearer_scheme):]

		if !services.NewJWTService().ValidateToken(token) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
