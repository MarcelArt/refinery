package entities

import (
	"fmt"
	"strings"

	"github.com/MarcelArt/refinery/pkg/jsonb"
	"gorm.io/gorm"
)

type Workflow struct {
	gorm.Model
	Title       string                       `gorm:"not null" json:"title"`
	Description string                       `gorm:"not null" json:"description"`
	Prompt      string                       `json:"prompt"`
	Schemas     jsonb.JSONB[workflowSchemas] `json:"schemas"`
	Type        string                       `gorm:"default:PDF_TEXT" json:"type"` // PDF_TEXT, PICTURE

	UserID uint `json:"userId"`

	User *User `json:"user,omitzero"`
}

type WorkflowSchema struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Example     string `json:"example"`
}
type workflowSchemas []WorkflowSchema

func (e workflowSchemas) ToMarkdownTable() string {
	var schemaStr strings.Builder
	schemaStr.WriteString("| Key    | Type          | Description              | Example             |\n")
	schemaStr.WriteString("| ------ | ------------- | ------------------------ | ------------------- |\n")
	for _, schema := range e {
		fmt.Fprintf(&schemaStr, "| %s | %s | %s | %s |\n", schema.Key, schema.Type, schema.Description, schema.Example)
	}

	return schemaStr.String()
}

func (e workflowSchemas) ToJSONExample() string {
	var jsonExample strings.Builder
	jsonExample.WriteString("[{\n")
	count := len(e)
	for i, schema := range e {
		switch schema.Type {
		case "boolean", "number":
			fmt.Fprintf(&jsonExample, `\t"%s": %s`, schema.Key, schema.Example)
		default:
			fmt.Fprintf(&jsonExample, `\t"%s": "%s"`, schema.Key, schema.Example)
		}

		if i < count-1 {
			fmt.Fprintf(&jsonExample, ",\n")
		}
	}
	jsonExample.WriteString("\n}]")

	return jsonExample.String()
}
