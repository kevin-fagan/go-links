package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/learn-gin/internal/repository"
	"github.com/kevin-fagan/learn-gin/internal/routes"
)

var (
	sqlite *repository.SQLContext
)

func init() {
	var err error

	sqlite, err = repository.Connect("compass-db")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("web/**")

	v1 := router.Group("/v1")
	routes.AddLinkRoutes(v1, sqlite)

	router.Run("localhost:8080")
}
