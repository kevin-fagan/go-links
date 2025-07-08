package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/repository"
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
		g.HTML(http.StatusBadRequest, "modal-error.html", gin.H{
			"Message": "missing short or long URL",
		})
		return
	}

	err := ls.linkRepository.CreateLink(short, long, g.ClientIP())
	if err != nil {
		g.HTML(http.StatusBadRequest, "modal-error.html", gin.H{
			"Message": err.Error(),
		})
		return
	}

	g.Header("HX-Trigger", "refresh")
	g.HTML(http.StatusOK, "clear.html", gin.H{})
}

func (ls *LinkService) UpdateLink(g *gin.Context) {
	short := g.PostForm("short-url")
	long := g.PostForm("long-url")

	if short == "" || long == "" {
		g.HTML(http.StatusBadRequest, "modal-error.html", gin.H{
			"Message": "missing short or long URL",
		})
		return
	}

	err := ls.linkRepository.UpdateLink(short, long, g.ClientIP())
	if err != nil {
		g.HTML(http.StatusBadRequest, "modal-error.html", gin.H{
			"Message": err.Error(),
		})
		return
	}

	g.Header("HX-Trigger", "refresh")
	g.HTML(http.StatusOK, "clear.html", gin.H{})
}

func (ls *LinkService) DeleteLink(g *gin.Context) {
	short := g.Param("link")
	err := ls.linkRepository.DeleteLink(short, g.ClientIP())

	if err != nil {
		g.HTML(http.StatusBadRequest, "modal-error.html", gin.H{
			"Message": err.Error(),
		})
		return
	}

	g.Header("HX-Trigger", "refresh")
	g.HTML(http.StatusOK, "clear.html", gin.H{})
}
