package repositories

import (
	"context"
	"fmt"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/entities"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/gin-gonic/gin"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type IWorkflowRepo interface {
	common.IBaseCrudRepo[entities.Workflow, models.WorkflowInput, models.WorkflowPage]
	GetByUserID(c *gin.Context, userID any) (paginate.Page, []models.WorkflowPage)
}

type WorkflowRepo struct {
	db        *gorm.DB
	pageQuery string
}

var _ IWorkflowRepo = &WorkflowRepo{}

func NewWorkflowRepo(db *gorm.DB) *WorkflowRepo {
	return &WorkflowRepo{
		db: db,
		pageQuery: `
			SELECT 
				*
			FROM workflows w
			where w.deleted_at isnull
		`,
	}
}

func (r *WorkflowRepo) Create(c context.Context, input models.WorkflowInput) (uint, error) {
	workflow, err := common.Cast[entities.Workflow](input)
	if err != nil {
		return 0, fmt.Errorf("cannot cast input: %w", err)
	}

	err = gorm.G[entities.Workflow](r.db).Create(c, &workflow)

	return workflow.ID, err
}

func (r *WorkflowRepo) Read(c *gin.Context) (paginate.Page, []models.WorkflowPage) {
	workflows := make([]models.WorkflowPage, 0)

	stmt := r.db.Raw(r.pageQuery)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request).Response(&workflows)

	return page, workflows
}

func (r *WorkflowRepo) Update(c context.Context, id any, input models.WorkflowInput) error {
	workflow, err := common.Cast[entities.Workflow](input)
	if err != nil {
		return fmt.Errorf("cannot cast input: %w", err)
	}

	_, err = gorm.G[entities.Workflow](r.db).Where("id = ?", id).Updates(c, workflow)

	return err
}

func (r *WorkflowRepo) Delete(c context.Context, id any) error {
	_, err := gorm.G[entities.Workflow](r.db).Where("id = ?", id).Delete(c)

	return err
}

func (r *WorkflowRepo) GetByID(c context.Context, id any) (entities.Workflow, error) {
	var workflow entities.Workflow

	workflow, err := gorm.G[entities.Workflow](r.db).Where("id = ?", id).First(c)

	return workflow, err
}

func (r *WorkflowRepo) GetByUserID(c *gin.Context, userID any) (paginate.Page, []models.WorkflowPage) {
	workflows := make([]models.WorkflowPage, 0)
	query := `
		SELECT 
			*
		FROM workflows w
		where w.deleted_at isnull
		and w.user_id = ?
	`

	stmt := r.db.Raw(query, userID)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request).Response(&workflows)

	return page, workflows
}
