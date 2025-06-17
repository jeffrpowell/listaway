package database

import (
	"database/sql"

	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/constants/random"
	_ "github.com/lib/pq"
)

// GetCollections retrieves all collections owned by a specific user
func GetCollections(userId int) ([]constants.Collection, error) {
	db := getDatabaseConnection()
	defer db.Close()
	rows, err := db.Query("SELECT id, name, description, sharecode FROM listaway.collection WHERE userid = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []constants.Collection
	for rows.Next() {
		var c constants.Collection

		err := rows.Scan(&c.Id, &c.Name, &c.Description, &c.ShareCode)
		if err != nil {
			return nil, err
		}
		collections = append(collections, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return collections, nil
}

// CollectionNameTaken checks if a collection name already exists for a user
func CollectionNameTaken(userId int, name string) (bool, error) {
	db := getDatabaseConnection()
	defer db.Close()
	row := db.QueryRow("SELECT COUNT(1) FROM listaway.collection WHERE userid = $1 AND name = $2", userId, name)
	var matches int
	err := row.Scan(&matches)
	if err != nil {
		return false, err
	}
	return matches != 0, nil
}

// CreateCollection creates a new collection for a user
func CreateCollection(userId int, name string, description string) (int, error) {
	db := getDatabaseConnection()
	defer db.Close()
	var newId int
	err := db.QueryRow(`INSERT INTO listaway.collection (userid, name, description) VALUES($1, $2, $3) RETURNING id`, userId, name, description).Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, nil
}

// UserOwnsCollection checks if a user owns a specific collection
func UserOwnsCollection(userId int, collectionId int) (bool, error) {
	db := getDatabaseConnection()
	defer db.Close()
	row := db.QueryRow("SELECT COUNT(1) FROM listaway.collection WHERE userid = $1 AND id = $2", userId, collectionId)
	var matches int
	err := row.Scan(&matches)
	if err != nil {
		return false, err
	}
	return matches != 0, nil
}

// GetCollection retrieves a collection by its ID
func GetCollection(collectionId int) (constants.Collection, error) {
	db := getDatabaseConnection()
	defer db.Close()
	row := db.QueryRow("SELECT id, name, description, sharecode FROM listaway.collection WHERE id = $1", collectionId)
	var collection constants.Collection
	err := row.Scan(&collection.Id, &collection.Name, &collection.Description, &collection.ShareCode)
	if err != nil {
		return constants.Collection{}, err
	}
	return collection, nil
}

// UpdateCollection updates a collection's details
func UpdateCollection(collectionId int, params constants.CollectionPostParams) error {
	db := getDatabaseConnection()
	defer db.Close()
	_, err := db.Exec(`UPDATE listaway.collection SET name = $1, description = $2 WHERE id = $3`, params.Name, params.Description, collectionId)
	return err
}

// DeleteCollection deletes a collection after confirming the name matches
func DeleteCollection(collectionId int, confirmationName string) (bool, error) {
	db := getDatabaseConnection()
	defer db.Close()
	matches, err := confirmNameMatchesCollectionName(db, collectionId, confirmationName)
	if err != nil {
		return false, err
	}
	if !matches {
		return false, nil
	}
	_, err = db.Exec(`DELETE FROM listaway.collection WHERE id = $1 AND name = $2`, collectionId, confirmationName)
	return true, err
}

// confirmNameMatchesCollectionName verifies the provided name matches the collection name
func confirmNameMatchesCollectionName(db *sql.DB, collectionId int, confirmationName string) (bool, error) {
	row := db.QueryRow("SELECT COUNT(1) FROM listaway.collection WHERE id = $1 AND name = $2", collectionId, confirmationName)
	var matches int
	err := row.Scan(&matches)
	if err != nil {
		return false, err
	}
	return matches != 0, nil
}

// GenerateCollectionShareCode generates a unique share code for a collection
func GenerateCollectionShareCode(collectionId int) (string, error) {
	db := getDatabaseConnection()
	defer db.Close()
	code, err := createUniqueCollectionShareCode(db)
	if err != nil {
		return "", err
	}
	_, err = db.Exec(`UPDATE listaway.collection SET sharecode = $1 WHERE id = $2`, code, collectionId)
	return code, err
}

// createUniqueCollectionShareCode creates a unique share code for collections
func createUniqueCollectionShareCode(db *sql.DB) (string, error) {
	var code string
	var err error
	var count int

	for {
		// Generate a random string
		code, err = random.String(constants.DefaultN, constants.CHARSET_UNAMBIGUOUS)
		if err != nil {
			return "", err
		}

		// Check if the generated code already exists in the database for collections
		row := db.QueryRow(`SELECT COUNT(1) FROM listaway.collection WHERE sharecode = $1`, code)
		err = row.Scan(&count)
		if err != nil {
			return "", err
		}

		// Also check if the code exists for lists to avoid collisions
		row = db.QueryRow(`SELECT COUNT(1) FROM listaway.list WHERE sharecode = $1`, code)
		var listCount int
		err = row.Scan(&listCount)
		if err != nil {
			return "", err
		}

		// If the code doesn't exist in either table (count is 0), it's unique
		if count == 0 && listCount == 0 {
			break
		}
	}

	return code, nil
}

// GetCollectionFromShareCode retrieves a collection by its share code
func GetCollectionFromShareCode(shareCode string) (constants.Collection, error) {
	db := getDatabaseConnection()
	defer db.Close()
	row := db.QueryRow("SELECT id, name, description, sharecode FROM listaway.collection WHERE sharecode = $1", shareCode)
	var collection constants.Collection
	err := row.Scan(&collection.Id, &collection.Name, &collection.Description, &collection.ShareCode)
	if err != nil {
		return constants.Collection{}, err
	}
	return collection, nil
}

// UnpublishCollectionShareCode removes the share code from a collection
func UnpublishCollectionShareCode(collectionId int) error {
	db := getDatabaseConnection()
	defer db.Close()
	_, err := db.Exec(`UPDATE listaway.collection SET sharecode = NULL WHERE id = $1`, collectionId)
	return err
}

// AddListToCollection adds a list to a collection
func AddListToCollection(collectionId int, listId int) error {
	db := getDatabaseConnection()
	defer db.Close()
	_, err := db.Exec(`INSERT INTO listaway.collection_list (collectionid, listid) VALUES($1, $2) 
	ON CONFLICT (collectionid, listid) DO NOTHING`, collectionId, listId)
	return err
}

// RemoveListFromCollection removes a list from a collection
func RemoveListFromCollection(collectionId int, listId int) error {
	db := getDatabaseConnection()
	defer db.Close()
	_, err := db.Exec(`DELETE FROM listaway.collection_list WHERE collectionid = $1 AND listid = $2`, collectionId, listId)
	return err
}

// GetCollectionLists retrieves all lists in a collection
func GetCollectionLists(collectionId int) ([]constants.List, error) {
	db := getDatabaseConnection()
	defer db.Close()

	rows, err := db.Query(`
		SELECT l.id, l.name, l.description, l.sharecode 
		FROM listaway.list l
		JOIN listaway.collection_list cl ON l.id = cl.listid
		WHERE cl.collectionid = $1
		ORDER BY l.name ASC
	`, collectionId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collectionLists []constants.List
	for rows.Next() {
		var cl constants.List

		err := rows.Scan(&cl.Id, &cl.Name, &cl.Description, &cl.ShareCode)
		if err != nil {
			return nil, err
		}
		collectionLists = append(collectionLists, cl)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return collectionLists, nil
}

// GetCollectionListIds retrieves all list ids in a collection
func GetCollectionListIds(collectionId int) ([]uint64, error) {
	db := getDatabaseConnection()
	defer db.Close()

	rows, err := db.Query(`
		SELECT l.id 
		FROM listaway.list l
		JOIN listaway.collection_list cl ON l.id = cl.listid
		WHERE cl.collectionid = $1
		ORDER BY l.name ASC
	`, collectionId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collectionLists []uint64
	for rows.Next() {
		var cl uint64

		err := rows.Scan(&cl)
		if err != nil {
			return nil, err
		}
		collectionLists = append(collectionLists, cl)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return collectionLists, nil
}

// ListInCollection checks if a list is already in a collection
func ListInCollection(collectionId int, listId int) (bool, error) {
	db := getDatabaseConnection()
	defer db.Close()
	row := db.QueryRow("SELECT COUNT(1) FROM listaway.collection_list WHERE collectionid = $1 AND listid = $2", collectionId, listId)
	var matches int
	err := row.Scan(&matches)
	if err != nil {
		return false, err
	}
	return matches != 0, nil
}
