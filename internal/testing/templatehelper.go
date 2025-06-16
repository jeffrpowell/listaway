package testing

import (
	"bytes"
	"html/template"
	"io/fs"
	"path/filepath"
	"testing"
)

// TemplateRenderer helps test HTML templates
type TemplateRenderer struct {
	Templates *template.Template
}

// NewTemplateRenderer creates a new template renderer with templates from the specified directory
func NewTemplateRenderer(t *testing.T, templateDir string) *TemplateRenderer {
	templates := template.New("")
	
	err := filepath.Walk(templateDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			_, err = templates.ParseFiles(path)
			if err != nil {
				t.Logf("Error parsing template %s: %v", path, err)
				return err
			}
		}
		
		return nil
	})
	
	if err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}
	
	return &TemplateRenderer{Templates: templates}
}

// RenderTemplate renders a template with the given data and returns the result
func (tr *TemplateRenderer) RenderTemplate(t *testing.T, templateName string, data interface{}) string {
	var buf bytes.Buffer
	
	err := tr.Templates.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		t.Fatalf("Failed to render template %s: %v", templateName, err)
	}
	
	return buf.String()
}

// AssertTemplateContains checks if rendered template contains expected content
func (tr *TemplateRenderer) AssertTemplateContains(t *testing.T, templateName string, data interface{}, expectedContent string) {
	rendered := tr.RenderTemplate(t, templateName, data)
	AssertBodyContains(t, rendered, expectedContent)
}

// MockTemplateData returns common data structures used in templates
// This will need to be expanded based on the actual data structures your templates use
func MockTemplateData() map[string]interface{} {
	return map[string]interface{}{
		"User": map[string]interface{}{
			"ID":            int64(1),
			"Email":         "test@example.com",
			"Name":          "Test User",
			"Admin":         true,
			"InstanceAdmin": false,
		},
		"Lists": []map[string]interface{}{
			{
				"ID":          int64(1),
				"Name":        "Test List",
				"Description": "A test list",
				"ShareCode":   "abc123",
			},
		},
		"Items": []map[string]interface{}{
			{
				"ID":       int64(1),
				"ListID":   int64(1),
				"Name":     "Test Item",
				"URL":      "https://example.com",
				"Notes":    "Test notes",
				"Priority": 1,
			},
		},
	}
}
