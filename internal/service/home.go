package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeService struct{}

func NewHomeService() HomeService {
	return HomeService{}
}

func (hs *HomeService) Index(g *gin.Context) {
	g.HTML(http.StatusOK, "index.html", gin.H{})
}

func (hs *HomeService) LinkPage(g *gin.Context) {
	g.HTML(http.StatusOK, "page-links.html", gin.H{})
}

func (hs *HomeService) AuditPage(g *gin.Context) {
	g.HTML(http.StatusOK, "page-audit.html", gin.H{})
}
