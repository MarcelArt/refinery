package repositories

import (
	"fmt"
	"io"
	"net/http"

	"github.com/MarcelArt/refinery/internal/configs"
	"github.com/MarcelArt/refinery/pkg/fetch"
)

type IN8NRepo interface {
	PostWebhookForm(path string, body io.Reader, contentType string) error
}

type N8NRepo struct {
	client  *http.Client
	baseURL string
}

func NewN8NRepo() *N8NRepo {
	return &N8NRepo{
		client:  &http.Client{},
		baseURL: configs.Env.N8NBaseURL,
	}
}

func (r *N8NRepo) PostWebhookForm(path string, body io.Reader, contentType string) error {
	url := fmt.Sprintf("%s/webhook-test/%s", r.baseURL, path)
	headers := make(map[string]string)
	headers["X-Webhook-Key"] = configs.Env.JwtSecret
	_, err := fetch.FormFile[any](r.client, http.MethodPost, url, body, contentType, headers)

	return err
}
