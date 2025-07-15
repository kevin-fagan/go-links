package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LinkPage(g *gin.Context) {
	g.HTML(http.StatusOK, "page-links.html", gin.H{})
}

func LogPage(g *gin.Context) {
	g.HTML(http.StatusOK, "page-logs.html", gin.H{})
}

func TagPage(g *gin.Context) {
	g.HTML(http.StatusOK, "page-tags.html", gin.H{})
}
