package web

import (
	"testing"

	testhelper "github.com/jeffrpowell/listaway/internal/testing"
)

func TestTemplateRendering(t *testing.T) {
	// Create a template renderer pointing to your templates directory
	// You may need to adjust the path based on your actual template location
	renderer := testhelper.NewTemplateRenderer(t, "./templates")
	
	// Test case: List template
	testData := testhelper.MockTemplateData()
	
	// Example test for the list page template
	renderer.AssertTemplateContains(t, "list.html", testData, "Test List")
	renderer.AssertTemplateContains(t, "list.html", testData, "Test Item")
	
	// Test case: Login template
	loginData := map[string]interface{}{
		"Error": "",
	}
	renderer.AssertTemplateContains(t, "login.html", loginData, "Login")
	
	// Test with error message
	loginWithError := map[string]interface{}{
		"Error": "Invalid credentials",
	}
	renderer.AssertTemplateContains(t, "login.html", loginWithError, "Invalid credentials")
}
