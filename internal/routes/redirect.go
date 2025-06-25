package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/repository"
	"github.com/kevin-fagan/go-links/internal/service"
)

func AddRedirectRoutes(rg *gin.RouterGroup, sqlite *repository.SQLContext) {
	rs := service.NewRedirectService(sqlite)

	rg.GET("/:link", rs.Redirect)
}
