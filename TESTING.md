# Listaway Testing Guide

This document outlines the testing approach for Listaway, covering both backend Go code testing and frontend testing for the server-side rendered application.

## Testing Overview

Listaway uses a comprehensive testing approach that includes:

1. **Backend Unit Tests** - Testing individual Go packages and functions
2. **Database Integration Tests** - Testing database operations with a test database
3. **HTTP Handler Tests** - Testing API endpoints and responses
4. **Template Tests** - Testing server-side HTML templates with mock data
5. **End-to-End Tests** - Testing full user flows with a headless browser

## Setup Requirements

### Local Development Testing

1. A running PostgreSQL instance for testing
2. Create a separate test database for automated tests:
    ```sql
    --connect to your postgres server with an admin role
    CREATE DATABASE listaway_test;
    GRANT CONNECT ON DATABASE listaway_test TO listaway; --reusing existing login from plain listaway database setup in README.md
    GRANT ALL PRIVILEGES ON DATABASE listaway_test TO listaway;
    --connect to your new listaway_test database with an admin role
    CREATE SCHEMA listaway;
    GRANT CREATE, USAGE ON SCHEMA listaway to listaway;
    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA listaway TO listaway;
    GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA listaway TO listaway;
    ```
3. Environment variables set for test database connection (already configured in devcontainer):
   ```
   TEST_POSTGRES_HOST=localhost
   TEST_POSTGRES_USER=listaway
   TEST_POSTGRES_PASSWORD=password
   TEST_POSTGRES_DATABASE=listaway_test
   TEST_POSTGRES_PORT=5432
   ```
4. For E2E tests: Chrome or Chromium installed

### CI/CD Testing

The GitHub Actions workflow automatically sets up all required dependencies.

## Running Tests

### Backend Tests

Run all backend tests:

```bash
go test ./internal/...
```

Run specific package tests:

```bash
go test ./internal/database/...
go test ./internal/handlers/...
```

### Template Tests

Run template rendering tests:

```bash
go test ./web/...
```

### End-to-End Tests

E2E tests are skipped by default in regular test runs. To run them:

```bash
go test -v ./tests/e2e/...
```

To run a specific E2E test:

```bash
go test -v ./tests/e2e/ -run TestLoginFlow
```

## Writing Tests

### Backend Unit Tests

1. Create a `*_test.go` file in the same package as the code you want to test
2. Use the helper functions in `internal/testing/testhelper.go` for database testing
3. Example:
   ```go
   func TestMyFunction(t *testing.T) {
       // Setup test database if needed
       db := testhelper.SetupTestDB(t)
       defer db.TeardownTestDB(t)
       
       // Test your function
       result := MyFunction(args)
       
       // Assertions
       if result != expected {
           t.Errorf("Expected %v, got %v", expected, result)
       }
   }
   ```

### HTTP Handler Tests

1. Create a test in the handlers package
2. Use the `internal/testing/httphelper.go` functions
3. Example:
   ```go
   func TestMyHandler(t *testing.T) {
       db, recorder := testhelper.SetupTestHandler(t)
       defer db.TeardownTestDB(t)
       
       req, _ := testhelper.CreateTestRequest("GET", "/path", nil)
       
       // Call your handler
       MyHandler(recorder, req)
       
       // Assertions
       testhelper.AssertStatusCode(t, recorder.Code, http.StatusOK)
       testhelper.AssertBodyContains(t, recorder.Body.String(), "expected content")
   }
   ```

### Template Tests

1. Use the `internal/testing/templatehelper.go` functions
2. Example:
   ```go
   func TestTemplate(t *testing.T) {
       renderer := testhelper.NewTemplateRenderer(t, "./templates")
       
       testData := testhelper.MockTemplateData()
       renderer.AssertTemplateContains(t, "template.html", testData, "expected content")
   }
   ```

### End-to-End Tests

1. Create a test in the `tests/e2e` package
2. Use the chromedp helper functions
3. Example:
   ```go
   func TestUserFlow(t *testing.T) {
       testServer := NewTestServer(t)
       testServer.RunE2ETest(t, func(ctx context.Context, t *testing.T) {
           // Test steps using chromedp
           chromedp.Run(ctx, chromedp.Navigate(testServer.URL))
           AssertElementExists(ctx, t, ".element-selector")
       })
   }
   ```

## Test Database

Tests that require a database use a completely separate test database (`listaway_test`). This approach provides several benefits:

1. Complete isolation between development and test environments
2. Independent schema management without risk of cross-contamination
3. Ability to run tests without affecting development data
4. Prevention of duplicate SQL schema definitions

### Test Database Management

The test environment is configured as follows:

1. A separate dedicated PostgreSQL database (`listaway_test`)
2. Using the same schema structure as the main application (`listaway` schema)
3. Test helpers handle connections and clean up automatically

Test helpers in `internal/testing/testhelper.go` manage the database connection:

- `SetupTestDB()`: Connects to the test database and prepares it for testing
- `TeardownTestDB()`: Cleans up by closing connections
- `CleanupTables()`: Truncates tables to provide a clean slate between tests

The database initialization SQL is run during tests in the test database, ensuring consistent schema structure with production.

## Best Practices

1. **Isolate Tests** - Each test should be independent and not rely on other tests
2. **Clean Up** - Always clean up resources (database connections, files) after tests
3. **Mock External Dependencies** - Use mocks for external services
4. **Test Edge Cases** - Include tests for error conditions and edge cases
5. **Keep Tests Fast** - Optimize tests to run quickly for developer productivity

## Continuous Integration

Tests run automatically on GitHub Actions for:
- All pushes to the main branch
- All pull requests targeting the main branch

The workflow:
1. Sets up a test PostgreSQL database
2. Builds the frontend assets
3. Runs backend unit tests
4. Runs template tests

End-to-end tests are not run in CI by default due to their complexity and resource requirements.
