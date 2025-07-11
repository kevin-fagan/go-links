package main

import (
	"log"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/repository"
	"github.com/kevin-fagan/go-links/internal/routes"
	"github.com/kevin-fagan/go-links/internal/tmpl"
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

	router.SetFuncMap(template.FuncMap{
		"formatDate": tmpl.FormatDate,
		"formatChip": tmpl.FormatChip,
	})

	router.Static("/assets", "./web/assets")
	router.StaticFile("/favicon.ico", "./web/assets/images/favicon.ico")
	router.LoadHTMLGlob("web/html/**/*.html")

	root := router.Group("/")
	routes.AddPageRoute(root)
	routes.AddLinksRoutes(root, sqlite)
	routes.AddRedirectRoutes(root, sqlite)
	routes.AddComponentRoutes(root, sqlite)

	router.Run("localhost:8080")
}
