package repositories

import (
	"fmt"
	"io"
	"log"
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
	log.Println("url :>> ", url)
	// log.Println("body :>> ", body)
	log.Println("contentType :>> ", contentType)
	_, err := fetch.FormFile[any](r.client, http.MethodPost, url, body, contentType, nil)
	// _, err := fetch.FormFile[any](r.client, http.MethodPost, "https://webhook.site/b1052ffc-daf9-427c-945f-6cfa6eb5b105", body, contentType, nil)

	return err
}
