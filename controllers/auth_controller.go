package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sampalm/projectapi/services"
)

type AuthController struct{}

func (ctrl AuthController) GetToken(c *gin.Context) {
	token, err := services.NewJWTService().GenerateToken()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Could not generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (ctrl AuthController) RefreshToken() {

}
