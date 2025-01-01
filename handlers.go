package main

import (
	"errors"
	"fmj/internal/utils"
	"github.com/angelofallars/htmx-go"
	"html/template"
	"net/http"

	"path/filepath"

	gowebly "github.com/gowebly/helpers"

	"github.com/gin-gonic/gin"
)

// indexViewHandler handles a view for the index page.
func indexViewHandler(c *gin.Context) {
	var data map[string]interface{}
	isAuthenticated := c.GetBool("isAuthenticated")
	// Define paths to the user templates.
	indexPage := filepath.Join("templates", "pages", "index.html")

	data = map[string]interface{}{
		"isAuthenticated": isAuthenticated,
	}
	utils.Render(c, indexPage, data)

}

func showDashboardHandler(c *gin.Context) {
	// Define paths to the templates
	dashboardLayout := filepath.Join("templates", "dashboard.html")
	dashboardPage := filepath.Join("templates", "pages", "dashboard_home.html")

	// Parse templates
	tmpl, err := template.ParseFiles(dashboardLayout, dashboardPage)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Render the template
	if err := tmpl.Execute(c.Writer, nil); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

// showContentAPIHandler handles an API endpoint to show content.
//func showContentAPIHandler(c *gin.Context) {
//	// Check, if the current request has a 'HX-Request' header.
//	// For more information, see https://htmx.org/docs/#request-headers
//	if !htmx.IsHTMX(c.Request) {
//		// If not, return HTTP 400 error.
//		c.AbortWithError(http.StatusBadRequest, errors.New("non-htmx request"))
//		return
//	}
//
//	// Write HTML content.
//	c.Writer.Write([]byte("<p>🎉 Yes, <strong>htmx</strong> is ready to use! (<code>GET /api/hello-world</code>)</p>"))
//
//	// Send htmx response.
//	htmx.NewResponse().Write(c.Writer)
//}

func showContentAPIHandler(c *gin.Context) {
	// Check if the current request has an 'HX-Request' header.
	if !htmx.IsHTMX(c.Request) {
		// If not, return HTTP 400 error.
		c.AbortWithError(http.StatusBadRequest, errors.New("non-htmx request"))
		return
	}

	// Define the path to the content template.
	testage := filepath.Join("templates", "partials", "test.html")

	// Parse the content template using gowebly helper.
	tmpl, err := gowebly.ParseTemplates(testage)
	if err != nil {
		// If parsing fails, return HTTP 500 error.
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Prepare the data to pass to the template.
	data := map[string]interface{}{
		"Title":   "HTMX Content",
		"Message": "🎉 Yes, htmx is ready to use! (<code>GET /api/hello-world</code>)",
	}

	// Render the template with the data.
	if err := tmpl.Execute(c.Writer, data); err != nil {
		// If rendering fails, return HTTP 500 error.
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Send htmx response.
	htmx.NewResponse().Write(c.Writer)
}
