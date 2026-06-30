package handlers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/MarcelArt/refinery/internal/v1/middlewares"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/MarcelArt/refinery/internal/web/viewmodels"
	"github.com/gin-gonic/gin"
)

type customResponseWriter struct {
	gin.ResponseWriter
	body    *bytes.Buffer
	status  int
	headers http.Header
}

func (w *customResponseWriter) Header() http.Header {
	return w.headers
}

func (w *customResponseWriter) WriteHeader(status int) {
	w.status = status
}

func (w *customResponseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *customResponseWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}

type ExtractionResultWebHandler struct {
	extractionResultService services.IExtractionResultService
	workflowService         services.IWorkflowService
	userService             services.IUserService
	rateLimitM              *middlewares.RateLimiterMiddleware
}

func NewExtractionResultWebHandler(
	extractionResultService services.IExtractionResultService,
	workflowService services.IWorkflowService,
	userService services.IUserService,
	rateLimitM *middlewares.RateLimiterMiddleware,
) *ExtractionResultWebHandler {
	return &ExtractionResultWebHandler{
		extractionResultService: extractionResultService,
		workflowService:         workflowService,
		userService:             userService,
		rateLimitM:              rateLimitM,
	}
}

// RateLimit acts as a delegate to the backend RateLimiterMiddleware
func (h *ExtractionResultWebHandler) RateLimit(c *gin.Context) {
	realWriter := c.Writer
	blw := &customResponseWriter{
		ResponseWriter: realWriter,
		body:           bytes.NewBuffer(nil),
		status:         http.StatusOK,
		headers:        make(http.Header),
	}
	c.Writer = blw

	h.rateLimitM.Limit(c)

	// Restore the real writer
	c.Writer = realWriter

	if c.IsAborted() {
		// Output the error HTML alert fragment instead of JSON
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
			"Error": "Daily extraction quota exceeded. Please try again tomorrow.",
		})
		return
	}

	// If not aborted, forward the headers and flush the buffered body
	for k, vv := range blw.headers {
		for _, v := range vv {
			realWriter.Header().Add(k, v)
		}
	}
	if blw.status != http.StatusOK {
		realWriter.WriteHeader(blw.status)
	}
	if blw.body.Len() > 0 {
		_, _ = realWriter.Write(blw.body.Bytes())
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
			Status:    p.Status,
		})
	}

	var selectedResultVM *viewmodels.ExtractionResultDetailsViewModel
	resultIDStr := c.Param("resultId")

	if resultIDStr != "" {
		res, err := h.extractionResultService.GetByID(c, resultIDStr)
		if err == nil && res.WorkflowID == workflow.ID {
			jsonArray, _ := res.Json.Deserialize()
			selectedResultVM = &viewmodels.ExtractionResultDetailsViewModel{
				ID:           res.ID,
				CreatedAt:    res.CreatedAt,
				Status:       res.Status,
				FinishedAt:   res.FinishedAt,
				Columns:      extractUniqueKeys(jsonArray),
				Rows:         jsonArray,
				Attachment:   res.Attachment,
				WorkflowType: workflow.Type,
			}
		}
	} else if len(pages) > 0 {
		// Default to the first result in the list
		firstResult := pages[0]
		jsonArray, _ := firstResult.Json.Deserialize()
		selectedResultVM = &viewmodels.ExtractionResultDetailsViewModel{
			ID:           firstResult.ID,
			CreatedAt:    firstResult.CreatedAt,
			Status:       firstResult.Status,
			FinishedAt:   firstResult.FinishedAt,
			Columns:      extractUniqueKeys(jsonArray),
			Rows:         jsonArray,
			Attachment:   firstResult.Attachment,
			WorkflowType: firstResult.WorkflowType,
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

	renderWorkflowTemplate(c, http.StatusOK, "extraction_results.html", "results", workflow.Title, workflow.Description, workflow.ID, user, gin.H{
		"WorkflowPrompt":  workflow.Prompt,
		"WorkflowSchemas": schemas,
		"Results":         resultsVM,
		"SelectedResult":  selectedResultVM,
		"Pagination":      paginationVM,
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

	workflow, err := h.workflowService.GetByID(c, res.WorkflowID)
	if err != nil {
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{"Error": "Workflow not found"})
		return
	}

	jsonArray, _ := res.Json.Deserialize()
	selectedResultVM := &viewmodels.ExtractionResultDetailsViewModel{
		ID:           res.ID,
		CreatedAt:    res.CreatedAt,
		Status:       res.Status,
		FinishedAt:   res.FinishedAt,
		Columns:      extractUniqueKeys(jsonArray),
		Rows:         jsonArray,
		Attachment:   res.Attachment,
		WorkflowType: workflow.Type,
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

	additionalPrompt := c.PostForm("additionalPrompt")
	metadata := c.PostForm("metadata")

	if metadata != "" {
		var temp json.RawMessage
		if err := json.Unmarshal([]byte(metadata), &temp); err != nil {
			renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
				"Error": "Metadata must be a valid JSON format: " + err.Error(),
			})
			return
		}
	}

	workflowOption := models.WorkflowStartOption{
		AdditionalPrompt: additionalPrompt,
		Metadata:         metadata,
	}

	if err := h.workflowService.UploadToWorkflow(c, workflowIDStr, formFile.Filename, file, workflowOption); err != nil {
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
			"Error": "Failed upload to workflow: " + err.Error(),
		})
		return
	}

	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Trigger", "workflow-started")
		c.Status(http.StatusOK)
	} else {
		c.Redirect(http.StatusSeeOther, "/workflows/"+workflowIDStr+"/results")
	}
}
