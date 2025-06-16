package database

import (
	"testing"

	testhelper "github.com/jeffrpowell/listaway/internal/testing"
)

func TestCreateUser(t *testing.T) {
	// Setup test database
	db := testhelper.SetupTestDB(t)
	defer db.TeardownTestDB(t)
	
	// Clean up tables before test
	db.CleanupTables(t)
	
	// Test data
	email := "test@example.com"
	name := "Test User"
	passwordHash := "hashed_password"
	isAdmin := true
	isInstanceAdmin := false
	
	// Execute create user function
	user, err := CreateUser(email, name, passwordHash, isAdmin, isInstanceAdmin)
	
	// Assertions
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}
	
	if user == nil {
		t.Fatal("Expected user to be created, got nil")
	}
	
	if user.Email != email {
		t.Errorf("Expected email %s, got %s", email, user.Email)
	}
	
	if user.Name != name {
		t.Errorf("Expected name %s, got %s", name, user.Name)
	}
	
	if user.Admin != isAdmin {
		t.Errorf("Expected admin %v, got %v", isAdmin, user.Admin)
	}
	
	if user.InstanceAdmin != isInstanceAdmin {
		t.Errorf("Expected instance admin %v, got %v", isInstanceAdmin, user.InstanceAdmin)
	}
	
	// Verify user was saved in database
	savedUser, err := GetUserByEmail(email)
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}
	
	if savedUser == nil {
		t.Fatal("User not found in database")
	}
	
	if savedUser.Email != email {
		t.Errorf("Expected saved email %s, got %s", email, savedUser.Email)
	}
}

func TestGetUserByEmail(t *testing.T) {
	// Setup test database
	db := testhelper.SetupTestDB(t)
	defer db.TeardownTestDB(t)
	
	// Clean up tables before test
	db.CleanupTables(t)
	
	// Test data
	email := "test@example.com"
	name := "Test User"
	passwordHash := "hashed_password"
	isAdmin := true
	isInstanceAdmin := false
	
	// Create a user first
	_, err := CreateUser(email, name, passwordHash, isAdmin, isInstanceAdmin)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	
	// Test getting user by email
	user, err := GetUserByEmail(email)
	
	// Assertions
	if err != nil {
		t.Fatalf("GetUserByEmail returned error: %v", err)
	}
	
	if user == nil {
		t.Fatal("Expected user to be retrieved, got nil")
	}
	
	if user.Email != email {
		t.Errorf("Expected email %s, got %s", email, user.Email)
	}
	
	// Test with non-existent email
	nonExistentUser, err := GetUserByEmail("nonexistent@example.com")
	
	if err != nil {
		t.Errorf("Expected nil error for non-existent user, got: %v", err)
	}
	
	if nonExistentUser != nil {
		t.Errorf("Expected nil for non-existent user, got: %+v", nonExistentUser)
	}
}
