package app

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// handleLicenses returns a gin.HandlerFunc that reads and aggregates licenses
func handleLicenses(dir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var builder strings.Builder

		// Check if directory exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			c.String(http.StatusNotFound, "Third-party licenses directory not found. Please run 'make license' to generate it.")
			return
		}

		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			// Strictly only include files containing "license" or "licence" in their name
			base := strings.ToLower(filepath.Base(path))
			if !strings.Contains(base, "license") && !strings.Contains(base, "licence") {
				return nil
			}

			// Read the license file
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Find the relative path of the file to determine the package name
			rel, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}

			// The package name is the directory structure relative to the licenses directory
			pkgName := filepath.Dir(rel)

			builder.WriteString(strings.Repeat("=", 80))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("Package:      %s\n", pkgName))
			builder.WriteString(fmt.Sprintf("License File: %s\n", filepath.Base(path)))
			builder.WriteString(strings.Repeat("=", 80))
			builder.WriteString("\n\n")
			builder.Write(content)
			builder.WriteString("\n\n")

			return nil
		})

		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read licenses: %v", err)
			return
		}

		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(http.StatusOK, builder.String())
	}
}
