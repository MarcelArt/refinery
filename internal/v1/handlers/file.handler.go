package handlers

import (
	"github.com/MarcelArt/refinery/internal/v1/services"
)

type FileHandler struct {
	service services.IFileService
}

func NewFileHandler(service services.IFileService) *FileHandler {
	return &FileHandler{
		service: service,
	}
}
