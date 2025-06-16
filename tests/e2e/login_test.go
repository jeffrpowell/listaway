package e2e

import (
	"context"
	"testing"

	"github.com/chromedp/chromedp"
)

func TestLoginFlow(t *testing.T) {
	// Skip this test in normal test runs unless specifically requested
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}
	
	testServer := NewTestServer(t)
	testServer.RunE2ETest(t, func(ctx context.Context, t *testing.T) {
		// Navigate to login page
		err := chromedp.Run(ctx,
			chromedp.Navigate(testServer.URL+"/login"),
			chromedp.WaitVisible(`form[action="/login"]`),
		)
		
		if err != nil {
			t.Fatalf("Failed to navigate to login page: %v", err)
		}
		
		// Verify login form elements
		AssertElementExists(ctx, t, `input[name="email"]`)
		AssertElementExists(ctx, t, `input[name="password"]`)
		AssertElementExists(ctx, t, `button[type="submit"]`)
		
		// Attempt login with invalid credentials
		err = chromedp.Run(ctx,
			chromedp.Clear(`input[name="email"]`),
			chromedp.SendKeys(`input[name="email"]`, "wrong@example.com"),
			chromedp.SendKeys(`input[name="password"]`, "wrongpassword"),
			chromedp.Submit(`form[action="/login"]`),
			// Wait for page to load after submission
			chromedp.WaitVisible(`body`),
		)
		
		if err != nil {
			t.Fatalf("Failed to submit login form with invalid credentials: %v", err)
		}
		
		// Check for error message
		AssertTextContains(ctx, t, `body`, "Invalid credentials")
		
		// Login with valid credentials
		err = chromedp.Run(ctx,
			chromedp.Clear(`input[name="email"]`),
			chromedp.SendKeys(`input[name="email"]`, testServer.Username),
			chromedp.SendKeys(`input[name="password"]`, testServer.Password),
			chromedp.Submit(`form[action="/login"]`),
			// Wait for redirect to complete and dashboard to load
			chromedp.WaitVisible(`a[href="/logout"]`),
		)
		
		if err != nil {
			t.Fatalf("Failed to login with valid credentials: %v", err)
		}
		
		// Verify we're on the dashboard
		AssertElementExists(ctx, t, `.dashboard-container`)
	})
}
