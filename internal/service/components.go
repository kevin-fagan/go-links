package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/repository"
)

type ComponentService struct {
	linkRepository repository.LinkRepository
}

func NewComponentService(ctx *repository.SQLContext) *ComponentService {
	return &ComponentService{
		linkRepository: *repository.NewLinkRepository(ctx),
	}
}

func (cs *ComponentService) LinkTable(g *gin.Context) {
	page, err := strconv.Atoi(g.Query("page"))
	if err != nil {
		g.String(http.StatusBadRequest, "invalid page")
		return
	}

	pageSize, err := strconv.Atoi(g.Query("pageSize"))
	if err != nil {
		g.String(http.StatusBadRequest, "invalid page size")
		return
	}

	links, err := cs.linkRepository.GetLinks(page, pageSize)
	if err != nil {
		g.String(http.StatusBadRequest, "error retrieving links")
		return
	}

	count, err := cs.linkRepository.Count()
	if err != nil {
		g.String(http.StatusBadRequest, "error retrieving links")
		return
	}

	g.HTML(http.StatusOK, "table.html", gin.H{
		"Links": &links,
		"Results": gin.H{
			"Current": min(count, (page+1)*pageSize),
			"Total":   count,
		},
		"Page": gin.H{
			"Size":     5,
			"Current":  page,
			"Previous": page - 1,
			"Next":     page + 1,
		},
	})
}

func (cs *ComponentService) FormCreate(g *gin.Context) {
	g.HTML(http.StatusOK, "create.html", gin.H{})
}

func (cs *ComponentService) FormClear(g *gin.Context) {
	g.HTML(http.StatusOK, "clear.html", gin.H{})
}

func (cs *ComponentService) FormUpdate(g *gin.Context) {
	short := g.Param("link")

	link, err := cs.linkRepository.GetLink(short)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.HTML(http.StatusOK, "update.html", &link)
}

func (cs *ComponentService) FormDelete(g *gin.Context) {
	short := g.Param("link")

	link, err := cs.linkRepository.GetLink(short)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.HTML(http.StatusOK, "delete.html", &link)
}
