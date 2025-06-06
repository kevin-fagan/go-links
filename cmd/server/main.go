package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/learn-gin/internal/repository"
	"github.com/kevin-fagan/learn-gin/internal/service"
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
	ls := service.NewLinkService(sqlite)

	router := gin.Default()
	router.LoadHTMLGlob("web/**")

	// Route(s) to handle home page
	home := router.Group("/")
	{
		home.GET("/", func(ctx *gin.Context) {
			ctx.HTML(200, "home.html", gin.H{})
		})
	}

	// Route(s) to handle HTMX components
	components := router.Group("/components")
	{
		components.GET("/table", ls.GetLinks)
		components.GET("/form/create")
		components.GET("/form/edit/:link")

	}

	// Route(s) to handle link CRUD operations
	links := router.Group("/links")
	{
		links.POST("/")
		links.PUT("/:link")
		links.DELETE("/:link")
	}

	// Route(s) to handle redirect
	redirect := router.Group("/redirect")
	{
		redirect.GET("/:link")
	}

	router.Run("localhost:8080")
}
