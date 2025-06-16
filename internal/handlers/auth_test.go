package handlers

import (
	"net/http"
	"net/url"
	"testing"

	testhelper "github.com/jeffrpowell/listaway/internal/testing"
)

func TestLoginHandler(t *testing.T) {
	// Setup test environment
	db, recorder := testhelper.SetupTestHandler(t)
	defer db.TeardownTestDB(t)
	
	// Clean up database tables
	db.CleanupTables(t)
	
	// Create a test user in the database
	// This will depend on how users are created in the system
	// You may need to adapt this based on your actual implementation
	
	// Test case 1: GET request - should return login form
	req, err := testhelper.CreateTestRequest("GET", "/login", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	// Call the login handler (this needs to be updated with the actual handler function)
	// LoginHandler(recorder, req)
	
	// Assert response
	testhelper.AssertStatusCode(t, recorder.Code, http.StatusOK)
	testhelper.AssertContentType(t, recorder.Header(), "text/html")
	testhelper.AssertBodyContains(t, recorder.Body.String(), "<form")
	
	// Test case 2: POST request with valid credentials
	formData := url.Values{}
	formData.Set("email", "test@example.com")
	formData.Set("password", "password123")
	
	req, err = testhelper.CreateTestRequest("POST", "/login", formData)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	// Reset recorder
	recorder = &testhelper.MockResponseRecorder{
		ResponseRecorder: recorder.ResponseRecorder,
	}
	
	// Call the login handler
	// LoginHandler(recorder, req)
	
	// Assert successful login redirect
	testhelper.AssertStatusCode(t, recorder.Code, http.StatusSeeOther)
	
	// Test case 3: POST request with invalid credentials
	formData = url.Values{}
	formData.Set("email", "test@example.com")
	formData.Set("password", "wrongpassword")
	
	req, err = testhelper.CreateTestRequest("POST", "/login", formData)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	// Reset recorder
	recorder = &testhelper.MockResponseRecorder{
		ResponseRecorder: recorder.ResponseRecorder,
	}
	
	// Call the login handler
	// LoginHandler(recorder, req)
	
	// Assert failed login
	testhelper.AssertStatusCode(t, recorder.Code, http.StatusOK) // Stay on login page
	testhelper.AssertBodyContains(t, recorder.Body.String(), "Invalid") // Should contain error message
}

// Additional handler tests would follow a similar pattern
