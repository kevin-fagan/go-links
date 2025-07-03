package main

import (
	"log"
	"text/template"
	"time"

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

	router.SetFuncMap(template.FuncMap{
		"formatDate": formatDate,
	})

	router.Static("/assets", "./web/assets")
	router.StaticFile("/favicon.ico", "./web/assets/images/favicon.ico")

	router.LoadHTMLFiles(
		"web/html/home.html",
		"web/html/modal/clear.html",
		"web/html/modal/create.html",
		"web/html/modal/delete.html",
		"web/html/modal/error.html",
		"web/html/modal/update.html",
		"web/html/table/links.html",
		"web/html/table/audit.html",
	)

	root := router.Group("/")
	routes.AddHomeRoutes(root)
	routes.AddLinksRoutes(root, sqlite)
	routes.AddRedirectRoutes(root, sqlite)
	routes.AddComponentRoutes(root, sqlite)

	router.Run("localhost:8080")
}

func formatDate(t time.Time) string {
	return t.Format("02 Jan 2006 15:04")
}
