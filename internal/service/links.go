package service

import (
	"fmt"
	"net/http"
	"strconv"

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

func (ls *LinkService) RedirectToLink(g *gin.Context) {
	short := g.Param("link")

	link, err := ls.linkRepository.GetLink(short)
	if err == repository.ErrLinkNotFound {
		g.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// Increment the visits field for a url. We can safely ignore
	// any errors that occur for this function call
	ls.linkRepository.IncrementVisits(short)

	// Redirecting the user to the long url
	g.Redirect(http.StatusFound, link.LongName)
}

func (ls *LinkService) GetLinks(g *gin.Context) {
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

	links, err := ls.linkRepository.GetLinks(page, pageSize)
	if err != nil {
		g.String(http.StatusBadRequest, "error retrieving links")
		return
	}

	fmt.Println(links)

	g.HTML(http.StatusOK, "table.html", links)
}
