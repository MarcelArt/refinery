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

type IApiKeyRepo interface {
	common.IBaseCrudRepo[entities.ApiKey, models.ApiKeyInput, models.ApiKeyPage]
	GetByKey(c context.Context, key string) (entities.ApiKey, error)
	GetByUserID(c *gin.Context, userID any) (paginate.Page, []models.ApiKeyPage)
}

type ApiKeyRepo struct {
	db        *gorm.DB
	pageQuery string
}

var _ IApiKeyRepo = &ApiKeyRepo{}

func NewApiKeyRepo(db *gorm.DB) *ApiKeyRepo {
	return &ApiKeyRepo{
		db: db,
		pageQuery: `
			SELECT 
				*
			FROM api_keys ak
			where ak.deleted_at isnull
		`,
	}
}

func (r *ApiKeyRepo) Create(c context.Context, input models.ApiKeyInput) (uint, error) {
	apiKey, err := common.Cast[entities.ApiKey](input)
	if err != nil {
		return 0, fmt.Errorf("cannot cast input: %w", err)
	}

	err = gorm.G[entities.ApiKey](r.db).Create(c, &apiKey)

	return apiKey.ID, err
}

func (r *ApiKeyRepo) Read(c *gin.Context) (paginate.Page, []models.ApiKeyPage) {
	apiKeys := make([]models.ApiKeyPage, 0)

	stmt := r.db.Raw(r.pageQuery)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request).Response(&apiKeys)

	return page, apiKeys
}

func (r *ApiKeyRepo) Update(c context.Context, id any, input models.ApiKeyInput) error {
	apiKey, err := common.Cast[entities.ApiKey](input)
	if err != nil {
		return fmt.Errorf("cannot cast input: %w", err)
	}

	_, err = gorm.G[entities.ApiKey](r.db).Where("id = ?", id).Updates(c, apiKey)

	return err
}

func (r *ApiKeyRepo) Delete(c context.Context, id any) error {
	_, err := gorm.G[entities.ApiKey](r.db).Where("id = ?", id).Delete(c)

	return err
}

func (r *ApiKeyRepo) GetByID(c context.Context, id any) (entities.ApiKey, error) {
	var apiKey entities.ApiKey

	apiKey, err := gorm.G[entities.ApiKey](r.db).Where("id = ?", id).First(c)

	return apiKey, err
}

func (r *ApiKeyRepo) GetByKey(c context.Context, key string) (entities.ApiKey, error) {
	var apiKey entities.ApiKey

	apiKey, err := gorm.G[entities.ApiKey](r.db).Where("key = ?", key).First(c)

	return apiKey, err
}

func (r *ApiKeyRepo) GetByUserID(c *gin.Context, userID any) (paginate.Page, []models.ApiKeyPage) {
	apiKeys := make([]models.ApiKeyPage, 0)
	query := `
		SELECT 
			*
		FROM api_keys ak
		where ak.deleted_at isnull
		and ak.user_id = ?
	`

	stmt := r.db.Raw(query, userID)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request).Response(&apiKeys)

	return page, apiKeys
}
