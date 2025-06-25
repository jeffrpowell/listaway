package testing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/sessions"
)

// MockResponseRecorder extends httptest.ResponseRecorder to include session handling
type MockResponseRecorder struct {
	*httptest.ResponseRecorder
	session *sessions.Session
}

// CreateTestRequest creates a new HTTP request for testing
func CreateTestRequest(method, path string, body interface{}) (*http.Request, error) {
	var reqBody *bytes.Buffer

	if body != nil {
		switch v := body.(type) {
		case string:
			reqBody = bytes.NewBufferString(v)
		case url.Values:
			reqBody = bytes.NewBufferString(v.Encode())
		default:
			jsonData, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			reqBody = bytes.NewBuffer(jsonData)
		}
	} else {
		reqBody = &bytes.Buffer{}
	}

	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		return nil, err
	}

	// Set content type based on body type
	if body != nil {
		switch body.(type) {
		case url.Values:
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		case string:
			if strings.HasPrefix(body.(string), "{") || strings.HasPrefix(body.(string), "[") {
				req.Header.Set("Content-Type", "application/json")
			} else {
				req.Header.Set("Content-Type", "text/plain")
			}
		default:
			req.Header.Set("Content-Type", "application/json")
		}
	}

	return req, nil
}

// SetupTestHandler creates a test handler environment
func SetupTestHandler(t *testing.T) (*TestDB, *MockResponseRecorder) {
	db := SetupTestDB(t)
	recorder := &MockResponseRecorder{
		ResponseRecorder: httptest.NewRecorder(),
	}
	return db, recorder
}

// SetTestSessionValues sets session values for authenticated requests
func SetTestSessionValues(req *http.Request, userID int64, email string, isAdmin bool, isInstanceAdmin bool) {
	// Note: This is a placeholder and will need to be updated
	// based on how Listaway's session handling works
}

// AssertStatusCode asserts that the response has the expected status code
func AssertStatusCode(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("Status code: got %d, want %d", got, want)
	}
}

// AssertContentType asserts that the response has the expected content type
func AssertContentType(t *testing.T, header http.Header, want string) {
	got := header.Get("Content-Type")
	if !strings.HasPrefix(got, want) {
		t.Errorf("Content-Type: got %q, want %q", got, want)
	}
}

// AssertBodyContains asserts that the response body contains the expected string
func AssertBodyContains(t *testing.T, body string, want string) {
	if !strings.Contains(body, want) {
		t.Errorf("Body does not contain %q, body: %q", want, body)
	}
}

// AssertBodyJSON asserts that the response body can be parsed as JSON
// and matches the expected structure
func AssertBodyJSON(t *testing.T, body string, target interface{}) {
	err := json.Unmarshal([]byte(body), target)
	if err != nil {
		t.Errorf("Failed to parse JSON from body: %v, body: %q", err, body)
	}
}
