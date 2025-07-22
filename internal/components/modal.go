package components

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ModalCreateLink renders the modal used for creating a new link.
func (s *Service) ModalCreateLink(g *gin.Context) {
	g.HTML(http.StatusOK, "link_create.html", gin.H{})
}

// ModaleUpdateLink renders the delete modal for a specific short link.
// Responds with 400 Bad Request if the link cannot be retrieved.
func (s *Service) ModaleUpdateLink(g *gin.Context) {
	short := g.Param("link")

	link, err := s.links.Read(short)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.HTML(http.StatusOK, "link_update.html", &link)
}

// ModalDeleteLink renders the delete modal for a specific short link.
// Responds with 400 Bad Request if the link cannot be retrieved.
func (s *Service) ModalDeleteLink(g *gin.Context) {
	short := g.Param("link")

	link, err := s.links.Read(short)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.HTML(http.StatusOK, "link_delete.html", &link)
}

// ModalDeleteTag renders the delete modal for a specific short link.
// Responds with 400 Bad Request if the link cannot be retrieved.
func (s *Service) ModalDeleteTag(g *gin.Context) {
	t := g.Param("tag")

	tag, err := s.tags.Read(t)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.HTML(http.StatusOK, "tag_delete.html", &tag)
}

// ModalUpdateTag renders the delete modal for a specific short link.
// Responds with 400 Bad Request if the link cannot be retrieved.
func (s *Service) ModalUpdateTag(g *gin.Context) {
	t := g.Param("tag")

	tag, err := s.tags.Read(t)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.HTML(http.StatusOK, "tag_update.html", &tag)
}

func (s *Service) ModalCreateTag(g *gin.Context) {
	g.HTML(http.StatusOK, "tag_create.html", gin.H{})
}

// ModalClear renders an empty modal to clear HTML content via HTMX.
func (s *Service) ModalClear(g *gin.Context) {
	g.HTML(http.StatusOK, "clear.html", gin.H{})
}
