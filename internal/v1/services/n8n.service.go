package services

import (
	"github.com/MarcelArt/refinery/internal/v1/repositories"
)

type IN8NService interface {
	PostWebhook(path string) error
}

type N8NService struct {
	nRepo repositories.IN8NRepo
}

func NewN8NService(nRepo repositories.IN8NRepo) *N8NService {
	return &N8NService{
		nRepo: nRepo,
	}
}

func (s *N8NService) PostWebhook(path string) error {
	return nil
	// return s.nRepo.PostWebhook(path)
}
