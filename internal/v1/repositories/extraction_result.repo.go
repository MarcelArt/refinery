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
	GetDailyThroughput(c context.Context, userID any) ([]models.ThroughputPoint, error)
	GetLatencyStats(c context.Context, userID any) (models.LatencyStats, error)
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

func (r *ExtractionResultRepo) GetDailyThroughput(c context.Context, userID any) ([]models.ThroughputPoint, error) {
	query := `
		WITH user_results AS (
			SELECT er.created_at, er.status
			FROM extraction_results er
			JOIN workflows w ON er.workflow_id = w.id
			WHERE er.deleted_at IS NULL
			AND w.user_id = ?
		)
		SELECT
			d.day AS bucket,
			COUNT(ur.created_at) FILTER (WHERE ur.status = 'DONE')        AS done,
			COUNT(ur.created_at) FILTER (WHERE ur.status = 'FAILED')      AS failed,
			COUNT(ur.created_at) FILTER (WHERE ur.status = 'IN_PROGRESS') AS in_progress
		FROM generate_series(
			date_trunc('day', now() - interval '13 days'),
			date_trunc('day', now()),
			interval '1 day'
		) AS d(day)
		LEFT JOIN user_results ur
			ON date_trunc('day', ur.created_at) = d.day
		GROUP BY d.day
		ORDER BY d.day
	`

	var throughputPoints []models.ThroughputPoint
	err := gorm.G[entities.ExtractionResult](r.db).Raw(query, userID).Scan(c, &throughputPoints)
	return throughputPoints, err
}

func (r *ExtractionResultRepo) GetLatencyStats(c context.Context, userID any) (models.LatencyStats, error) {
	query := `
		WITH latencies AS (
			SELECT EXTRACT(EPOCH FROM (er.finished_at - er.created_at)) AS seconds
			FROM extraction_results er
			JOIN workflows w ON er.workflow_id = w.id
			WHERE er.deleted_at IS NULL
			AND w.user_id = ?
			AND er.finished_at IS NOT NULL
			AND er.created_at >= date_trunc('day', now() - interval '29 days')
		)
		SELECT
			COUNT(*)                                                       AS completed,
			AVG(seconds)                                                   AS avg_seconds,
			percentile_cont(0.5)  WITHIN GROUP (ORDER BY seconds)          AS p50_seconds,
			percentile_cont(0.95) WITHIN GROUP (ORDER BY seconds)          AS p95_seconds
		FROM latencies
	`

	var latencyStats models.LatencyStats
	err := gorm.G[entities.ExtractionResult](r.db).Raw(query, userID).Scan(c, &latencyStats)
	return latencyStats, err
}
