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

func CreateItem(item constants.ItemInsert) (int, error) {
	db := getDatabaseConnection()
	defer db.Close()
	var newId int
	err := db.QueryRow(`INSERT INTO listaway.item (Name, ListId, URL, Notes, Priority) VALUES($1, $2, $3, $4, $5) RETURNING Id`, item.Name, item.ListId, item.URL, item.Notes, item.Priority).Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, nil
}
