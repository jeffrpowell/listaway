package database

import (
	"database/sql"

	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/constants/random"
	_ "github.com/lib/pq"
)

func GetLists(userId int) ([]constants.List, error) {
	db := getDatabaseConnection()
	defer db.Close()
	rows, err := db.Query("SELECT id, name, description, shareCode FROM "+constants.DB_TABLE_LIST+" WHERE userId = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []constants.List
	for rows.Next() {
		var l constants.List

		err := rows.Scan(&l.Id, &l.Name, &l.Description, &l.ShareCode)
		if err != nil {
			return nil, err
		}
		lists = append(lists, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lists, nil
}

func ListNameTaken(userId int, name string) (bool, error) {
	db := getDatabaseConnection()
	defer db.Close()
	row := db.QueryRow("SELECT COUNT(1) FROM "+constants.DB_TABLE_LIST+" WHERE userId = $1 AND name = $2", userId, name)
	var matches int
	err := row.Scan(&matches)
	if err != nil {
		return false, err
	}
	return matches != 0, nil
}

func CreateList(userId int, name string, description string) (int, error) {
	db := getDatabaseConnection()
	defer db.Close()
	var newId int
	err := db.QueryRow(`INSERT INTO listaway.list (userId, name, description) VALUES($1, $2, $3) RETURNING id`, userId, name, description).Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, nil
}

func UserOwnsList(userId int, listId int) (bool, error) {
	db := getDatabaseConnection()
	defer db.Close()
	row := db.QueryRow("SELECT COUNT(1) FROM "+constants.DB_TABLE_LIST+" WHERE userId = $1 AND id = $2", userId, listId)
	var matches int
	err := row.Scan(&matches)
	if err != nil {
		return false, err
	}
	return matches != 0, nil
}

func GetList(listId int) (constants.List, error) {
	db := getDatabaseConnection()
	defer db.Close()
	row := db.QueryRow("SELECT id, name, description, sharecode FROM "+constants.DB_TABLE_LIST+" WHERE id = $1", listId)
	var list constants.List
	err := row.Scan(&list.Id, &list.Name, &list.Description, &list.ShareCode)
	if err != nil {
		return constants.List{}, err
	}
	return list, nil
}

func UpdateList(listId int, params constants.ListPostParams) error {
	db := getDatabaseConnection()
	defer db.Close()
	_, err := db.Exec(`UPDATE listaway.list SET name = $1, description = $2 WHERE id = $3`, params.Name, params.Description, listId)
	return err
}

func DeleteList(listId int, confirmationName string) (bool, error) {
	db := getDatabaseConnection()
	defer db.Close()
	
	// Start a transaction to ensure atomicity
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}
	
	// Verify the list exists and the name matches
	matches, err := confirmNameMatchesListName(db, listId, confirmationName)
	if err != nil {
		tx.Rollback()
		return false, err
	}
	if !matches {
		tx.Rollback()
		return false, nil
	}
	
	// Delete all items associated with this list first
	_, err = tx.Exec("DELETE FROM "+constants.DB_TABLE_ITEM+" WHERE listid = $1", listId)
	if err != nil {
		tx.Rollback()
		return false, err
	}
	
	// Then delete the list itself
	_, err = tx.Exec("DELETE FROM "+constants.DB_TABLE_LIST+" WHERE id = $1 AND name = $2", listId, confirmationName)
	if err != nil {
		tx.Rollback()
		return false, err
	}
	
	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return false, err
	}
	
	return true, nil
}

func confirmNameMatchesListName(db *sql.DB, listId int, confirmationName string) (bool, error) {
	row := db.QueryRow("SELECT COUNT(1) FROM "+constants.DB_TABLE_LIST+" WHERE id = $1 AND name = $2", listId, confirmationName)
	var matches int
	err := row.Scan(&matches)
	if err != nil {
		return false, err
	}
	return matches != 0, nil
}

func GenerateShareCode(listId int) (string, error) {
	db := getDatabaseConnection()
	defer db.Close()
	code, err := createUniqueShareCode(db)
	if err != nil {
		return "", err
	}
	_, err = db.Exec(`UPDATE listaway.list SET sharecode = $1 WHERE id = $2`, code, listId)
	return code, err
}

func createUniqueShareCode(db *sql.DB) (string, error) {
	var code string
	var err error
	var count int

	for {
		// Generate a random string
		code, err = random.String(constants.DefaultN, constants.CHARSET_UNAMBIGUOUS)
		if err != nil {
			return "", err
		}

		// Check if the generated code already exists in the database
		row := db.QueryRow("SELECT COUNT(1) FROM "+constants.DB_TABLE_LIST+" WHERE sharecode = $1", code)
		err = row.Scan(&count)
		if err != nil {
			return "", err
		}

		// If the code doesn't exist (count is 0), it's unique, so return it
		if count == 0 {
			break
		}
	}

	return code, nil
}

func GetListFromShareCode(shareCode string) (constants.List, error) {
	db := getDatabaseConnection()
	defer db.Close()
	row := db.QueryRow("SELECT id, name, description, sharecode FROM "+constants.DB_TABLE_LIST+" WHERE sharecode = $1", shareCode)
	var list constants.List
	err := row.Scan(&list.Id, &list.Name, &list.Description, &list.ShareCode)
	if err != nil {
		return constants.List{}, err
	}
	return list, nil
}

func UnpublishShareCode(listId int) error {
	db := getDatabaseConnection()
	defer db.Close()
	_, err := db.Exec(`UPDATE listaway.list SET sharecode = NULL WHERE id = $1`, listId)
	return err
}

// GetListIdsWithShareCode retrieves all list IDs for a user that have share codes
func GetListIdsWithShareCode(userId int) ([]uint64, error) {
	db := getDatabaseConnection()
	defer db.Close()
	rows, err := db.Query("SELECT id FROM "+constants.DB_TABLE_LIST+" WHERE userId = $1 AND sharecode IS NOT NULL", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listIdsWithShareCode []uint64
	for rows.Next() {
		var id uint64
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		listIdsWithShareCode = append(listIdsWithShareCode, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return listIdsWithShareCode, nil
}
