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
	router.LoadHTMLGlob("web/**/*.html")
	router.Static("/assets", "./web/assets")

	home := router.Group("/")

	routes.AddHomeRoutes(home)
	routes.AddLinksRoutes(home, sqlite)
	routes.AddRedirectRoutes(home, sqlite)
	routes.AddComponentRoutes(home, sqlite)

	router.Run("localhost:8080")
}
