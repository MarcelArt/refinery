package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/entities"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/repositories"
	"github.com/gin-gonic/gin"
	"github.com/morkid/paginate"
)

type IWorkflowService interface {
	common.IBaseCrudService[entities.Workflow, models.WorkflowInput, models.WorkflowPage]
	UploadToWorkflow(c context.Context, id any, filename string, file multipart.File) error
}

type WorkflowService struct {
	repo  repositories.IWorkflowRepo
	nRepo repositories.IN8NRepo
}

var _ IWorkflowService = &WorkflowService{}

func NewWorkflowService(repo repositories.IWorkflowRepo, nRepo repositories.IN8NRepo) *WorkflowService {
	return &WorkflowService{
		repo:  repo,
		nRepo: nRepo,
	}
}

func (s *WorkflowService) Create(c context.Context, input models.WorkflowInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *WorkflowService) Read(c *gin.Context) (paginate.Page, []models.WorkflowPage) {
	return s.repo.Read(c)
}

func (s *WorkflowService) Update(c context.Context, id any, input models.WorkflowInput) error {
	return s.repo.Update(c, id, input)
}

func (s *WorkflowService) Delete(c context.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *WorkflowService) GetByID(c context.Context, id any) (entities.Workflow, error) {
	return s.repo.GetByID(c, id)
}

func (s *WorkflowService) UploadToWorkflow(c context.Context, id any, filename string, file multipart.File) error {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("failed creating form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed copying file: %w", err)
	}

	workflow, err := s.GetByID(c, id)
	if err != nil {
		return err
	}

	schemas, err := workflow.Schemas.Deserialize()
	if err != nil {
		return fmt.Errorf("failed deserialize schema: %w", err)
	}
	var schemaStr strings.Builder
	schemaStr.WriteString("| Key    | Type          | Description              | Example             |\n")
	schemaStr.WriteString("| ------ | ------------- | ------------------------ | ------------------- |\n")
	for _, schema := range schemas {
		fmt.Fprintf(&schemaStr, "| %s | %s | %s | %s |\n", schema.Key, schema.Type, schema.Description, schema.Example)
	}
	prompt := fmt.Sprintf("based on markdown text above, extract the information into json format without markdown code block with this specification:\n\n%s", schemaStr.String())

	writer.WriteField("prompt", prompt)
	writer.WriteField("system", "you may only reply in json format without markdown code block")

	contentType := writer.FormDataContentType()

	writer.Close()
	return s.nRepo.PostWebhookForm("48c2f9e5-a3a5-4582-9f47-7792c790d701", &requestBody, contentType)
}
