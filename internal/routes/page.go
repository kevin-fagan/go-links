package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/service"
)

func AddPageRoute(rg *gin.RouterGroup) {
	hs := service.NewPageService()

	rg.GET("/home", hs.LinkPage)
	rg.GET("/home/links", hs.LinkPage)
	rg.GET("/home/audit", hs.AuditPage)
}
