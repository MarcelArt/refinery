package services

import (
	"context"
	"fmt"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/entities"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/repositories"
	"github.com/gin-gonic/gin"
	"github.com/morkid/paginate"
)

type IWebhookService interface {
	common.IBaseCrudService[entities.Webhook, models.WebhookInput, models.WebhookPage]
	GetByWorkflowID(c *gin.Context, workflowID any) (paginate.Page, []models.WebhookPage)
	CreateWithHMAC(c context.Context, input models.WebhookInput) (string, error)
}

type WebhookService struct {
	repo repositories.IWebhookRepo
}

var _ IWebhookService = &WebhookService{}

func NewWebhookService(repo repositories.IWebhookRepo) *WebhookService {
	return &WebhookService{
		repo: repo,
	}
}

func (s *WebhookService) Create(c context.Context, input models.WebhookInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *WebhookService) Read(c *gin.Context) (paginate.Page, []models.WebhookPage) {
	return s.repo.Read(c)
}

func (s *WebhookService) Update(c context.Context, id any, input models.WebhookInput) error {
	return s.repo.Update(c, id, input)
}

func (s *WebhookService) Delete(c context.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *WebhookService) GetByID(c context.Context, id any) (entities.Webhook, error) {
	return s.repo.GetByID(c, id)
}

func (s *WebhookService) GetByWorkflowID(c *gin.Context, workflowID any) (paginate.Page, []models.WebhookPage) {
	return s.repo.GetByWorkflowID(c, workflowID)
}

func (s *WebhookService) CreateWithHMAC(c context.Context, input models.WebhookInput) (string, error) {
	input.HmacKey = common.GenerateSecureKey("whsec_")
	if _, err := s.repo.Create(c, input); err != nil {
		return "", fmt.Errorf("failed saving webhook setting: %w", err)
	}

	return input.HmacKey, nil
}
