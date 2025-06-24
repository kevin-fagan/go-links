package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/learn-gin/internal/repository"
)

type LinkService struct {
	linkRepository repository.LinkRepository
}

func NewLinkService(ctx *repository.SQLContext) *LinkService {
	return &LinkService{
		linkRepository: *repository.NewLinkRepository(ctx),
	}
}

func (ls *LinkService) CreateLink(g *gin.Context) {
	short := g.PostForm("short-url")
	long := g.PostForm("long-url")

	if short == "" || long == "" {
		g.String(http.StatusBadRequest, "missing short or long URL")
		return
	}

	err := ls.linkRepository.CreateLink(short, long)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.Header("hx-trigger", "refresh")
	g.HTML(http.StatusOK, "clear.html", gin.H{})
}

func (ls *LinkService) UpdateLink(g *gin.Context) {
	short := g.PostForm("short-url")
	long := g.PostForm("long-url")

	if short == "" || long == "" {
		g.String(http.StatusBadRequest, "missing short or long URL")
		return
	}

	err := ls.linkRepository.UpdateLink(short, long)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.Header("hx-trigger", "refresh")
	g.HTML(http.StatusOK, "clear.html", gin.H{})
}

func (ls *LinkService) DeleteLink(g *gin.Context) {
	short := g.Param("link")
	err := ls.linkRepository.DeleteLink(short)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.Header("hx-trigger", "refresh")
	g.HTML(http.StatusOK, "clear.html", gin.H{})
}
