package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/learn-gin/internal/repository"
	"github.com/kevin-fagan/learn-gin/internal/service"
)

func AddLinksRoutes(rg *gin.RouterGroup, sqlite *repository.SQLContext) {
	ls := service.NewLinkService(sqlite)

	rg.POST("/links", ls.CreateLink)
	rg.POST("/links/:link", ls.UpdateLink)
	rg.DELETE("/links/:link", ls.DeleteLink)
}
