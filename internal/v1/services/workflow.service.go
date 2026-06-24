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

type IWorkflowService interface {
	common.IBaseCrudService[entities.Workflow, models.WorkflowInput, models.WorkflowPage]
}

type WorkflowService struct {
	repo repositories.IWorkflowRepo
}

var _ IWorkflowService = &WorkflowService{}

func NewWorkflowService(repo repositories.IWorkflowRepo) *WorkflowService {
	return &WorkflowService{
		repo: repo,
	}
}

func (s *WorkflowService) Create(c context.Context, input models.WorkflowInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *WorkflowService) Read(c *gin.Context) (paginate.Page, []models.WorkflowPage) {
	return s.repo.Read(c)
}

func (s *WorkflowService) Update(c context.Context, id any, input models.WorkflowInput) error {
	return s.repo.Update(c, id, input)
}

func (s *WorkflowService) Delete(c context.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *WorkflowService) GetByID(c context.Context, id any) (entities.Workflow, error) {
	return s.repo.GetByID(c, id)
}
