package links

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/db"
	"github.com/kevin-fagan/go-links/internal/htmx"
)

type Service struct {
	repository Repository
}

func NewService(ctx *db.SQLiteContext) *Service {
	return &Service{repository: *NewRepository(ctx)}
}

// Create handles HTTP POST requests to create a new shortened link.
// On error, returns an error modal, otherwise triggers a UI refresh.
func (s *Service) Create(g *gin.Context) {
	long := g.PostForm("long-url")
	short := g.PostForm("short-url")

	if short == "" || long == "" {
		htmx.ModalError(g, "missing short or long url")
		return
	}

	err := s.repository.Create(short, long, g.ClientIP())
	if err != nil {
		htmx.ModalError(g, err.Error())
		return
	}

	htmx.Refresh(g)
}

// Update handles HTTP POST requests to update an existing shortened link.
// On error, returns an error modal, otherwise triggers a UI refresh.
func (s *Service) Update(g *gin.Context) {
	long := g.PostForm("long-url")
	short := g.PostForm("short-url")

	if short == "" || long == "" {
		htmx.ModalError(g, "missing short or long url")
		return
	}

	err := s.repository.Update(short, long, g.ClientIP())
	if err != nil {
		htmx.ModalError(g, err.Error())
		return
	}

	htmx.Refresh(g)
}

// Delete handles HTTP DELETE requests to update an existing shortened link.
// On error, returns an error modal, otherwise triggers a UI refresh.
func (s *Service) Delete(g *gin.Context) {
	short := g.Param("link")

	err := s.repository.Delete(short, g.ClientIP())
	if err != nil {
		htmx.ModalError(g, err.Error())
		return
	}

	htmx.Refresh(g)
}

func (s *Service) Redirect(g *gin.Context) {
	short := g.Param("link")

	link, err := s.repository.Read(short)
	if err == ErrLinkNotFound {
		g.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	s.repository.CountVisit(short)
	g.Redirect(http.StatusFound, link.LongURL)
}
