package database

import (
	"log"

	"github.com/jeffrpowell/listaway/internal/constants"
	_ "github.com/lib/pq"
)

func GetLists() []constants.List {
	db := getDatabaseConnection()
	defer db.Close()
	rows, err := db.Query("SELECT Id, Name, ShareURL FROM " + constants.DB_TABLE_LIST)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var lists []constants.List
	for rows.Next() {
		var l constants.List

		err := rows.Scan(&l.Id, &l.Name, &l.ShareURL)
		if err != nil {
			log.Fatal(err)
		}
		lists = append(lists, l)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return lists
}
