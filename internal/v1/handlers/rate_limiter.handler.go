package handlers

import (
	"net/http"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/gin-gonic/gin"
	_ "github.com/morkid/paginate"
)

type RateLimiterHandler struct {
	service services.IRateLimiterService
}

func NewRateLimiterHandler(service services.IRateLimiterService) *RateLimiterHandler {
	return &RateLimiterHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create a new rate limiter
// @Description  Create a new rate limiter config with the provided details
// @Tags         rate-limiters
// @Accept       json
// @Produce      json
// @Param        rateLimiter  body      models.RateLimiterInput  true  "Rate limiter details"
// @Success      201   {object}  common.Result[uint]
// @Failure      400   {object}  common.Result[string]
// @Failure      401   {object}  common.Result[string]
// @Failure      500   {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/rate-limiters [post]
func (h *RateLimiterHandler) Create(c *gin.Context) {
	var rateLimiter models.RateLimiterInput
	if err := c.ShouldBindJSON(&rateLimiter); err != nil {
		_, res := common.ResultErr(err, "failed parsing json")
		c.JSON(http.StatusBadRequest, res)
		return
	}

	id, err := h.service.Create(c, rateLimiter)
	if err != nil {
		c.JSON(common.ResultErr(err, "failed creating rate limiter"))
		return
	}

	c.JSON(http.StatusCreated, common.ResultOk(id, "Rate limiter created"))
}

// Read godoc
// @Summary      List rate limiters
// @Description  Get a paginated list of rate limiters
// @Tags         rate-limiters
// @Accept       json
// @Produce      json
// @Param        page     query     int     false  "Page"
// @Param        size     query     int     false  "Size"
// @Param        sort     query     string  false  "Sort"
// @Param        filters  query     string  false  "Filter"
// @Success      200      {object}  paginate.Page{items=[]models.RateLimiterPage}
// @Failure      401      {object}  common.Result[string]
// @Failure      500      {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/rate-limiters [get]
func (h *RateLimiterHandler) Read(c *gin.Context) {
	rateLimiters, _ := h.service.Read(c)

	c.JSON(http.StatusOK, rateLimiters)
}

// Update godoc
// @Summary      Update rate limiter
// @Description  Update an existing rate limiter's details
// @Tags         rate-limiters
// @Accept       json
// @Produce      json
// @Param        id           path      string                   true  "Rate Limiter ID"
// @Param        rateLimiter  body      models.RateLimiterInput  true  "Updated rate limiter details"
// @Success      200   {object}  common.Result[any]
// @Failure      400   {object}  common.Result[string]
// @Failure      401   {object}  common.Result[string]
// @Failure      500   {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/rate-limiters/{id} [put]
func (h *RateLimiterHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var rateLimiter models.RateLimiterInput
	if err := c.ShouldBindJSON(&rateLimiter); err != nil {
		_, res := common.ResultErr(err, "failed parsing json")
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if err := h.service.Update(c, id, rateLimiter); err != nil {
		c.JSON(common.ResultErr(err, "failed updating rate limiter"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk[any](nil, "Rate limiter updated"))
}

// Delete godoc
// @Summary      Delete rate limiter
// @Description  Delete a rate limiter config by ID
// @Tags         rate-limiters
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Rate Limiter ID"
// @Success      200  {object}  common.Result[any]
// @Failure      401  {object}  common.Result[string]
// @Failure      500  {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/rate-limiters/{id} [delete]
func (h *RateLimiterHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c, c.Param("id")); err != nil {
		c.JSON(common.ResultErr(err, "failed deleting rate limiter"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk[any](nil, "Rate limiter deleted"))
}

// GetByID godoc
// @Summary      Get rate limiter by ID
// @Description  Get detailed information about a rate limiter config by its ID
// @Tags         rate-limiters
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Rate Limiter ID"
// @Success      200  {object}  common.Result[entities.RateLimiter]
// @Failure      401  {object}  common.Result[string]
// @Failure      500  {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/rate-limiters/{id} [get]
func (h *RateLimiterHandler) GetByID(c *gin.Context) {
	rateLimiter, err := h.service.GetByID(c, c.Param("id"))
	if err != nil {
		c.JSON(common.ResultErr(err, "failed getting rate limiter"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk(rateLimiter, "Rate limiter found"))
}
