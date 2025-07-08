package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/repository"
	"github.com/kevin-fagan/go-links/internal/service"
)

func AddComponentRoutes(rg *gin.RouterGroup, sqlite *repository.SQLContext) {
	cs := service.NewComponentService(sqlite)

	rg.GET("/components/table/links", cs.LinkTable)
	rg.GET("/components/table/audit", cs.AuditTable)

	rg.GET("/components/modal/clear", cs.ModalClear)
	rg.GET("/components/modal/create", cs.ModalCreate)
	rg.GET("/components/modal/update/:link", cs.ModalUpdate)
	rg.GET("/components/modal/delete/:link", cs.ModalDelete)
}
