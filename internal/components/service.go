package components

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/db"
	"github.com/kevin-fagan/go-links/internal/htmx"
	"github.com/kevin-fagan/go-links/internal/links"
	"github.com/kevin-fagan/go-links/internal/logs"
	"github.com/kevin-fagan/go-links/internal/tags"
)

type Service struct {
	tags  tags.Repository
	logs  logs.Repository
	links links.Repository
}

func NewService(ctx *db.SQLiteContext) *Service {
	return &Service{
		tags:  *tags.NewRepository(ctx),
		logs:  *logs.NewRepository(ctx),
		links: *links.NewRepository(ctx),
	}
}

// LogTable retrieves paginated logs records from the repository
// and renders them into the "table-logs.html" template.
func (s *Service) LogTable(g *gin.Context) {
	renderTable(g, s.logs.ReadAll, "table_logs.html", "Logs")
}

// LinkTable retrieves paginated link records from the repository
// and renders them into the "table-links.html" template.
func (s *Service) LinkTable(g *gin.Context) {
	renderTable(g, s.links.ReadAll, "table_links.html", "Links")
}

// TagTable retrieves paginated link records from the repository
// and renders them into the "table-links.html" template.
func (s *Service) TagTable(g *gin.Context) {
	renderTable(g, s.tags.ReadAll, "table_tags.html", "Tags")
}

// ModalCreate renders the modal used for creating a new link.
func (s *Service) ModalCreate(g *gin.Context) {
	g.HTML(http.StatusOK, "create.html", gin.H{})
}

// ModalClear renders an empty modal to clear HTML content via HTMX.
func (s *Service) ModalClear(g *gin.Context) {
	g.HTML(http.StatusOK, "clear.html", gin.H{})
}

// ModalUpdate renders the delete modal for a specific short link.
// Responds with 400 Bad Request if the link cannot be retrieved.
func (s *Service) ModalUpdate(g *gin.Context) {
	short := g.Param("link")

	link, err := s.links.Read(short)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.HTML(http.StatusOK, "update.html", &link)
}

// ModalDelete renders the delete modal for a specific short link.
// Responds with 400 Bad Request if the link cannot be retrieved.
func (s *Service) ModalDelete(g *gin.Context) {
	short := g.Param("link")

	link, err := s.links.Read(short)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.HTML(http.StatusOK, "delete.html", &link)
}

// renderTable handles paginated fetching and rendering of data using the given fetch function and template.
// It manages page bounds, extracts query parameters, and passes data plus pagination info to the template.
func renderTable[T any](g *gin.Context, fetchData func(page, pageSize int, search string) ([]T, int, error), templateName string, key string) {
	page, err := strconv.Atoi(g.Query("page"))
	if err != nil || page < 0 {
		page = 0
	}

	pageSize, err := strconv.Atoi(g.Query("pageSize"))
	if err != nil || pageSize <= 0 {
		pageSize = 25
	}

	search := g.Query("search")

	// Fetch total count (using a minimal fetch) to calculate page bounds
	_, totalCount, err := fetchData(0, 1, search)
	if err != nil {
		htmx.ModalError(g, err.Error())
		return
	}

	// Clamp page to maxPage based on total count to avoid empty pages
	maxPage := max(0, (totalCount-1)/pageSize)
	if page > maxPage {
		page = maxPage
	}

	// Fething the real page data
	data, count, err := fetchData(page, pageSize, search)
	if err != nil {
		htmx.ModalError(g, err.Error())
		return
	}

	// Calculate pagination result boundaries for UI
	start := 0
	if count != 0 {
		start = max(0, page*pageSize+1)
	}
	end := min(count, (page+1)*pageSize)

	// Prepare pagination info for template
	pageRange := gin.H{
		"Start": start,
		"End":   end,
		"Total": count,
	}

	pageInfo := gin.H{
		"Size":     pageSize,
		"Current":  page,
		"Previous": page - 1,
		"Next":     page + 1,
	}

	// Assemble template data and render
	dataMap := gin.H{
		key:       data,
		"Results": pageRange,
		"Page":    pageInfo,
	}

	g.HTML(http.StatusOK, templateName, dataMap)
}
