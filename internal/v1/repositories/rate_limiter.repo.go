package repositories

import (
	"context"
	"fmt"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/entities"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/pkg/date"
	"github.com/devfeel/mapper"
	"github.com/gin-gonic/gin"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type IRateLimiterRepo interface {
	common.IBaseCrudRepo[entities.RateLimiter, models.RateLimiterInput, models.RateLimiterPage]
	GetTodayByUserID(c context.Context, userID any) (entities.RateLimiter, error)
}

type RateLimiterRepo struct {
	db        *gorm.DB
	pageQuery string
}

var _ IRateLimiterRepo = &RateLimiterRepo{}

func NewRateLimiterRepo(db *gorm.DB) *RateLimiterRepo {
	return &RateLimiterRepo{
		db: db,
		pageQuery: `
			SELECT 
				*
			FROM rate_limiters rl
			where rl.deleted_at isnull
		`,
	}
}

func (r *RateLimiterRepo) Create(c context.Context, input models.RateLimiterInput) (uint, error) {
	var rateLimiter entities.RateLimiter
	if err := mapper.AutoMapper(&input, &rateLimiter); err != nil {
		return 0, fmt.Errorf("cannot map input: %w", err)
	}

	err := gorm.G[entities.RateLimiter](r.db).Create(c, &rateLimiter)

	return rateLimiter.ID, err
}

func (r *RateLimiterRepo) Read(c *gin.Context) (paginate.Page, []models.RateLimiterPage) {
	rateLimiters := make([]models.RateLimiterPage, 0)

	stmt := r.db.Raw(r.pageQuery)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request).Response(&rateLimiters)

	return page, rateLimiters
}

func (r *RateLimiterRepo) Update(c context.Context, id any, input models.RateLimiterInput) error {
	var rateLimiter entities.RateLimiter
	if err := mapper.AutoMapper(&input, &rateLimiter); err != nil {
		return fmt.Errorf("cannot map input: %w", err)
	}

	_, err := gorm.G[entities.RateLimiter](r.db).Where("id = ?", id).Updates(c, rateLimiter)

	return err
}

func (r *RateLimiterRepo) Delete(c context.Context, id any) error {
	_, err := gorm.G[entities.RateLimiter](r.db).Where("id = ?", id).Delete(c)

	return err
}

func (r *RateLimiterRepo) GetByID(c context.Context, id any) (entities.RateLimiter, error) {
	var rateLimiter entities.RateLimiter

	rateLimiter, err := gorm.G[entities.RateLimiter](r.db).Where("id = ?", id).First(c)

	return rateLimiter, err
}

func (r *RateLimiterRepo) GetTodayByUserID(c context.Context, userID any) (entities.RateLimiter, error) {
	now := date.Local()
	today := now.Formats("YYYY-MM-DD")

	query := `
		select
			rl.*
		from rate_limiters rl 
		where rl.deleted_at isnull
		and rl.created_at::date = ?
		and rl.user_id = ?
		order by count desc
	`

	return gorm.G[entities.RateLimiter](r.db).Raw(query, today, userID).First(c)
}
