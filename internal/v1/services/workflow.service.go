package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"

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
	GetByUserID(c *gin.Context, userID any) (paginate.Page, []models.WorkflowPage)
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
	schemaStr := schemas.ToMarkdownTable()

	// Text below "Source Text" heading are empty because n8n workflow will appends it
	prompt := fmt.Sprintf(`
		%s

		### Extraction Specification
		%s

		### Source Text
	`, workflow.Prompt, schemaStr)

	writer.WriteField("prompt", prompt)
	writer.WriteField("system", `
		You are a precise data extraction engine. Analyze the source text and extract all matching entities based on the user's specification table. 

		CRITICAL: Your entire response must be a single JSON Array containing one or more JSON Objects. Even if you only extract a single row or item, it MUST be wrapped inside a JSON Array. Do not use an outer wrapper object.

		Example output structure:
		[
			{ "Key1": "Value1", "Key2": "Value2" }
		]
	`)
	writer.WriteField("workflowId", strconv.Itoa(int(workflow.ID)))

	contentType := writer.FormDataContentType()

	writer.Close()
	return s.nRepo.PostWebhookForm("48c2f9e5-a3a5-4582-9f47-7792c790d701", &requestBody, contentType)
}

func (s *WorkflowService) GetByUserID(c *gin.Context, userID any) (paginate.Page, []models.WorkflowPage) {
	return s.repo.GetByUserID(c, userID)
}
