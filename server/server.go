package server

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sampalm/projectapi/routes"
)

type Server struct {
	port   string
	engine *gin.Engine
}

func NewServer() Server {
	return Server{
		port:   os.Getenv("PORT"),
		engine: gin.Default(),
	}
}

func (s *Server) Run() {
	router := routes.GetRoutes(s.engine)

	log.Fatal(router.Run(":" + s.port))
}
