package database

import (
	"testing"

	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/dbtest"
)

func TestCreateUser(t *testing.T) {
	// Setup test database
	db := dbtest.SetupTestDB(t, GetInitSQL())
	defer db.TeardownTestDB(t)

	// Clean up tables before test
	db.CleanupTables(t)

	// Test data
	email := "test@example.com"
	name := "Test User"
	password := "hashed_password"
	isAdmin := true
	isInstanceAdmin := false

	// Get the next group ID
	nextGroupId, err := GetNextAvailableGroupId()
	if err != nil {
		t.Fatalf("Failed to get next available group ID: %v", err)
	}

	// Create user using RegisterUser
	userRegister := constants.UserRegister{
		GroupId:       nextGroupId,
		Email:         email,
		Name:          name,
		Password:      password,
		Admin:         isAdmin,
		InstanceAdmin: isInstanceAdmin,
	}

	// Execute register user function
	err = RegisterUser(userRegister)

	// Assertions
	if err != nil {
		t.Fatalf("RegisterUser returned error: %v", err)
	}

	// Get the created user
	user, err := GetUser(email)
	if err != nil {
		t.Fatalf("GetUser returned error: %v", err)
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
}

func TestGetUserByEmail(t *testing.T) {
	// Setup test database
	db := dbtest.SetupTestDB(t, GetInitSQL())
	defer db.TeardownTestDB(t)

	// Clean up tables before test
	db.CleanupTables(t)

	// Test data
	email := "test@example.com"
	name := "Test User"
	password := "hashed_password"
	isAdmin := true
	isInstanceAdmin := false

	// Get the next group ID
	nextGroupId, err := GetNextAvailableGroupId()
	if err != nil {
		t.Fatalf("Failed to get next available group ID: %v", err)
	}

	// Create a user first
	userRegister := constants.UserRegister{
		GroupId:       nextGroupId,
		Email:         email,
		Name:          name,
		Password:      password,
		Admin:         isAdmin,
		InstanceAdmin: isInstanceAdmin,
	}

	err = RegisterUser(userRegister)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test getting user by email
	userId, err := GetUserByEmail(email)

	// Assertions
	if err != nil {
		t.Fatalf("GetUserByEmail returned error: %v", err)
	}

	if userId == -1 {
		t.Fatal("Expected user ID to be retrieved, got -1")
	}

	// Get the complete user object
	user, err := GetUser(email)
	if err != nil {
		t.Fatalf("GetUser returned error: %v", err)
	}

	if user == nil {
		t.Fatal("Expected user to be retrieved, got nil")
	}

	if user.Email != email {
		t.Errorf("Expected email %s, got %s", email, user.Email)
	}

	// Test with non-existent email
	nonExistentUserId, err := GetUserByEmail("nonexistent@example.com")

	if err != nil {
		t.Errorf("Expected nil error for non-existent user, got: %v", err)
	}

	if nonExistentUserId != -1 {
		t.Errorf("Expected -1 for non-existent user, got: %d", nonExistentUserId)
	}

	nonExistentUser, err := GetUser("nonexistent@example.com")
	if err != nil {
		t.Errorf("Expected nil error for non-existent user, got: %v", err)
	}

	if nonExistentUser != nil {
		t.Errorf("Expected nil for non-existent user, got: %+v", nonExistentUser)
	}
}
