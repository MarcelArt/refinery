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

type IExtractionResultRepo interface {
	common.IBaseCrudRepo[entities.ExtractionResult, models.ExtractionResultInput, models.ExtractionResultPage]
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
				*
			FROM extraction_results er
			where er.deleted_at isnull
		`,
	}
}

func (r *ExtractionResultRepo) Create(c context.Context, input models.ExtractionResultInput) (uint, error) {
	extractionResult, err := common.Cast[entities.ExtractionResult](input)
	if err != nil {
		return 0, fmt.Errorf("cannot cast input: %w", err)
	}

	err = gorm.G[entities.ExtractionResult](r.db).Create(c, &extractionResult)

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
	extractionResult, err := common.Cast[entities.ExtractionResult](input)
	if err != nil {
		return fmt.Errorf("cannot cast input: %w", err)
	}

	_, err = gorm.G[entities.ExtractionResult](r.db).Where("id = ?", id).Updates(c, extractionResult)

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
