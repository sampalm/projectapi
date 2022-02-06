package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sampalm/projectapi/models"
	"gopkg.in/validator.v2"
)

type ApiController struct{}

var apiModel = new(models.APIModel)

func (ctrl ApiController) Show(c *gin.Context) {
	projectName := c.Param("project")
	apiName := c.Param("name")

	api, err := apiModel.Find(apiName, projectName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{
			"message": "Could not find API",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": api,
	})
}

func (ctrl ApiController) All(c *gin.Context) {

	apis, err := apiModel.All()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{
			"message": "No API was found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": apis,
	})
}

func (ctrl ApiController) Create(c *gin.Context) {

	apiObj := models.API{}
	if err := c.ShouldBindJSON(&apiObj); err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"message": "JSON is not valid ",
		})
		return
	}

	if err := validator.Validate(apiObj); err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"message": err.Error(),
		})
		return
	}

	err := apiModel.Save(apiObj.ProjectName, apiObj)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"message": "API could not be created",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": apiObj,
	})
}

func (ctrl ApiController) Update(c *gin.Context) {
	projectName := c.Param("project")
	apiName := c.Param("name")

	var apiObj struct {
		Version     string `json:"version" validate:"nonzero,max=10"`
		OpenApiFile string `json:"openapi_file" validate:"nonzero"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&apiObj); err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"message": "JSON is not valid ",
		})
		return
	}

	if err := validator.Validate(apiObj); err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"message": err.Error(),
		})
		return
	}

	payload := models.API{
		Version:     apiObj.Version,
		Description: apiObj.Description,
		OpenApiFile: apiObj.OpenApiFile,
	}

	err := apiModel.Update(apiName, projectName, payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"message": "Could not update, try again later => " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project updated",
	})
}

func (ctrl ApiController) Destroy(c *gin.Context) {
	projectName := c.Param("project")
	apiName := c.Param("name")

	err := apiModel.Destroy(apiName, projectName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"message": "Could not delete, try again later",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "API deleted",
	})
}
