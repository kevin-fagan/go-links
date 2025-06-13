package routes

import "github.com/gin-gonic/gin"

func AddHomeRoutes(rg *gin.RouterGroup) {
	rg.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "home.html", gin.H{})
	})
}
