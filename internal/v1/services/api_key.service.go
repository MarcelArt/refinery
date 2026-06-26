package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/entities"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/repositories"
	"github.com/gin-gonic/gin"
	"github.com/morkid/paginate"
)

type IApiKeyService interface {
	common.IBaseCrudService[entities.ApiKey, models.ApiKeyInput, models.ApiKeyPage]
	Generate(c context.Context, input models.ApiKeyInput) (string, error)
	Regenerate(c context.Context, id any) (string, error)
	GetByKey(c context.Context, key string) (entities.ApiKey, error)
	GetByUserID(c *gin.Context, userID any) (paginate.Page, []models.ApiKeyPage)
}

type ApiKeyService struct {
	repo repositories.IApiKeyRepo
}

var _ IApiKeyService = &ApiKeyService{}

func NewApiKeyService(repo repositories.IApiKeyRepo) *ApiKeyService {
	return &ApiKeyService{
		repo: repo,
	}
}

func (s *ApiKeyService) Create(c context.Context, input models.ApiKeyInput) (uint, error) {
	if input.Key == "" {
		input.Key = generateSecureKey()
	}
	return s.repo.Create(c, input)
}

func (s *ApiKeyService) Read(c *gin.Context) (paginate.Page, []models.ApiKeyPage) {
	return s.repo.Read(c)
}

func (s *ApiKeyService) Update(c context.Context, id any, input models.ApiKeyInput) error {
	return s.repo.Update(c, id, input)
}

func (s *ApiKeyService) Delete(c context.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *ApiKeyService) GetByID(c context.Context, id any) (entities.ApiKey, error) {
	return s.repo.GetByID(c, id)
}

func (s *ApiKeyService) Generate(c context.Context, input models.ApiKeyInput) (string, error) {
	input.Key = generateSecureKey()
	if _, err := s.repo.Create(c, input); err != nil {
		return "", fmt.Errorf("failed saving api key: %w", err)
	}

	return input.Key, nil
}

func (s *ApiKeyService) Regenerate(c context.Context, id any) (string, error) {
	oldKey, err := s.repo.GetByID(c, id)
	if err != nil {
		return "", err
	}

	input := models.ApiKeyInput{
		Key:    generateSecureKey(),
		Scopes: oldKey.Scopes,
	}

	if err := s.repo.Update(c, id, input); err != nil {
		return "", fmt.Errorf("failed regenerating api key: %w", err)
	}

	return input.Key, nil
}

func (s *ApiKeyService) GetByKey(c context.Context, key string) (entities.ApiKey, error) {
	return s.repo.GetByKey(c, key)
}

func (s *ApiKeyService) GetByUserID(c *gin.Context, userID any) (paginate.Page, []models.ApiKeyPage) {
	return s.repo.GetByUserID(c, userID)
}

func generateSecureKey() string {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return "rf_pat_" + hex.EncodeToString(b)
}
