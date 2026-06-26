package handlers

import (
	"log"
	"net/http"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/gin-gonic/gin"
	_ "github.com/morkid/paginate"
)

type WorkflowHandler struct {
	service  services.IWorkflowService
	nService services.IN8NService
}

func NewWorkflowHandler(service services.IWorkflowService, nService services.IN8NService) *WorkflowHandler {
	return &WorkflowHandler{
		service:  service,
		nService: nService,
	}
}

// Create godoc
// @Summary      Create a new workflow
// @Description  Create a new workflow with the provided details
// @Tags         workflows
// @Accept       json
// @Produce      json
// @Param        workflow  body      models.WorkflowInput  true  "Workflow details"
// @Success      201   {object}  common.Result[uint]
// @Failure      400   {object}  common.Result[string]
// @Failure      401   {object}  common.Result[string]
// @Failure      500   {object}  common.Result[string]
// @Security     ApiKeyAuth
// @Router       /v1/workflows [post]
func (h *WorkflowHandler) Create(c *gin.Context) {
	userID, err := common.MustGet[float64](c, "userId")
	if err != nil {
		_, res := common.ResultErr(err, "token invalid")
		c.JSON(http.StatusUnauthorized, res)
		return
	}
	log.Println("userID :>> ", userID)

	var workflow models.WorkflowInput
	if err := c.ShouldBindJSON(&workflow); err != nil {
		_, res := common.ResultErr(err, "failed parsing json")
		c.JSON(http.StatusBadRequest, res)
		return
	}
	workflow.UserID = uint(userID)
	log.Println("workflow :>> ", workflow)

	id, err := h.service.Create(c, workflow)
	if err != nil {
		c.JSON(common.ResultErr(err, "failed creating workflow"))
		return
	}

	c.JSON(http.StatusCreated, common.ResultOk(id, "Workflow created"))
}

// Read godoc
// @Summary      List workflows
// @Description  Get a paginated list of workflows
// @Tags         workflows
// @Accept       json
// @Produce      json
// @Param        page     query     int     false  "Page"
// @Param        size     query     int     false  "Size"
// @Param        sort     query     string  false  "Sort"
// @Param        filters  query     string  false  "Filter"
// @Success      200      {object}  paginate.Page{items=[]models.WorkflowPage}
// @Failure      401      {object}  common.Result[string]
// @Failure      500      {object}  common.Result[string]
// @Security     ApiKeyAuth
// @Router       /v1/workflows [get]
func (h *WorkflowHandler) Read(c *gin.Context) {
	workflows, _ := h.service.Read(c)

	c.JSON(http.StatusOK, workflows)
}

// Update godoc
// @Summary      Update workflow
// @Description  Update an existing workflow's details
// @Tags         workflows
// @Accept       json
// @Produce      json
// @Param        id        path      string                true  "Workflow ID"
// @Param        workflow  body      models.WorkflowInput  true  "Updated workflow details"
// @Success      200   {object}  common.Result[any]
// @Failure      400   {object}  common.Result[string]
// @Failure      401   {object}  common.Result[string]
// @Failure      500   {object}  common.Result[string]
// @Security     ApiKeyAuth
// @Router       /v1/workflows/{id} [put]
func (h *WorkflowHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var workflow models.WorkflowInput
	if err := c.ShouldBindJSON(&workflow); err != nil {
		_, res := common.ResultErr(err, "failed parsing json")
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if err := h.service.Update(c, id, workflow); err != nil {
		c.JSON(common.ResultErr(err, "failed updating workflow"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk[any](nil, "Workflow updated"))
}

// Delete godoc
// @Summary      Delete workflow
// @Description  Delete a workflow by ID
// @Tags         workflows
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Workflow ID"
// @Success      200  {object}  common.Result[any]
// @Failure      401  {object}  common.Result[string]
// @Failure      500  {object}  common.Result[string]
// @Security     ApiKeyAuth
// @Router       /v1/workflows/{id} [delete]
func (h *WorkflowHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c, c.Param("id")); err != nil {
		c.JSON(common.ResultErr(err, "failed deleting workflow"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk[any](nil, "Workflow deleted"))
}

// GetByID godoc
// @Summary      Get workflow by ID
// @Description  Get detailed information about a workflow by its ID
// @Tags         workflows
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Workflow ID"
// @Success      200  {object}  common.Result[entities.Workflow]
// @Failure      401  {object}  common.Result[string]
// @Failure      500  {object}  common.Result[string]
// @Security     ApiKeyAuth
// @Router       /v1/workflows/{id} [get]
func (h *WorkflowHandler) GetByID(c *gin.Context) {
	workflow, err := h.service.GetByID(c, c.Param("id"))
	if err != nil {
		c.JSON(common.ResultErr(err, "failed getting workflow"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk(workflow, "Workflow found"))
}

// Upload godoc
// @Summary      Upload file to workflow
// @Description  Upload a file to be processed by a workflow by ID
// @Tags         workflows
// @Accept       multipart/form-data
// @Produce      json
// @Param        id    path      string  true  "Workflow ID"
// @Param        file  formData  file    true  "File to upload"
// @Param        file  formData  string    false  "Additional prompt to apply during runtime"
// @Success      200   {object}  common.Result[any]
// @Failure      400   {object}  common.Result[string]
// @Failure      401   {object}  common.Result[string]
// @Failure      500   {object}  common.Result[string]
// @Security     ApiKeyAuth
// @Router       /v1/workflows/{id}/upload [post]
func (h *WorkflowHandler) Upload(c *gin.Context) {
	id := c.Param("id")
	workflowOption := models.WorkflowStartOption{
		AdditionalPrompt: c.PostForm("additionalPrompt"),
	}

	formFile, err := c.FormFile("file")
	if err != nil {
		_, res := common.ResultErr(err, "failed uploading file")
		c.JSON(http.StatusBadRequest, res)
		return
	}
	file, err := formFile.Open()
	if err != nil {
		c.JSON(common.ResultErr(err, "failed to open file"))
		return
	}
	defer file.Close()

	if err := h.service.UploadToWorkflow(c, id, formFile.Filename, file, workflowOption); err != nil {
		c.JSON(common.ResultErr(err, "failed upload to workflow"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk[any](nil, "Workflow started"))
}
