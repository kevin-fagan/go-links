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
		g.HTML(http.StatusBadRequest, "error.html", "unable to retrieve links")
		return
	}

	start := 0
	if count != 0 {
		start = max(0, page*pageSize+1)
	}

	end := min(count, (page+1)*pageSize)

	g.HTML(http.StatusOK, "links.html", gin.H{
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
