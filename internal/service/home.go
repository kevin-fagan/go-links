package service

import "github.com/gin-gonic/gin"

type HomeService struct{}

func NewHomeService() HomeService {
	return HomeService{}
}

func (hs *HomeService) HomePage(g *gin.Context) {
	g.HTML(200, "home.html", gin.H{})
}
