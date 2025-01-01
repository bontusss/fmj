package utils

import (
	"github.com/gin-gonic/gin"
	gowebly "github.com/gowebly/helpers"
	"log/slog"
	"net/http"
)

// Render encapsulates template rendering logic for handlers.
func Render(c *gin.Context, templatePath string, data interface{}) {
	tmpl, err := gowebly.ParseTemplates(templatePath)
	if err != nil {
		// Log error and return HTTP 400 error.
		slog.Error("Error parsing template", "path", templatePath, "error", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		// Log error and return HTTP 500 error.
		slog.Error("Error rendering template", "path", templatePath, "error", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
