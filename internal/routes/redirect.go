package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/learn-gin/internal/repository"
)

func AddRedirectRoutes(rg *gin.RouterGroup, sqlite *repository.SQLContext) {
	rg.GET("/:link")
}
