package database

import (
	"database/sql"

	"github.com/jeffrpowell/listaway/internal/constants"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func AdminUserExists() bool {
	db := getDatabaseConnection()
	defer db.Close()
	row := db.QueryRow("SELECT COUNT(1) FROM " + constants.DB_TABLE_USER + " WHERE Admin = true")

	var numAdmins int
	err := row.Scan(&numAdmins)
	if err != nil {
		return true //Assume safest option
	}
	return numAdmins > 0
}

func RegisterUser(user constants.UserRegister) error {
	hash, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	hashedUser := constants.UserRegister{
		Email:    user.Email,
		Name:     user.Name,
		Password: hash,
		Admin:    user.Admin,
	}
	db := getDatabaseConnection()
	defer db.Close()
	_, err = db.Exec("INSERT INTO "+constants.DB_TABLE_USER+" (Email, Name, PasswordHash, Admin) VALUES ($1, $2, $3, $4)",
		hashedUser.Email, hashedUser.Name, hashedUser.Password, hashedUser.Admin)
	if err != nil {
		return err
	}
	return nil
}

func LoginUser(email, password string) (bool, error) {
	db := getDatabaseConnection()
	defer db.Close()
	row := db.QueryRow("SELECT PasswordHash FROM "+constants.DB_TABLE_USER+" WHERE Email = $1", email)
	var passwordHash string
	err := row.Scan(&passwordHash)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		} else {
			return false, nil
		}
	}
	return checkPasswordHash(password, passwordHash), nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
