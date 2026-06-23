package services

import (
	"github.com/MarcelArt/refinery/internal/v1/repositories"
)

type IFileService interface {
	PostWebhook(path string) error
}

type FileService struct {
	nRepo repositories.IN8NRepo
}

func NewFileService(nRepo IFileService) *FileService {
	return &FileService{
		nRepo: nRepo,
	}
}

func (s *FileService) PostWebhook(path string) error {
	return s.nRepo.PostWebhook(path)
}
