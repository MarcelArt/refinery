package common

import (
	"fmt"
	"strings"
	"text/template"
)

func TextTemplating[T any](text string, data T) (string, error) {
	tmpl, err := template.New("text-template").Parse(text)
	if err != nil {
		return "", fmt.Errorf("failed creating template: %w", err)
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("failed to put data into template: %w", err)
	}

	return result.String(), nil
}
