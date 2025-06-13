package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/learn-gin/internal/repository"
)

type RedirectService struct {
	linkRepository repository.LinkRepository
}

func NewRedirectService(ctx *repository.SQLContext) *RedirectService {
	return &RedirectService{
		linkRepository: *repository.NewLinkRepository(ctx),
	}
}

func (rs *RedirectService) Redirect(g *gin.Context) {
	short := g.Param("link")

	link, err := rs.linkRepository.GetLink(short)
	if err == repository.ErrLinkNotFound {
		g.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	rs.linkRepository.IncrementVisits(short)
	g.Redirect(http.StatusFound, link.LongURL)
}
