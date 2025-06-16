package e2e

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestServer holds data about the test server
type TestServer struct {
	URL      string
	Username string
	Password string
}

// NewTestServer creates a TestServer instance
func NewTestServer(t *testing.T) *TestServer {
	// Get server URL from environment or use default
	serverURL := getEnvOrDefault("TEST_SERVER_URL", "http://localhost:8080")
	username := getEnvOrDefault("TEST_ADMIN_USER", "test@example.com")
	password := getEnvOrDefault("TEST_ADMIN_PASSWORD", "testpassword")

	return &TestServer{
		URL:      serverURL,
		Username: username,
		Password: password,
	}
}

// RunE2ETest runs an end-to-end test with a headless browser
func (ts *TestServer) RunE2ETest(t *testing.T, testFunc func(context.Context, *testing.T)) {
	// Create a Chrome instance
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.WindowSize(1280, 1024),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// Create a timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Run the test
	if err := chromedp.Run(ctx); err != nil {
		t.Fatalf("Failed to initialize browser: %v", err)
	}

	testFunc(ctx, t)
}

// LoginToApp logs into the application
func (ts *TestServer) LoginToApp(ctx context.Context, t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.Navigate(ts.URL+"/login"),
		chromedp.WaitVisible(`form[action="/login"]`),
		chromedp.SendKeys(`input[name="email"]`, ts.Username),
		chromedp.SendKeys(`input[name="password"]`, ts.Password),
		chromedp.Submit(`form[action="/login"]`),
		chromedp.WaitVisible(`a[href="/logout"]`), // Wait for logout link to appear, indicating successful login
	)

	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}
}

// AssertElementExists verifies an element exists on the page
func AssertElementExists(ctx context.Context, t *testing.T, selector string) {
	var exists bool
	err := chromedp.Run(ctx,
		chromedp.WaitVisible(selector),
		chromedp.Evaluate(fmt.Sprintf(`document.querySelector('%s') !== null`, selector), &exists),
	)

	if err != nil {
		t.Fatalf("Error checking for element %s: %v", selector, err)
	}

	if !exists {
		t.Fatalf("Element %s not found", selector)
	}
}

// AssertTextContains verifies text content contains expected string
func AssertTextContains(ctx context.Context, t *testing.T, selector, expectedText string) {
	var text string
	err := chromedp.Run(ctx,
		chromedp.WaitVisible(selector),
		chromedp.Text(selector, &text),
	)

	if err != nil {
		t.Fatalf("Error getting text from %s: %v", selector, err)
	}

	if text == "" || !contains(text, expectedText) {
		t.Fatalf("Expected text to contain %q, but got %q", expectedText, text)
	}
}

// contains is a helper function to check if a string contains another
func contains(s, substr string) bool {
	return s != "" && substr != "" && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || s != "" && substr != "" && len(s) > len(substr) && contains(s[1:], substr)))
}

// helper function to get environment variable with default fallback
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
