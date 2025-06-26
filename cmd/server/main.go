package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/repository"
	"github.com/kevin-fagan/go-links/internal/routes"
)

var (
	sqlite *repository.SQLContext
)

func init() {
	var err error

	sqlite, err = repository.Connect("golinks-db")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("web/**/*.html")
	router.Static("/assets", "./web/assets")

	root := router.Group("/")

	routes.AddHomeRoutes(root)
	routes.AddLinksRoutes(root, sqlite)
	routes.AddRedirectRoutes(root, sqlite)
	routes.AddComponentRoutes(root, sqlite)

	router.Run("localhost:8080")
}
