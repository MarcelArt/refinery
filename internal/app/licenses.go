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

		// Append HTMX, Alpine.js, and Chart.js licenses which are not golang packages
		appendManualLicense := func(pkgName, licenseContent string) {
			builder.WriteString(strings.Repeat("=", 80))
			builder.WriteString("\n")
			fmt.Fprintf(&builder, "Package:      %s\n", pkgName)
			builder.WriteString("License File: LICENSE\n")
			builder.WriteString(strings.Repeat("=", 80))
			builder.WriteString("\n\n")
			builder.WriteString(strings.TrimSpace(licenseContent))
			builder.WriteString("\n\n")
		}

		appendManualLicense("github.com/bigskysoftware/htmx", htmxLicense)
		appendManualLicense("github.com/alpinejs/alpine", alpineJSLicense)
		appendManualLicense("github.com/chartjs/Chart.js", chartJSLicense)

		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(http.StatusOK, builder.String())
	}
}

const htmxLicense = `
Zero-Clause BSD
=============

Permission to use, copy, modify, and/or distribute this software for
any purpose with or without fee is hereby granted.

THE SOFTWARE IS PROVIDED “AS IS” AND THE AUTHOR DISCLAIMS ALL
WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES
OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE
FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY
DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN
AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT
OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
`

const alpineJSLicense = `
# MIT License

Copyright © 2019-2025 Caleb Porzio and contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
`

const chartJSLicense = `
The MIT License (MIT)

Copyright (c) 2014-2024 Chart.js Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
`
