package database

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/jeffrpowell/listaway/internal/constants"
)

// CreatePasswordResetToken generates a new reset token for the given email
func CreatePasswordResetToken(email string) (string, error) {
	// Generate a secure random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)

	// Set expiry time (1 hour from now)
	createdAt := time.Now()
	expiresAt := createdAt.Add(1 * time.Hour)

	// Store in database
	db := getDatabaseConnection()
	defer db.Close()

	// Delete any existing tokens for this email
	_, err := db.Exec("DELETE FROM "+constants.DB_TABLE_RESET+" WHERE email = $1", email)
	if err != nil {
		return "", err
	}

	// Insert new token
	_, err = db.Exec(
		"INSERT INTO "+constants.DB_TABLE_RESET+" (token, email, created_at, expires_at) VALUES ($1, $2, $3, $4)",
		token, email, createdAt, expiresAt,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidatePasswordResetToken checks if a token is valid and returns the associated email
func ValidatePasswordResetToken(token string) (string, bool, error) {
	db := getDatabaseConnection()
	defer db.Close()

	// Query for token
	row := db.QueryRow(
		"SELECT email, expires_at FROM "+constants.DB_TABLE_RESET+" WHERE token = $1",
		token,
	)

	var email string
	var expiresAt time.Time
	err := row.Scan(&email, &expiresAt)
	if err != nil {
		return "", false, err
	}

	// Check if token has expired
	if time.Now().After(expiresAt) {
		// Token expired, delete it
		_, _ = db.Exec("DELETE FROM "+constants.DB_TABLE_RESET+" WHERE token = $1", token)
		return "", false, nil
	}

	return email, true, nil
}

// InvalidatePasswordResetToken removes a token from the database
func InvalidatePasswordResetToken(token string) error {
	db := getDatabaseConnection()
	defer db.Close()

	_, err := db.Exec("DELETE FROM "+constants.DB_TABLE_RESET+" WHERE token = $1", token)
	return err
}
