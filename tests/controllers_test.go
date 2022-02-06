package tests

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sampalm/projectapi/controllers"
	"github.com/sampalm/projectapi/database"
	"github.com/sampalm/projectapi/database/migrations"
	"github.com/sampalm/projectapi/middlewares"
	"github.com/sampalm/projectapi/models"
	"github.com/stretchr/testify/assert"
)

var payloadProject = models.Project{
	Name:        "projetodeteste",
	DisplayName: "Projeto de Teste",
	Description: "Descricao de teste",
}

var payloadApi = models.API{
	Name:        "apideteste",
	Description: "Api de Teste",
	Version:     "1.0",
	ProjectName: "projetodeteste",
	OpenApiFile: MockFile(),
}

var token string

func MockFile() string {
	f, _ := os.Open("../openapi.yaml")

	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	encoded := "data:@file/x-yaml;base64," + base64.StdEncoding.EncodeToString(content)

	return encoded
}

func SetupRouter() *gin.Engine {

	r := gin.Default()
	gin.SetMode(gin.TestMode)

	auth := new(controllers.AuthController)
	r.GET("/api/v1/auth", auth.GetToken)

	v1 := r.Group("/api/v1", middlewares.Auth())
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

	return r

}

func TestInitDB(t *testing.T) {
	err := godotenv.Load("../test.env")
	if err != nil {
		log.Fatal("error: failed to load the env file")
	}
	database.StartDB()
	migrations.RunMigrations(database.GetDatabase())
}

func TestAuthorizationtokenGenerator(t *testing.T) {
	r := SetupRouter()

	req, err := http.NewRequest("GET", "/api/v1/auth", nil)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res struct {
		Token string `json:"token"`
	}

	json.Unmarshal(body, &res)

	token = res.Token

	r.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)

}

func TestCreateProjectFail(t *testing.T) {
	r := SetupRouter()

	paylod := models.Project{}

	data, _ := json.Marshal(paylod)

	req, err := http.NewRequest("POST", "/api/v1/project", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res models.Project
	json.Unmarshal(body, &res)

	assert.Equal(t, http.StatusNotAcceptable, resp.Code)
}

func TestCreateProjectPass(t *testing.T) {
	r := SetupRouter()

	data, _ := json.Marshal(payloadProject)

	req, err := http.NewRequest("POST", "/api/v1/project", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestGetProjectPass(t *testing.T) {
	r := SetupRouter()

	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/project/%s", payloadProject.Name), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

}

func TestGetAllProjectPass(t *testing.T) {
	r := SetupRouter()

	req, err := http.NewRequest("GET", "/api/v1/projects", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

}

func TestUpdateProjectPass(t *testing.T) {
	r := SetupRouter()

	payloadUpdate := models.Project{}
	payloadUpdate.DisplayName = "DisplayName Updated"
	payloadUpdate.Description = "Description Updated"

	data, _ := json.Marshal(payloadUpdate)

	req, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/project/%s", payloadProject.Name), bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res models.Project
	json.Unmarshal(body, &res)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestCreateApiPass(t *testing.T) {
	r := SetupRouter()

	data, _ := json.Marshal(payloadApi)

	req, err := http.NewRequest("POST", "/api/v1/api", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

}
func TestCreateApiFail(t *testing.T) {

	r := SetupRouter()

	data, _ := json.Marshal(payloadApi)

	req, err := http.NewRequest("POST", "/api/v1/api", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotAcceptable, resp.Code)

}
func TestGetApiPass(t *testing.T) {
	r := SetupRouter()

	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/api/%s/%s", payloadApi.ProjectName, payloadApi.Name), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

}

func TestUpdateApiPass(t *testing.T) {

	r := SetupRouter()

	payloadUpdate := models.API{
		Version:     "1.1",
		Description: "Description updated",
		OpenApiFile: MockFile(),
	}

	data, _ := json.Marshal(payloadUpdate)

	req, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/api/%s/%s", payloadApi.ProjectName, payloadApi.Name), bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

}

func TestGetAllApiPass(t *testing.T) {
	r := SetupRouter()

	req, err := http.NewRequest("GET", "/api/v1/apis", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
func TestDeleteApiPass(t *testing.T) {
	r := SetupRouter()

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/api/%s/%s", payloadApi.ProjectName, payloadApi.Name), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestDeleteProjectPass(t *testing.T) {
	r := SetupRouter()

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/project/%s", payloadProject.Name), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

}
