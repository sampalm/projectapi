package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sampalm/projectapi/models"
	"gopkg.in/validator.v2"
)

type ProjectController struct{}

var projectModel = new(models.ProjectModel)

func (ctrl ProjectController) Show(c *gin.Context) {

	projectName := c.Param("name")

	project, err := projectModel.Find(projectName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{
			"message": "Project not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": project,
	})
}

func (ctrl ProjectController) All(c *gin.Context) {

	projects, err := projectModel.All()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{
			"message": "No project was found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": projects,
	})
}

func (ctrl ProjectController) Create(c *gin.Context) {

	project := models.Project{}

	if err := c.ShouldBindJSON(&project); err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"message": "JSON is not valid ",
		})
		return
	}

	if err := validator.Validate(project); err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"message": err.Error(),
		})
		return
	}

	err := projectModel.Save(project)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"message": "Project could not be created",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": project,
	})
}

func (ctrl ProjectController) Update(c *gin.Context) {
	projectName := c.Param("name")

	var projectObj struct {
		DisplayName string `json:"display_name" validate:"min=3,max=40"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&projectObj); err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"message": "JSON is not valid ",
		})
		return
	}

	if err := validator.Validate(projectObj); err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"message": err.Error(),
		})
		return
	}

	project := models.Project{
		DisplayName: projectObj.DisplayName,
		Description: projectObj.Description,
	}

	err := projectModel.Update(projectName, project)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"message": "Could not update, try again later",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project updated",
	})
}

func (ctrl ProjectController) Destroy(c *gin.Context) {

	projectName := c.Param("name")

	err := projectModel.Destroy(projectName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"message": "Could not delete, try again later",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project deleted",
	})
}
