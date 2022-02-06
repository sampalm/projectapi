package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sampalm/projectapi/controllers"
	"github.com/sampalm/projectapi/middlewares"
)

func GetRoutes(router *gin.Engine) *gin.Engine {

	router.Use(middlewares.CORSMiddleware())

	auth := new(controllers.AuthController)
	router.GET("/api/v1/auth", auth.GetToken)

	v1 := router.Group("/api/v1", middlewares.Auth())
	{
		project := new(controllers.ProjectController)
		v1.GET("/project/:name", project.Show)
		v1.GET("/projects", project.All)
		v1.POST("/project", project.Create)
		v1.PUT("/project/:name", project.Update)
		v1.DELETE("/project/:name", project.Destroy)

		api := new(controllers.ApiController)
		v1.GET("/api/:project/:name", api.Show)
		v1.GET("/apis", api.All)
		v1.POST("/api", api.Create)
		v1.PUT("/api/:project/:name", api.Update)
		v1.DELETE("/api/:project/:name", api.Destroy)
	}

	return router
}
