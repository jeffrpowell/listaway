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

func LoginUser(email, password string) (int, error) {
	db := getDatabaseConnection()
	defer db.Close()
	row := db.QueryRow("SELECT Id, PasswordHash FROM "+constants.DB_TABLE_USER+" WHERE Email = $1", email)
	var userId int
	var passwordHash string
	err := row.Scan(&userId, &passwordHash)
	if err != nil {
		if err != sql.ErrNoRows {
			return -1, err
		} else {
			return -1, nil
		}
	}
	if checkPasswordHash(password, passwordHash) {
		return userId, nil
	}
	return -1, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func UserIsAdmin(userId int) (bool, error) {
	db := getDatabaseConnection()
	defer db.Close()
	row := db.QueryRow("SELECT Admin FROM "+constants.DB_TABLE_USER+" WHERE Id = $1", userId)
	var admin bool
	err := row.Scan(&admin)
	if err != nil {
		return false, err
	}
	return admin, nil
}

func GetAllUsers() ([]constants.UserRead, error) {
	db := getDatabaseConnection()
	defer db.Close()
	rows, err := db.Query("SELECT Id, Email, Name, Admin FROM " + constants.DB_TABLE_USER)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []constants.UserRead
	for rows.Next() {
		var user constants.UserRead

		err := rows.Scan(&user.Id, &user.Email, &user.Name, &user.Admin)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
