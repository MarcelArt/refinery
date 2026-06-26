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

func renderWorkflowTemplate(c *gin.Context, status int, tabFile string, activeTab string, workflowTitle string, workflowDesc string, workflowID uint, user any, extra gin.H) {
	tmpl, err := template.ParseFiles(
		filepath.Join("internal", "web", "views", "layout.html"),
		filepath.Join("internal", "web", "views", "workflow_layout.html"),
		filepath.Join("internal", "web", "views", tabFile),
	)
	if err != nil {
		c.String(http.StatusInternalServerError, "Template parsing error: %s", err.Error())
		return
	}

	data := gin.H{
		"Title":               workflowTitle,
		"User":                user,
		"WorkflowID":          workflowID,
		"WorkflowTitle":       workflowTitle,
		"WorkflowDescription": workflowDesc,
		"ActiveTab":           activeTab,
		"ActiveMenu":          "workflows",
	}

	for k, v := range extra {
		data[k] = v
	}

	err = tmpl.ExecuteTemplate(c.Writer, "layout.html", data)
	if err != nil {
		c.String(http.StatusInternalServerError, "Template execution error: %s", err.Error())
	}
}
