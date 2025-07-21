package tags

import (
	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/db"
)

func AddRoutes(rg *gin.RouterGroup, sqlite *db.SQLiteContext) {
	s := NewService(sqlite)

	rg.POST("/tags", s.Create)
	rg.POST("/tags/:tag", s.Update)
	rg.DELETE("/tags/:tag", s.Delete)
}
