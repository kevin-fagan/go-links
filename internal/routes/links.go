package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/learn-gin/internal/repository"
)

func AddLinkRoutes(rg *gin.RouterGroup, ctx *repository.SQLContext) {
	links := rg.Group("/links")

	links.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{})
	})

	links.GET("/:link")
	links.PUT("/:link")
	links.POST("/:link")
	links.DELETE("/:link")
}
