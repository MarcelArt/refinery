package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/MarcelArt/refinery/internal/web/viewmodels"
	"github.com/gin-gonic/gin"
)

type ExtractionResultWebHandler struct {
	extractionResultService services.IExtractionResultService
	workflowService         services.IWorkflowService
	userService             services.IUserService
}

func NewExtractionResultWebHandler(
	extractionResultService services.IExtractionResultService,
	workflowService services.IWorkflowService,
	userService services.IUserService,
) *ExtractionResultWebHandler {
	return &ExtractionResultWebHandler{
		extractionResultService: extractionResultService,
		workflowService:         workflowService,
		userService:             userService,
	}
}

// ShowResults renders the extraction results layout for a workflow
func (h *ExtractionResultWebHandler) ShowResults(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	user, err := h.userService.GetByID(c, userId)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	workflowIDStr := c.Param("id")
	workflow, err := h.workflowService.GetByID(c, workflowIDStr)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/workflows")
		return
	}

	schemas, _ := workflow.Schemas.Deserialize()

	// Fetch paginated extraction results for the workflow
	pageInfo, pages := h.extractionResultService.GetByWorkflowID(c, workflow.ID)

	resultsVM := make([]viewmodels.ExtractionResultRowViewModel, 0, len(pages))
	for _, p := range pages {
		resultsVM = append(resultsVM, viewmodels.ExtractionResultRowViewModel{
			ID:        p.ID,
			CreatedAt: p.CreatedAt,
		})
	}

	var selectedResultVM *viewmodels.ExtractionResultDetailsViewModel
	resultIDStr := c.Param("resultId")

	if resultIDStr != "" {
		res, err := h.extractionResultService.GetByID(c, resultIDStr)
		if err == nil && res.WorkflowID == workflow.ID {
			jsonArray, _ := res.Json.Deserialize()
			selectedResultVM = &viewmodels.ExtractionResultDetailsViewModel{
				ID:        res.ID,
				CreatedAt: res.CreatedAt,
				Columns:   extractUniqueKeys(jsonArray),
				Rows:      jsonArray,
			}
		}
	} else if len(pages) > 0 {
		// Default to the first result in the list
		firstResult := pages[0]
		jsonArray, _ := firstResult.Json.Deserialize()
		selectedResultVM = &viewmodels.ExtractionResultDetailsViewModel{
			ID:        firstResult.ID,
			CreatedAt: firstResult.CreatedAt,
			Columns:   extractUniqueKeys(jsonArray),
			Rows:      jsonArray,
		}
	}

	// Compute item bounds
	start := pageInfo.Page*pageInfo.Size + 1
	end := start + int64(len(resultsVM)) - 1
	if len(resultsVM) == 0 {
		start = 0
		end = 0
	}

	paginationVM := viewmodels.PaginationViewModel{
		Total:      pageInfo.Total,
		Page:       pageInfo.Page,
		Size:       pageInfo.Size,
		TotalPages: pageInfo.TotalPages,
		Last:       pageInfo.Last,
		First:      pageInfo.First,
		PrevPage:   pageInfo.Page - 1,
		NextPage:   pageInfo.Page + 1,
		Start:      start,
		End:        end,
	}

	renderTemplate(c, http.StatusOK, "extraction_results.html", gin.H{
		"Title":           "Extraction Results",
		"User":            user,
		"WorkflowID":      workflow.ID,
		"WorkflowTitle":   workflow.Title,
		"WorkflowPrompt":  workflow.Prompt,
		"WorkflowSchemas": schemas,
		"Results":         resultsVM,
		"SelectedResult":  selectedResultVM,
		"Pagination":      paginationVM,
		"ActiveMenu":      "workflows", // Highlights workflows menu
	})
}

// ShowResultDetails renders the right-hand details table as an HTML fragment
func (h *ExtractionResultWebHandler) ShowResultDetails(c *gin.Context) {
	resultIDStr := c.Param("resultId")
	res, err := h.extractionResultService.GetByID(c, resultIDStr)
	if err != nil {
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{"Error": "Result not found"})
		return
	}

	jsonArray, _ := res.Json.Deserialize()
	selectedResultVM := &viewmodels.ExtractionResultDetailsViewModel{
		ID:        res.ID,
		CreatedAt: res.CreatedAt,
		Columns:   extractUniqueKeys(jsonArray),
		Rows:      jsonArray,
	}

	// Parse extraction_results.html to execute the result_details block directly
	tmpl, err := template.ParseFiles(filepath.Join("internal", "web", "views", "extraction_results.html"))
	if err != nil {
		c.String(http.StatusInternalServerError, "Template parsing error: %s", err.Error())
		return
	}

	err = tmpl.ExecuteTemplate(c.Writer, "result_details", selectedResultVM)
	if err != nil {
		c.String(http.StatusInternalServerError, "Template execution error: %s", err.Error())
	}
}

// extractUniqueKeys gathers all distinct keys in the JSON records slice to represent dynamic columns
func extractUniqueKeys(jsonArray []map[string]any) []string {
	keysMap := make(map[string]bool)
	var keys []string
	for _, item := range jsonArray {
		for k := range item {
			if !keysMap[k] {
				keysMap[k] = true
				keys = append(keys, k)
			}
		}
	}
	return keys
}

// Upload processes the file upload to start a new extraction run for the workflow
func (h *ExtractionResultWebHandler) Upload(c *gin.Context) {
	_, exists := c.Get("userId")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	workflowIDStr := c.Param("id")
	formFile, err := c.FormFile("file")
	if err != nil {
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
			"Error": "Failed uploading file: " + err.Error(),
		})
		return
	}

	file, err := formFile.Open()
	if err != nil {
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
			"Error": "Failed to open file: " + err.Error(),
		})
		return
	}
	defer file.Close()

	if err := h.workflowService.UploadToWorkflow(c, workflowIDStr, formFile.Filename, file); err != nil {
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
			"Error": "Failed upload to workflow: " + err.Error(),
		})
		return
	}

	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", "/workflows/"+workflowIDStr+"/results")
		c.Status(http.StatusOK)
	} else {
		c.Redirect(http.StatusSeeOther, "/workflows/"+workflowIDStr+"/results")
	}
}
