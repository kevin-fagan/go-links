package main

import (
	"log"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/components"
	"github.com/kevin-fagan/go-links/internal/db"
	"github.com/kevin-fagan/go-links/internal/links"
	"github.com/kevin-fagan/go-links/internal/pages"
	"github.com/kevin-fagan/go-links/internal/tags"
	"github.com/kevin-fagan/go-links/internal/tmpl"
)

var (
	sqlite *db.SQLiteContext
)

func init() {
	var err error

	sqlite, err = db.Connect("golinks-db")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	router := gin.Default()

	// Setting up Go templating functions
	router.SetFuncMap(template.FuncMap{
		"FormatDate": tmpl.FormatDate,
		"FormatChip": tmpl.FormatChip,
	})

	// Loading static assets (HTML, CSS, JS)
	router.Static("/assets", "./web/assets")
	router.StaticFile("/favicon.ico", "./web/assets/images/logo.png")
	router.LoadHTMLGlob("web/html/**/*.html")

	// Setting up routes (URL paths)
	root := router.Group("/")
	pages.AddRoutes(root)
	links.AddRoutes(root, sqlite)
	tags.AddRoutes(root, sqlite)
	components.AddRoutes(root, sqlite)

	// Starting Gin
	router.Run("localhost:8080")
}
