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
	row := db.QueryRow("SELECT COUNT(1) FROM " + constants.DB_TABLE_USER + " WHERE admin = true")

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
		GroupId:       user.GroupId,
		Email:         user.Email,
		Name:          user.Name,
		Password:      hash,
		Admin:         user.Admin,
		InstanceAdmin: user.InstanceAdmin,
	}
	db := getDatabaseConnection()
	defer db.Close()
	_, err = db.Exec("INSERT INTO "+constants.DB_TABLE_USER+" (groupid, email, name, passwordhash, admin, instanceadmin) VALUES ($1, $2, $3, $4, $5, $6)",
		hashedUser.GroupId, hashedUser.Email, hashedUser.Name, hashedUser.Password, hashedUser.Admin, hashedUser.InstanceAdmin)
	if err != nil {
		return err
	}
	return nil
}

func LoginUser(email, password string) (int, error) {
	db := getDatabaseConnection()
	defer db.Close()
	row := db.QueryRow("SELECT id, passwordhash FROM "+constants.DB_TABLE_USER+" WHERE email = $1", email)
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
	row := db.QueryRow("SELECT admin FROM "+constants.DB_TABLE_USER+" WHERE id = $1", userId)
	var admin bool
	err := row.Scan(&admin)
	if err != nil {
		return false, err
	}
	return admin, nil
}

func UserIsInstanceAdmin(userId int) (bool, error) {
	db := getDatabaseConnection()
	defer db.Close()
	row := db.QueryRow("SELECT instanceadmin FROM "+constants.DB_TABLE_USER+" WHERE id = $1", userId)
	var instanceAdmin bool
	err := row.Scan(&instanceAdmin)
	if err != nil {
		return false, err
	}
	return instanceAdmin, nil
}

func GetAllUsers() ([]constants.UserRead, error) {
	db := getDatabaseConnection()
	defer db.Close()
	rows, err := db.Query("SELECT id, groupid, email, name, admin, instanceadmin FROM " + constants.DB_TABLE_USER)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []constants.UserRead
	for rows.Next() {
		var user constants.UserRead

		err := rows.Scan(&user.Id, &user.GroupId, &user.Email, &user.Name, &user.Admin, &user.InstanceAdmin)
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

func GetUsersInSameGroupAsUser(userId int) ([]constants.UserRead, error) {
	db := getDatabaseConnection()
	defer db.Close()
	rows, err := db.Query("SELECT id, groupid, email, name, admin, instanceadmin FROM "+constants.DB_TABLE_USER+" WHERE groupid = (SELECT groupid FROM listaway.user WHERE id = $1)", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []constants.UserRead
	for rows.Next() {
		var user constants.UserRead

		err := rows.Scan(&user.Id, &user.GroupId, &user.Email, &user.Name, &user.Admin, &user.InstanceAdmin)
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

func GetUserGroupId(userId int) (int, error) {
	db := getDatabaseConnection()
	defer db.Close()
	row := db.QueryRow("SELECT groupid FROM "+constants.DB_TABLE_USER+" WHERE id = $1", userId)
	var groupId int
	err := row.Scan(&groupId)
	if err != nil {
		return -1, err
	}
	return groupId, nil
}

func DeleteUser(userId int) error {
	db := getDatabaseConnection()
	defer db.Close()
	_, err := db.Exec(`DELETE FROM listaway.user WHERE id = $1`, userId)
	return err
}

func SetUserAdmin(userId int, admin bool) error {
	db := getDatabaseConnection()
	defer db.Close()
	_, err := db.Exec(`UPDATE listaway.user SET admin = $1 WHERE id = $2`, admin, userId)
	return err
}

func SetUserInstanceAdmin(userId int, instanceAdmin bool) error {
	db := getDatabaseConnection()
	defer db.Close()
	_, err := db.Exec(`UPDATE listaway.user SET instanceadmin = $1 WHERE id = $2`, instanceAdmin, userId)
	return err
}

func GetAllGroupAdmins() ([]constants.UserRead, error) {
	db := getDatabaseConnection()
	defer db.Close()
	// Get all admin users who are not instance admins (group admins)
	rows, err := db.Query("SELECT id, groupid, email, name, admin, instanceadmin FROM " + constants.DB_TABLE_USER + " WHERE admin = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []constants.UserRead
	for rows.Next() {
		var user constants.UserRead

		err := rows.Scan(&user.Id, &user.GroupId, &user.Email, &user.Name, &user.Admin, &user.InstanceAdmin)
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

// GetNextAvailableGroupId returns the next available group ID (max + 1)
func GetNextAvailableGroupId() (int, error) {
	db := getDatabaseConnection()
	defer db.Close()
	
	row := db.QueryRow("SELECT COALESCE(MAX(groupid), 0) + 1 FROM " + constants.DB_TABLE_USER)
	
	var nextGroupId int
	if err := row.Scan(&nextGroupId); err != nil {
		return -1, err
	}
	
	return nextGroupId, nil
}
