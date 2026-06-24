package handlers

import (
	"github.com/MarcelArt/refinery/internal/v1/services"
)

type FileHandler struct {
	service services.IN8NService
}

func NewFileHandler(service services.IN8NService) *FileHandler {
	return &FileHandler{
		service: service,
	}
}
