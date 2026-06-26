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

type IWebhookRepo interface {
	common.IBaseCrudRepo[entities.Webhook, models.WebhookInput, models.WebhookPage]
	GetByWorkflowID(c *gin.Context, workflowID any) (paginate.Page, []models.WebhookPage)
}

type WebhookRepo struct {
	db        *gorm.DB
	pageQuery string
}

var _ IWebhookRepo = &WebhookRepo{}

func NewWebhookRepo(db *gorm.DB) *WebhookRepo {
	return &WebhookRepo{
		db: db,
		pageQuery: `
			SELECT 
				*
			FROM webhooks w
			where w.deleted_at isnull
		`,
	}
}

func (r *WebhookRepo) Create(c context.Context, input models.WebhookInput) (uint, error) {
	webhook, err := common.Cast[entities.Webhook](input)
	if err != nil {
		return 0, fmt.Errorf("cannot cast input: %w", err)
	}

	err = gorm.G[entities.Webhook](r.db).Create(c, &webhook)

	return webhook.ID, err
}

func (r *WebhookRepo) Read(c *gin.Context) (paginate.Page, []models.WebhookPage) {
	webhooks := make([]models.WebhookPage, 0)

	stmt := r.db.Raw(r.pageQuery)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request).Response(&webhooks)

	return page, webhooks
}

func (r *WebhookRepo) GetByWorkflowID(c *gin.Context, workflowID any) (paginate.Page, []models.WebhookPage) {
	webhooks := make([]models.WebhookPage, 0)
	query := `
		SELECT 
			*
		FROM webhooks w
		where w.deleted_at isnull
		and workflow_id = ?
	`

	stmt := r.db.Raw(query, workflowID)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request).Response(&webhooks)

	return page, webhooks
}

func (r *WebhookRepo) Update(c context.Context, id any, input models.WebhookInput) error {
	webhook, err := common.Cast[entities.Webhook](input)
	if err != nil {
		return fmt.Errorf("cannot cast input: %w", err)
	}

	_, err = gorm.G[entities.Webhook](r.db).Where("id = ?", id).Updates(c, webhook)

	return err
}

func (r *WebhookRepo) Delete(c context.Context, id any) error {
	_, err := gorm.G[entities.Webhook](r.db).Where("id = ?", id).Delete(c)

	return err
}

func (r *WebhookRepo) GetByID(c context.Context, id any) (entities.Webhook, error) {
	var webhook entities.Webhook

	webhook, err := gorm.G[entities.Webhook](r.db).Where("id = ?", id).First(c)

	return webhook, err
}
