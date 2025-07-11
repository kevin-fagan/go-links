package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PageService struct{}

func NewPageService() PageService {
	return PageService{}
}

func (hs *PageService) LinkPage(g *gin.Context) {
	g.HTML(http.StatusOK, "page-links.html", gin.H{})
}

func (hs *PageService) AuditPage(g *gin.Context) {
	g.HTML(http.StatusOK, "page-audit.html", gin.H{})
}
