package links

import (
	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/db"
)

func AddRoutes(rg *gin.RouterGroup, sqlite *db.SQLiteContext) {
	s := NewService(sqlite)

	// Handles redirect from short to long URL
	rg.GET("/:link", s.Redirect)

	rg.POST("/links", s.Create)
	rg.POST("/links/:link", s.Update)
	rg.DELETE("/links/:link", s.Delete)
}
