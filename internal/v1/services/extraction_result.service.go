package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/entities"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/repositories"
	"github.com/MarcelArt/refinery/pkg/fetch"
	"github.com/MarcelArt/refinery/pkg/jsonb"
	"github.com/gin-gonic/gin"
	"github.com/morkid/paginate"
)

type IExtractionResultService interface {
	common.IBaseCrudService[entities.ExtractionResult, models.ExtractionResultInput, models.ExtractionResultPage]
	SaveFromLLM(c context.Context, id uint, input models.ContentLLM) error
	GetByWorkflowID(c *gin.Context, workflowID any) (paginate.Page, []models.ExtractionResultPage)
}

type ExtractionResultService struct {
	repo  repositories.IExtractionResultRepo
	wRepo repositories.IWebhookRepo
	c     *http.Client
}

var _ IExtractionResultService = &ExtractionResultService{}

func NewExtractionResultService(repo repositories.IExtractionResultRepo, wRepo repositories.IWebhookRepo) *ExtractionResultService {
	return &ExtractionResultService{
		repo:  repo,
		wRepo: wRepo,
		c:     &http.Client{},
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

func (s *ExtractionResultService) SaveFromLLM(c context.Context, id uint, input models.ContentLLM) error {
	var content entities.ExtractionJSON
	if err := json.Unmarshal([]byte(input.Content), &content); err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	jsonContent, err := jsonb.New(content)
	if err != nil {
		return fmt.Errorf("failed to serialize json: %w", err)
	}

	result := models.ExtractionResultInput{
		Raw:        input.Content,
		Json:       jsonContent,
		Source:     input.Source,
		Status:     "DONE",
		FinishedAt: new(time.Now()),
	}
	if err := s.repo.Update(c, id, result); err != nil {
		return fmt.Errorf("failed to update extraction result: %w", err)
	}

	return s.triggerWebhooks(c, id, result, input)
}

func (s *ExtractionResultService) GetByWorkflowID(c *gin.Context, workflowID any) (paginate.Page, []models.ExtractionResultPage) {
	return s.repo.GetByWorkflowID(c, workflowID)
}

func (s *ExtractionResultService) triggerWebhooks(c context.Context, id uint, result models.ExtractionResultInput, metadata models.ContentLLM) error {
	extraction, err := s.repo.GetByID(c, id)
	if err != nil {
		return err
	}

	webhooks, err := s.wRepo.GetAllByWorkflowID(c, extraction.WorkflowID)
	if err != nil {
		return fmt.Errorf("failed retrieving webhooks: %w", err)
	}

	extraction.Raw = result.Raw
	extraction.Json = result.Json
	extraction.Source = result.Source
	extraction.Status = result.Status
	extraction.FinishedAt = result.FinishedAt

	data := make(map[string]any)
	data["extractionResult"] = extraction
	data["metadata"] = metadata.Metadata

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed marshalling webhook body: %w", err)
	}

	for _, webhook := range webhooks {
		mac := hmac.New(sha256.New, []byte(webhook.HmacKey))
		mac.Write(body)
		signature := hex.EncodeToString(mac.Sum(nil))
		headers := make(map[string]string)
		headers["X-Refinery-Signature"] = signature
		fetch.Fetch[any](s.c, webhook.Method, webhook.URL, data, headers)
	}

	return nil
}
