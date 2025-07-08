package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/repository"
)

const (
	formLong  = "long-url"
	formShort = "short-url"

	paramLink = "link"
)

// LinkService provides business logic for creating, updating, and deleting shortened links.
// It relies on a LinkRepository for persistent storage.
type LinkService struct {
	linkRepository repository.LinkRepository
}

// NewLinkService initializes a LinkService with a database context.
func NewLinkService(ctx *repository.SQLContext) *LinkService {
	return &LinkService{
		linkRepository: *repository.NewLinkRepository(ctx),
	}
}

// CreateLink handles HTTP POST requests to create a new shortened link.
// On error, returns an error modal, otherwise triggers a UI refresh.
func (ls *LinkService) CreateLink(g *gin.Context) {
	long := g.PostForm(formLong)
	short := g.PostForm(formShort)

	if short == "" || long == "" {
		triggerModalError(g, "missing short or long URL")
		return
	}

	err := ls.linkRepository.CreateLink(short, long, g.ClientIP())
	if err != nil {
		triggerModalError(g, err.Error())
		return
	}

	triggerRefresh(g)
}

// UpdateLink handles HTTP POST requests to update an existing shortened link.
// On error, returns an error modal, otherwise triggers a UI refresh.
func (ls *LinkService) UpdateLink(g *gin.Context) {
	long := g.PostForm(formLong)
	short := g.PostForm(formShort)

	if short == "" || long == "" {
		triggerModalError(g, "missing short or long URL")
		return
	}

	err := ls.linkRepository.UpdateLink(short, long, g.ClientIP())
	if err != nil {
		triggerModalError(g, err.Error())
		return
	}

	triggerRefresh(g)
}

// DeleteLink handles HTTP DELETE requests to update an existing shortened link.
// On error, returns an error modal, otherwise triggers a UI refresh.
func (ls *LinkService) DeleteLink(g *gin.Context) {
	short := g.Param(paramLink)

	err := ls.linkRepository.DeleteLink(short, g.ClientIP())
	if err != nil {
		triggerModalError(g, err.Error())
		return
	}

	triggerRefresh(g)
}

// triggerRefresh sends an HTMX trigger header to instruct the client to refresh UI components,
// then renders the modal-clear.html template to close any open modal dialogs.
func triggerRefresh(g *gin.Context) {
	g.Header("HX-Trigger", "refresh")
	g.HTML(http.StatusOK, "modal-clear.html", gin.H{})
}

// triggerModalError renders a modal-error.html template displaying the provided error message.
// Used to inform users of validation or repository errors.
func triggerModalError(g *gin.Context, message string) {
	g.HTML(http.StatusBadRequest, "modal-error.html", gin.H{
		"Message": message,
	})
}
