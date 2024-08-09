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

func CreateItem(item constants.ItemInsert) error {
	db := getDatabaseConnection()
	defer db.Close()
	_, err := db.Exec(`INSERT INTO listaway.item (Name, ListId, URL, Notes, Priority) VALUES($1, $2, $3, $4, $5)`, item.Name, item.ListId, item.URL, item.Notes, item.Priority)
	return err
}

func DeleteItem(itemId int) error {
	db := getDatabaseConnection()
	defer db.Close()
	_, err := db.Exec(`DELETE FROM listaway.item WHERE Id = $1`, itemId)
	return err
}

func GetItem(itemId int) (constants.Item, error) {
	db := getDatabaseConnection()
	defer db.Close()
	row := db.QueryRow("SELECT Id, Name, URL, Notes, Priority FROM "+constants.DB_TABLE_ITEM+" WHERE Id = $1", itemId)
	var item constants.Item
	err := row.Scan(&item.Id, &item.Name, &item.URL, &item.Notes, &item.Priority)
	if err != nil {
		return constants.Item{}, err
	}
	return item, nil
}

func UpdateItem(itemId int, item constants.ItemInsert) error {
	db := getDatabaseConnection()
	defer db.Close()
	_, err := db.Exec(`UPDATE listaway.item SET Name = $1, URL = $2, Priority = $3, Notes = $4 WHERE Id = $5`, item.Name, item.URL, item.Priority, item.Notes, itemId)
	return err
}
