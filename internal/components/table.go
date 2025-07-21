package components

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kevin-fagan/go-links/internal/htmx"
)

// TableLog retrieves paginated logs records from the repository
// and renders them into the "table_logs.html" template.
func (s *Service) TableLog(g *gin.Context) {
	renderTable(g, s.logs.ReadAll, "table_logs.html", "Logs")
}

// TableLink retrieves paginated link records from the repository
// and renders them into the "table_links.html" template.
func (s *Service) TableLink(g *gin.Context) {
	renderTable(g, s.links.ReadAll, "table_links.html", "Links")
}

// TableTag retrieves paginated link records from the repository
// and renders them into the "table_tags.html" template.
func (s *Service) TableTag(g *gin.Context) {
	renderTable(g, s.tags.ReadAll, "table_tags.html", "Tags")
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
