package database

import (
	"github.com/jeffrpowell/listaway/internal/constants"
	_ "github.com/lib/pq"
)

func GetLists(userId int) ([]constants.List, error) {
	db := getDatabaseConnection()
	defer db.Close()
	rows, err := db.Query("SELECT Id, Name, ShareURL FROM "+constants.DB_TABLE_LIST+" WHERE UserId = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []constants.List
	for rows.Next() {
		var l constants.List

		err := rows.Scan(&l.Id, &l.Name, &l.ShareURL)
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
	row := db.QueryRow("SELECT COUNT(1) FROM "+constants.DB_TABLE_LIST+" WHERE UserId = $1 AND Name = $2", userId, name)
	var matches int
	err := row.Scan(&matches)
	if err != nil {
		return false, err
	}
	return matches != 0, nil
}

func CreateList(userId int, name string) (int, error) {
	db := getDatabaseConnection()
	defer db.Close()
	var newId int
	err := db.QueryRow(`INSERT INTO listaway.list (UserId, Name) VALUES($1, $2) RETURNING Id`, userId, name).Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, nil
}
