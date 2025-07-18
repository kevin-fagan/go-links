package components

import (
	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/db"
)

func AddRoutes(rg *gin.RouterGroup, sqlite *db.SQLiteContext) {
	s := NewService(sqlite)
	rg.GET("/components/table/links", s.LinkTable)
	rg.GET("/components/table/logs", s.LogTable)
	rg.GET("/components/table/tags", s.TagTable)

	rg.GET("/components/modal/clear", s.ModalClear)
	rg.GET("/components/modal/create", s.ModalCreate)
	rg.GET("/components/modal/update/:link", s.ModalUpdate)
	rg.GET("/components/modal/delete/:link", s.ModalDelete)
}
