package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/repository"
	"github.com/kevin-fagan/go-links/internal/service"
)

func AddComponentRoutes(rg *gin.RouterGroup, sqlite *repository.SQLContext) {
	cs := service.NewComponentService(sqlite)

	rg.GET("/components/table/links", cs.LinkTable)
	rg.GET("/components/form/clear", cs.FormClear)
	rg.GET("/components/form/create", cs.FormCreate)
	rg.GET("/components/form/update/:link", cs.FormUpdate)
	rg.GET("/components/form/delete/:link", cs.FormDelete)
}
