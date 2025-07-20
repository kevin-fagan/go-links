package components

import (
	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/db"
)

func AddRoutes(rg *gin.RouterGroup, sqlite *db.SQLiteContext) {
	s := NewService(sqlite)
	rg.GET("/components/table/links", s.TableLink)
	rg.GET("/components/table/logs", s.TableLog)
	rg.GET("/components/table/tags", s.TableTag)

	rg.GET("/components/modal/clear", s.ModalClear)

	rg.GET("/components/modal/links/create", s.ModalCreateLink)
	rg.GET("/components/modal/links/update/:link", s.ModaleUpdateLink)
	rg.GET("/components/modal/links/delete/:link", s.ModalDeleteLink)

	rg.GET("/components/modal/tags/create", s.ModalCreateTag)
}
