package pages

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(rg *gin.RouterGroup) {
	rg.GET("/home", LinkPage)
	rg.GET("/home/links", LinkPage)
	rg.GET("/home/logs", LogPage)
	rg.GET("/home/tags", TagPage)
}
