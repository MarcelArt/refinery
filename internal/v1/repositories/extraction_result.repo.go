package repositories

import (
	"context"
	"fmt"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/entities"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/devfeel/mapper"
	"github.com/gin-gonic/gin"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type IExtractionResultRepo interface {
	common.IBaseCrudRepo[entities.ExtractionResult, models.ExtractionResultInput, models.ExtractionResultPage]
	GetByWorkflowID(c *gin.Context, workflowID any) (paginate.Page, []models.ExtractionResultPage)
	GetStatusCount(c context.Context, status string, userID any) (float32, error)
}

type ExtractionResultRepo struct {
	db        *gorm.DB
	pageQuery string
}

var _ IExtractionResultRepo = &ExtractionResultRepo{}

func NewExtractionResultRepo(db *gorm.DB) *ExtractionResultRepo {
	return &ExtractionResultRepo{
		db: db,
		pageQuery: `
			SELECT 
				er.*,
				w.type workflow_type
			FROM extraction_results er
			join workflows w on er.workflow_id = w.id 
			where er.deleted_at isnull
		`,
	}
}

func (r *ExtractionResultRepo) Create(c context.Context, input models.ExtractionResultInput) (uint, error) {
	var extractionResult entities.ExtractionResult
	if err := mapper.AutoMapper(&input, &extractionResult); err != nil {
		return 0, fmt.Errorf("cannot map input: %w", err)
	}

	err := gorm.G[entities.ExtractionResult](r.db).Create(c, &extractionResult)

	return extractionResult.ID, err
}

func (r *ExtractionResultRepo) Read(c *gin.Context) (paginate.Page, []models.ExtractionResultPage) {
	extractionResults := make([]models.ExtractionResultPage, 0)

	stmt := r.db.Raw(r.pageQuery)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request).Response(&extractionResults)

	return page, extractionResults
}

func (r *ExtractionResultRepo) Update(c context.Context, id any, input models.ExtractionResultInput) error {
	var extractionResult entities.ExtractionResult
	if err := mapper.AutoMapper(&input, &extractionResult); err != nil {
		return fmt.Errorf("cannot map input: %w", err)
	}

	_, err := gorm.G[entities.ExtractionResult](r.db).Where("id = ?", id).Updates(c, extractionResult)

	return err
}

func (r *ExtractionResultRepo) Delete(c context.Context, id any) error {
	_, err := gorm.G[entities.ExtractionResult](r.db).Where("id = ?", id).Delete(c)

	return err
}

func (r *ExtractionResultRepo) GetByID(c context.Context, id any) (entities.ExtractionResult, error) {
	var extractionResult entities.ExtractionResult

	extractionResult, err := gorm.G[entities.ExtractionResult](r.db).Where("id = ?", id).First(c)

	return extractionResult, err
}

func (r *ExtractionResultRepo) GetByWorkflowID(c *gin.Context, workflowID any) (paginate.Page, []models.ExtractionResultPage) {
	extractionResults := make([]models.ExtractionResultPage, 0)
	query := `
		SELECT 
			er.*,
			w.type workflow_type
		FROM extraction_results er
		join workflows w on er.workflow_id = w.id 
		where er.deleted_at isnull
		and er.workflow_id = ?
		order by er.id desc
	`

	stmt := r.db.Raw(query, workflowID)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request).Response(&extractionResults)

	return page, extractionResults
}

func (r *ExtractionResultRepo) GetStatusCount(c context.Context, status string, userID any) (float32, error) {
	var count float32
	query := `
		select
			count(1)
		from extraction_results er
		join workflows w on er.workflow_id = w.id
		where er.deleted_at isnull
		and w.user_id = ?
		and er.status = ?
	`

	err := gorm.G[entities.ExtractionResult](r.db).Raw(query, userID, status).Scan(c, &count)
	return count, err
}
