package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/learn-gin/internal/service"
)

func AddHomeRoutes(rg *gin.RouterGroup) {
	hs := service.NewHomeService()

	rg.GET("/home", hs.HomePage)
}
