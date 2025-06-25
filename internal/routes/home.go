package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/service"
)

func AddHomeRoutes(rg *gin.RouterGroup) {
	hs := service.NewHomeService()

	rg.GET("/home", hs.HomePage)
}
