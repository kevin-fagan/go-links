package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LinkPage(g *gin.Context) {
	g.HTML(http.StatusOK, "links.html", gin.H{})
}

func LogPage(g *gin.Context) {
	g.HTML(http.StatusOK, "logs.html", gin.H{})
}

func TagPage(g *gin.Context) {
	g.HTML(http.StatusOK, "tags.html", gin.H{})
}
