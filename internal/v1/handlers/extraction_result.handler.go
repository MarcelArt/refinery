package handlers

import (
	"net/http"
	"strconv"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/gin-gonic/gin"
	_ "github.com/morkid/paginate"
)

type ExtractionResultHandler struct {
	service services.IExtractionResultService
}

func NewExtractionResultHandler(service services.IExtractionResultService) *ExtractionResultHandler {
	return &ExtractionResultHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create a new extraction result
// @Description  Create a new extraction result with the provided details
// @Tags         extraction-results
// @Accept       json
// @Produce      json
// @Param        extractionResult  body      models.ExtractionResultInput  true  "Extraction result details"
// @Success      201   {object}  common.Result[uint]
// @Failure      400   {object}  common.Result[string]
// @Failure      401   {object}  common.Result[string]
// @Failure      500   {object}  common.Result[string]
// @Security     ApiKeyAuth
// @Router       /v1/extraction-results [post]
func (h *ExtractionResultHandler) Create(c *gin.Context) {
	var extractionResult models.ExtractionResultInput
	if err := c.ShouldBindJSON(&extractionResult); err != nil {
		_, res := common.ResultErr(err, "failed parsing json")
		c.JSON(http.StatusBadRequest, res)
		return
	}

	id, err := h.service.Create(c, extractionResult)
	if err != nil {
		c.JSON(common.ResultErr(err, "failed creating extraction result"))
		return
	}

	c.JSON(http.StatusCreated, common.ResultOk(id, "Extraction result created"))
}

// Read godoc
// @Summary      List extraction results
// @Description  Get a paginated list of extraction results
// @Tags         extraction-results
// @Accept       json
// @Produce      json
// @Param        page     query     int     false  "Page"
// @Param        size     query     int     false  "Size"
// @Param        sort     query     string  false  "Sort"
// @Param        filters  query     string  false  "Filter"
// @Success      200      {object}  paginate.Page{items=[]models.ExtractionResultPage}
// @Failure      401      {object}  common.Result[string]
// @Failure      500      {object}  common.Result[string]
// @Security     ApiKeyAuth
// @Router       /v1/extraction-results [get]
func (h *ExtractionResultHandler) Read(c *gin.Context) {
	extractionResults, _ := h.service.Read(c)

	c.JSON(http.StatusOK, extractionResults)
}

// Update godoc
// @Summary      Update extraction result
// @Description  Update an existing extraction result's details
// @Tags         extraction-results
// @Accept       json
// @Produce      json
// @Param        id                path      string                        true  "Extraction Result ID"
// @Param        extractionResult  body      models.ExtractionResultInput  true  "Updated extraction result details"
// @Success      200   {object}  common.Result[any]
// @Failure      400   {object}  common.Result[string]
// @Failure      401   {object}  common.Result[string]
// @Failure      500   {object}  common.Result[string]
// @Security     ApiKeyAuth
// @Router       /v1/extraction-results/{id} [put]
func (h *ExtractionResultHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var extractionResult models.ExtractionResultInput
	if err := c.ShouldBindJSON(&extractionResult); err != nil {
		_, res := common.ResultErr(err, "failed parsing json")
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if err := h.service.Update(c, id, extractionResult); err != nil {
		c.JSON(common.ResultErr(err, "failed updating extraction result"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk[any](nil, "Extraction result updated"))
}

// Delete godoc
// @Summary      Delete extraction result
// @Description  Delete an extraction result by ID
// @Tags         extraction-results
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Extraction Result ID"
// @Success      200  {object}  common.Result[any]
// @Failure      401  {object}  common.Result[string]
// @Failure      500  {object}  common.Result[string]
// @Security     ApiKeyAuth
// @Router       /v1/extraction-results/{id} [delete]
func (h *ExtractionResultHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c, c.Param("id")); err != nil {
		c.JSON(common.ResultErr(err, "failed deleting extraction result"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk[any](nil, "Extraction result deleted"))
}

// GetByID godoc
// @Summary      Get extraction result by ID
// @Description  Get detailed information about an extraction result by its ID
// @Tags         extraction-results
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Extraction Result ID"
// @Success      200  {object}  common.Result[entities.ExtractionResult]
// @Failure      401  {object}  common.Result[string]
// @Failure      500  {object}  common.Result[string]
// @Security     ApiKeyAuth
// @Router       /v1/extraction-results/{id} [get]
func (h *ExtractionResultHandler) GetByID(c *gin.Context) {
	extractionResult, err := h.service.GetByID(c, c.Param("id"))
	if err != nil {
		c.JSON(common.ResultErr(err, "failed getting extraction result"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk(extractionResult, "Extraction result found"))
}

// Webhook godoc
// @Summary      Webhook to save extraction result from LLM
// @Description  Save an extraction result using data received from the LLM webhook
// @Tags         extraction-results
// @Accept       json
// @Produce      json
// @Param        id     path      string             true  "Workflow ID"
// @Param        input  body      models.ContentLLM  true  "Content from LLM"
// @Success      201    {object}  common.Result[uint]
// @Failure      400    {object}  common.Result[string]
// @Failure      500    {object}  common.Result[string]
// @Router       /v1/extraction-results/{id}/webhook [post]
func (h *ExtractionResultHandler) Webhook(c *gin.Context) {
	id := c.Param("id")
	workflowID, err := strconv.Atoi(id)
	if err != nil {
		_, res := common.ResultErr(err, "invalid workflow id")
		c.JSON(http.StatusBadRequest, res)
		return
	}

	var input models.ContentLLM
	if err := c.ShouldBindJSON(&input); err != nil {
		_, res := common.ResultErr(err, "failed parsing json")
		c.JSON(http.StatusBadRequest, res)
		return
	}

	erID, err := h.service.SaveFromLLM(c, uint(workflowID), input)
	if err != nil {
		c.JSON(common.ResultErr(err, "failed saving extraction result"))
		return
	}

	c.JSON(http.StatusCreated, common.ResultOk(erID, "Extraction result saved"))
}
