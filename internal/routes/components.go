package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/learn-gin/internal/repository"
	"github.com/kevin-fagan/learn-gin/internal/service"
)

func AddComponentRoutes(rg *gin.RouterGroup, sqlite *repository.SQLContext) {
	cs := service.NewComponentService(sqlite)

	rg.GET("/components/table", cs.LinkTable)
	rg.GET("/components/form/create", cs.FormCreate)
	rg.GET("/components/form/update/:link", cs.FormUpdate)
	rg.GET("/components/form/delete/:link", cs.FormDelete)
}
