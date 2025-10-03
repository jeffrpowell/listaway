package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jeffrpowell/listaway/internal/constants"
)

// GetUserByOIDC retrieves a user by OIDC provider and subject
func GetUserByOIDC(provider, subject string) (int, error) {
	var userId int
	query := fmt.Sprintf("SELECT id FROM %s WHERE oidc_provider = $1 AND oidc_subject = $2", constants.DB_TABLE_USER)
	db := getDatabaseConnection()
	defer db.Close()
	err := db.QueryRow(query, provider, subject).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil // User not found
		}
		return -1, fmt.Errorf("error querying user by OIDC: %v", err)
	}

	return userId, nil
}

// CreateOIDCUser creates a new user with OIDC authentication
func CreateOIDCUser(email, name, provider, subject, oidcEmail string) (int, error) {
	// Get next available group ID for new OIDC user
	groupId, err := GetNextAvailableGroupId()
	if err != nil {
		return -1, fmt.Errorf("error getting next group ID: %v", err)
	}

	var userId int
	query := fmt.Sprintf(`
		INSERT INTO %s (groupid, email, name, passwordhash, admin, instanceadmin, oidc_provider, oidc_subject, oidc_email) 
		VALUES ($1, $2, $3, '', true, false, $4, $5, $6) 
		RETURNING id`, constants.DB_TABLE_USER)
	db := getDatabaseConnection()
	defer db.Close()
	err = db.QueryRow(query, groupId, email, name, provider, subject, oidcEmail).Scan(&userId)
	if err != nil {
		return -1, fmt.Errorf("error creating OIDC user: %v", err)
	}

	log.Printf("Created new OIDC user with ID %d (group %d) for provider %s", userId, groupId, provider)
	return userId, nil
}

// LinkOIDCToExistingUser links OIDC authentication to an existing user account
func LinkOIDCToExistingUser(userID int, provider, subject, oidcEmail string) error {
	query := fmt.Sprintf(`
		UPDATE %s 
		SET oidc_provider = $1, oidc_subject = $2, oidc_email = $3 
		WHERE id = $4`, constants.DB_TABLE_USER)
	db := getDatabaseConnection()
	defer db.Close()
	result, err := db.Exec(query, provider, subject, oidcEmail, userID)
	if err != nil {
		return fmt.Errorf("error linking OIDC to user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found for linking")
	}

	log.Printf("Linked OIDC provider %s to user ID %d", provider, userID)
	return nil
}

// GetUserByEmailForOIDCLinking retrieves a user by email for OIDC account linking
func GetUserByEmailForOIDCLinking(email string) (int, bool, error) {
	var userId int
	var hasOIDC bool
	query := fmt.Sprintf(`
		SELECT id, 
		       CASE WHEN oidc_provider IS NOT NULL AND oidc_subject IS NOT NULL THEN true ELSE false END as has_oidc
		FROM %s 
		WHERE email = $1`, constants.DB_TABLE_USER)
	db := getDatabaseConnection()
	defer db.Close()
	err := db.QueryRow(query, email).Scan(&userId, &hasOIDC)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, false, nil // User not found
		}
		return -1, false, fmt.Errorf("error querying user by email: %v", err)
	}

	return userId, hasOIDC, nil
}

// CreateOrUpdateOIDCUser creates a new OIDC user or updates existing one
func CreateOrUpdateOIDCUser(email, name, provider, subject, oidcEmail string) (int, error) {
	// First, check if user exists with this OIDC provider/subject
	existingUserID, err := GetUserByOIDC(provider, subject)
	if err != nil {
		return -1, fmt.Errorf("error checking existing OIDC user: %v", err)
	}
	db := getDatabaseConnection()
	defer db.Close()
	if existingUserID != -1 {
		// User exists with this OIDC identity, update their info
		query := fmt.Sprintf(`
			UPDATE %s 
			SET name = $1, oidc_email = $2 
			WHERE id = $3`, constants.DB_TABLE_USER)

		_, err := db.Exec(query, name, oidcEmail, existingUserID)
		if err != nil {
			return -1, fmt.Errorf("error updating existing OIDC user: %v", err)
		}

		log.Printf("Updated existing OIDC user with ID %d", existingUserID)
		return existingUserID, nil
	}

	// Check if user exists with this email but no OIDC
	emailUserID, hasOIDC, err := GetUserByEmailForOIDCLinking(email)
	if err != nil {
		return -1, fmt.Errorf("error checking user by email: %v", err)
	}

	if emailUserID != -1 && !hasOIDC {
		// User exists with email but no OIDC, link the accounts
		err := LinkOIDCToExistingUser(emailUserID, provider, subject, oidcEmail)
		if err != nil {
			return -1, fmt.Errorf("error linking OIDC to existing user: %v", err)
		}
		return emailUserID, nil
	}

	// Create new user
	return CreateOIDCUser(email, name, provider, subject, oidcEmail)
}

// UnlinkOIDCFromUser removes OIDC authentication from a user account
func UnlinkOIDCFromUser(userID int) error {
	query := fmt.Sprintf(`
		UPDATE %s 
		SET oidc_provider = NULL, oidc_subject = NULL, oidc_email = NULL 
		WHERE id = $1`, constants.DB_TABLE_USER)
	db := getDatabaseConnection()
	defer db.Close()
	result, err := db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("error unlinking OIDC from user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found for unlinking")
	}

	log.Printf("Unlinked OIDC from user ID %d", userID)
	return nil
}

// GetUserOIDCInfo retrieves OIDC information for a user
func GetUserOIDCInfo(userID int) (string, string, string, error) {
	var provider, subject, oidcEmail sql.NullString
	query := fmt.Sprintf(`
		SELECT oidc_provider, oidc_subject, oidc_email 
		FROM %s 
		WHERE id = $1`, constants.DB_TABLE_USER)
	db := getDatabaseConnection()
	defer db.Close()
	err := db.QueryRow(query, userID).Scan(&provider, &subject, &oidcEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", "", fmt.Errorf("user not found")
		}
		return "", "", "", fmt.Errorf("error querying user OIDC info: %v", err)
	}

	return provider.String, subject.String, oidcEmail.String, nil
}
