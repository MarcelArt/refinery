package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func renderTemplate(c *gin.Context, status int, templateFile string, data gin.H) {
	// Parse layout.html and the specific page template
	tmpl, err := template.ParseFiles(
		filepath.Join("internal", "web", "views", "layout.html"),
		filepath.Join("internal", "web", "views", templateFile),
	)
	if err != nil {
		c.String(http.StatusInternalServerError, "Template parsing error: %s", err.Error())
		return
	}

	err = tmpl.ExecuteTemplate(c.Writer, "layout.html", data)
	if err != nil {
		c.String(http.StatusInternalServerError, "Template execution error: %s", err.Error())
	}
}

func renderFragment(c *gin.Context, status int, templateFile string, data gin.H) {
	tmpl, err := template.ParseFiles(filepath.Join("internal", "web", "views", templateFile))
	if err != nil {
		c.String(http.StatusInternalServerError, "Template parsing error: %s", err.Error())
		return
	}

	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		c.String(http.StatusInternalServerError, "Template execution error: %s", err.Error())
	}
}
