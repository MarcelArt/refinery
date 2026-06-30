package services

import (
	"context"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/entities"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/repositories"
	"github.com/gin-gonic/gin"
	"github.com/morkid/paginate"
)

type IRateLimiterService interface {
	common.IBaseCrudService[entities.RateLimiter, models.RateLimiterInput, models.RateLimiterPage]
	GetTodayByUserID(c context.Context, userID any) (entities.RateLimiter, error)
}

type RateLimiterService struct {
	repo repositories.IRateLimiterRepo
}

var _ IRateLimiterService = &RateLimiterService{}

func NewRateLimiterService(repo repositories.IRateLimiterRepo) *RateLimiterService {
	return &RateLimiterService{
		repo: repo,
	}
}

func (s *RateLimiterService) Create(c context.Context, input models.RateLimiterInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *RateLimiterService) Read(c *gin.Context) (paginate.Page, []models.RateLimiterPage) {
	return s.repo.Read(c)
}

func (s *RateLimiterService) Update(c context.Context, id any, input models.RateLimiterInput) error {
	return s.repo.Update(c, id, input)
}

func (s *RateLimiterService) Delete(c context.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *RateLimiterService) GetByID(c context.Context, id any) (entities.RateLimiter, error) {
	return s.repo.GetByID(c, id)
}

func (s *RateLimiterService) GetTodayByUserID(c context.Context, userID any) (entities.RateLimiter, error) {
	return s.repo.GetTodayByUserID(c, userID)
}
