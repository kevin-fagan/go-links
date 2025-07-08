package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/repository"
)

// RedirectService handles the logic for redirecting short URLs
// to their corresponding long URLs. It uses LinkRepository
// to fetch link data and track visits.
type RedirectService struct {
	linkRepository repository.LinkRepository
}

// NewRedirectService initializes a RedirectService with the given database context.
func NewRedirectService(ctx *repository.SQLContext) *RedirectService {
	return &RedirectService{
		linkRepository: *repository.NewLinkRepository(ctx),
	}
}

// Redirect resolves a short URL to its long URL,
// counts the visit, and issues a 302 redirect.
func (rs *RedirectService) Redirect(g *gin.Context) {
	short := g.Param(paramLink)

	link, err := rs.linkRepository.GetLink(short)
	if err == repository.ErrLinkNotFound {
		g.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	rs.linkRepository.CountLinkVisit(short)
	g.Redirect(http.StatusFound, link.LongURL)
}
