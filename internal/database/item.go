package database

import (
	"github.com/jeffrpowell/listaway/internal/constants"
	_ "github.com/lib/pq"
)

func GetListItems(listId int) ([]constants.Item, error) {
	db := getDatabaseConnection()
	defer db.Close()
	rows, err := db.Query("SELECT Id, Name, URL, Priority, Notes FROM "+constants.DB_TABLE_ITEM+" WHERE ListId = $1", listId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []constants.Item
	for rows.Next() {
		var i constants.Item

		err := rows.Scan(&i.Id, &i.Name, &i.URL, &i.Priority, &i.Notes)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
