package service

import "github.com/gin-gonic/gin"

type HomeService struct{}

func NewHomeService() HomeService {
	return HomeService{}
}

func (hs *HomeService) LinkPage(g *gin.Context) {
	g.HTML(200, "page-links.html", gin.H{})
}

func (hs *HomeService) AuditPage(g *gin.Context) {
	g.HTML(200, "page-audit.html", gin.H{})
}
