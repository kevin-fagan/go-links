package htmx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Refresh sends an HTMX trigger header to instruct the client to refresh UI components,
// then renders the modal-clear.html template to close any open modal dialogs.
func Refresh(g *gin.Context) {
	g.Header("HX-Trigger", "refresh")
	g.HTML(http.StatusOK, "clear.html", gin.H{})
}

// ModalError renders a modal-error.html template displaying the provided error message.
// Used to inform users of validation or repository errors.
func ModalError(g *gin.Context, message string) {
	g.HTML(http.StatusBadRequest, "error.html", gin.H{
		"Message": message,
	})
}
