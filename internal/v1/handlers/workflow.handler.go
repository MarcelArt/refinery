package handlers

import (
	"net/http"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/entities"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/gin-gonic/gin"
	_ "github.com/morkid/paginate"
)

type WorkflowHandler struct {
	service services.IWorkflowService
}

var _ = entities.Workflow{}

func NewWorkflowHandler(service services.IWorkflowService) *WorkflowHandler {
	return &WorkflowHandler{
		service: service,
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
	var workflow models.WorkflowInput
	if err := c.ShouldBindJSON(&workflow); err != nil {
		_, res := common.ResultErr(err, "failed parsing json")
		c.JSON(http.StatusBadRequest, res)
		return
	}

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
