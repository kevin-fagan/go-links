package tags

import (
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

func (s *Service) Create(g *gin.Context) {
	tag := g.PostForm("tag")

	if tag == "" {
		htmx.ModalError(g, "missing tag")
		return
	}

	err := s.repository.Create(tag, g.ClientIP())
	if err != nil {
		htmx.ModalError(g, err.Error())
		return
	}

	htmx.Refresh(g)
}

func (s *Service) Update(g *gin.Context) {
	oldTag := g.Param("tag")
	newTag := g.PostForm("tag")

	if oldTag == newTag {
		return
	}

	if oldTag == "" {
		htmx.ModalError(g, "url param 'tag' missing")
		return
	}

	if newTag == "" {
		htmx.ModalError(g, "form data 'tag' missing")
		return
	}

	err := s.repository.Update(oldTag, newTag, g.ClientIP())
	if err != nil {
		htmx.ModalError(g, err.Error())
		return
	}

	htmx.Refresh(g)
}

func (s *Service) Delete(g *gin.Context) {
	tag := g.Param("tag")

	err := s.repository.Delete(tag, g.ClientIP())
	if err != nil {
		htmx.ModalError(g, err.Error())
		return
	}

	htmx.Refresh(g)
}
