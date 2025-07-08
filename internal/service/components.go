package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/repository"
)

// ComponentService provides HTTP handlers for rendering modals and tables
// related to links and audit records. It depends on a LinkRepository for data access.
type ComponentService struct {
	linkRepository repository.LinkRepository
}

// NewComponentService creates and returns a new ComponentService using the provided SQL context.
func NewComponentService(ctx *repository.SQLContext) *ComponentService {
	return &ComponentService{
		linkRepository: *repository.NewLinkRepository(ctx),
	}
}

// TODO: Simplify
func (cs *ComponentService) AuditTable(g *gin.Context) {
	page, err := strconv.Atoi(g.DefaultQuery("page", "0"))
	if err != nil || page < 0 {
		page = 0
	}

	pageSize, err := strconv.Atoi(g.DefaultQuery("pageSize", "5"))
	if err != nil || pageSize < 0 {
		pageSize = 5
	}

	search := g.Query("search")
	audits, count, err := cs.linkRepository.GetAudit(page, pageSize, search)
	if err != nil {
		g.HTML(http.StatusBadRequest, "modal-error.html", "unable to retrieve audits")
		return
	}

	start := 0
	if count != 0 {
		start = max(0, page*pageSize+1)
	}

	end := min(count, (page+1)*pageSize)

	g.HTML(http.StatusOK, "table-audit.html", gin.H{
		"Audit": &audits,
		"Results": gin.H{
			"Start": start,
			"End":   end,
			"Total": count,
		},
		"Page": gin.H{
			"Size":     5,
			"Current":  page,
			"Previous": page - 1,
			"Next":     page + 1,
		},
	})
}

// TODO: Simplify
func (cs *ComponentService) LinkTable(g *gin.Context) {
	page, err := strconv.Atoi(g.DefaultQuery("page", "0"))
	if err != nil || page < 0 {
		page = 0
	}

	pageSize, err := strconv.Atoi(g.DefaultQuery("pageSize", "5"))
	if err != nil || pageSize < 0 {
		pageSize = 5
	}

	search := g.Query("search")
	links, count, err := cs.linkRepository.GetLinks(search, page, pageSize)
	if err != nil {
		g.HTML(http.StatusBadRequest, "modal-error.html", "unable to retrieve links")
		return
	}

	start := 0
	if count != 0 {
		start = max(0, page*pageSize+1)
	}

	end := min(count, (page+1)*pageSize)

	g.HTML(http.StatusOK, "table-links.html", gin.H{
		"Links": &links,
		"Results": gin.H{
			"Start": start,
			"End":   end,
			"Total": count,
		},
		"Page": gin.H{
			"Size":     5,
			"Current":  page,
			"Previous": page - 1,
			"Next":     page + 1,
		},
	})
}

// ModalCreate renders the modal used for creating a new link.
func (cs *ComponentService) ModalCreate(g *gin.Context) {
	g.HTML(http.StatusOK, "modal-create.html", gin.H{})
}

// ModalClear renders an empty modal to clear HTML content via HTMX.
func (cs *ComponentService) ModalClear(g *gin.Context) {
	g.HTML(http.StatusOK, "modal-clear.html", gin.H{})
}

// ModalUpdate renders the delete modal for a specific short link.
// Responds with 400 Bad Request if the link cannot be retrieved.
func (cs *ComponentService) ModalUpdate(g *gin.Context) {
	short := g.Param(paramLink)

	link, err := cs.linkRepository.GetLink(short)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.HTML(http.StatusOK, "modal-update.html", &link)
}

// ModalDelete renders the delete modal for a specific short link.
// Responds with 400 Bad Request if the link cannot be retrieved.
func (cs *ComponentService) ModalDelete(g *gin.Context) {
	short := g.Param(paramLink)

	link, err := cs.linkRepository.GetLink(short)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.HTML(http.StatusOK, "modal-delete.html", &link)
}
