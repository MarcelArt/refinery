package repositories

import (
	"fmt"
	"net/http"

	"github.com/MarcelArt/refinery/internal/configs"
	"github.com/MarcelArt/refinery/pkg/fetch"
)

type IN8NRepo interface {
	PostWebhook(path string) error
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

func (r *N8NRepo) PostWebhook(path string) error {
	url := fmt.Sprintf("%s/%s", r.baseURL, path)
	_, err := fetch.Fetch[any](r.client, http.MethodPost, url, nil, nil)

	return err
}
