package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/entities"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/repositories"
	"github.com/MarcelArt/refinery/pkg/jsonb"
	"github.com/gin-gonic/gin"
	"github.com/morkid/paginate"
)

type IExtractionResultService interface {
	common.IBaseCrudService[entities.ExtractionResult, models.ExtractionResultInput, models.ExtractionResultPage]
	SaveFromLLM(c context.Context, id uint, input models.ContentLLM) (uint, error)
}

type ExtractionResultService struct {
	repo repositories.IExtractionResultRepo
}

var _ IExtractionResultService = &ExtractionResultService{}

func NewExtractionResultService(repo repositories.IExtractionResultRepo) *ExtractionResultService {
	return &ExtractionResultService{
		repo: repo,
	}
}

func (s *ExtractionResultService) Create(c context.Context, input models.ExtractionResultInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *ExtractionResultService) Read(c *gin.Context) (paginate.Page, []models.ExtractionResultPage) {
	return s.repo.Read(c)
}

func (s *ExtractionResultService) Update(c context.Context, id any, input models.ExtractionResultInput) error {
	return s.repo.Update(c, id, input)
}

func (s *ExtractionResultService) Delete(c context.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *ExtractionResultService) GetByID(c context.Context, id any) (entities.ExtractionResult, error) {
	return s.repo.GetByID(c, id)
}

func (s *ExtractionResultService) SaveFromLLM(c context.Context, id uint, input models.ContentLLM) (uint, error) {
	var content entities.ExtractionJSON
	if err := json.Unmarshal([]byte(input.Content), &content); err != nil {
		return 0, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	jsonContent, err := jsonb.New(content)
	if err != nil {
		return 0, fmt.Errorf("failed to serialize json: %w", err)
	}

	result := models.ExtractionResultInput{
		Raw:        input.Content,
		Json:       jsonContent,
		WorkflowID: id,
	}
	return s.repo.Create(c, result)
}
