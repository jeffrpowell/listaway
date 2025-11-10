package database

import (
	"database/sql"

	"github.com/jeffrpowell/listaway/internal/constants"
	_ "github.com/lib/pq"
)

// GetGroupSharingEnabled returns whether group sharing is enabled for a given group
func GetGroupSharingEnabled(groupId int) (bool, error) {
	db := getDatabaseConnection()
	defer db.Close()
	
	var enabled bool
	err := db.QueryRow("SELECT group_sharing_enabled FROM "+constants.DB_TABLE_GROUP_SETTINGS+" WHERE groupid = $1", groupId).Scan(&enabled)
	
	if err == sql.ErrNoRows {
		// Group settings don't exist yet, default is false
		return false, nil
	}
	if err != nil {
		return false, err
	}
	
	return enabled, nil
}

// SetGroupSharingEnabled sets whether group sharing is enabled for a given group
func SetGroupSharingEnabled(groupId int, enabled bool) error {
	db := getDatabaseConnection()
	defer db.Close()
	
	// Use UPSERT to handle cases where the group settings don't exist yet
	_, err := db.Exec(`
		INSERT INTO `+constants.DB_TABLE_GROUP_SETTINGS+` (groupid, group_sharing_enabled) 
		VALUES ($1, $2)
		ON CONFLICT (groupid) 
		DO UPDATE SET group_sharing_enabled = $2
	`, groupId, enabled)
	
	return err
}

// GetListsSharedWithGroup returns all lists shared with a user's group
// Returns list ID, list name, owner user ID, owner name, and whether the group can edit
func GetListsSharedWithGroup(userId int) ([]constants.ListSharedWithGroup, error) {
	db := getDatabaseConnection()
	defer db.Close()
	
	rows, err := db.Query(`
		SELECT l.id, l.name, l.description, l.userid, u.name, l.group_can_edit
		FROM `+constants.DB_TABLE_LIST+` l
		JOIN `+constants.DB_TABLE_USER+` u ON l.userid = u.id
		WHERE l.share_with_group = true
		AND u.groupid = (SELECT groupid FROM `+constants.DB_TABLE_USER+` WHERE id = $1)
		AND l.userid != $1
		ORDER BY u.name, l.name
	`, userId)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var lists []constants.ListSharedWithGroup
	for rows.Next() {
		var list constants.ListSharedWithGroup
		err := rows.Scan(&list.Id, &list.Name, &list.Description, &list.OwnerId, &list.OwnerName, &list.GroupCanEdit)
		if err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return lists, nil
}

// SetListGroupSharing sets the group sharing settings for a list
func SetListGroupSharing(listId int, shareWithGroup bool, groupCanEdit bool) error {
	db := getDatabaseConnection()
	defer db.Close()
	
	_, err := db.Exec(`
		UPDATE `+constants.DB_TABLE_LIST+` 
		SET share_with_group = $1, group_can_edit = $2 
		WHERE id = $3
	`, shareWithGroup, groupCanEdit, listId)
	
	return err
}

// GetListGroupSharing returns the group sharing settings for a list
func GetListGroupSharing(listId int) (shareWithGroup bool, groupCanEdit bool, err error) {
	db := getDatabaseConnection()
	defer db.Close()
	
	err = db.QueryRow(`
		SELECT share_with_group, group_can_edit 
		FROM `+constants.DB_TABLE_LIST+` 
		WHERE id = $1
	`, listId).Scan(&shareWithGroup, &groupCanEdit)
	
	if err != nil {
		return false, false, err
	}
	
	return shareWithGroup, groupCanEdit, nil
}

// UserCanEditList returns true if the user can edit the list
// This is true if the user owns the list OR if the list is shared with their group with edit permissions
func UserCanEditList(userId int, listId int) (bool, error) {
	db := getDatabaseConnection()
	defer db.Close()
	
	// Check if user owns the list
	owns, err := UserOwnsList(userId, listId)
	if err != nil {
		return false, err
	}
	if owns {
		return true, nil
	}
	
	// Check if the list is shared with the user's group with edit permissions
	var canEdit bool
	err = db.QueryRow(`
		SELECT l.group_can_edit
		FROM `+constants.DB_TABLE_LIST+` l
		JOIN `+constants.DB_TABLE_USER+` u ON l.userid = u.id
		WHERE l.id = $1
		AND l.share_with_group = true
		AND u.groupid = (SELECT groupid FROM `+constants.DB_TABLE_USER+` WHERE id = $2)
	`, listId, userId).Scan(&canEdit)
	
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	
	return canEdit, nil
}

// UserCanViewList returns true if the user can view the list
// This is true if the user owns the list OR if the list is shared with their group
func UserCanViewList(userId int, listId int) (bool, error) {
	db := getDatabaseConnection()
	defer db.Close()
	
	// Check if user owns the list
	owns, err := UserOwnsList(userId, listId)
	if err != nil {
		return false, err
	}
	if owns {
		return true, nil
	}
	
	// Check if the list is shared with the user's group
	var count int
	err = db.QueryRow(`
		SELECT COUNT(1)
		FROM `+constants.DB_TABLE_LIST+` l
		JOIN `+constants.DB_TABLE_USER+` u ON l.userid = u.id
		WHERE l.id = $1
		AND l.share_with_group = true
		AND u.groupid = (SELECT groupid FROM `+constants.DB_TABLE_USER+` WHERE id = $2)
	`, listId, userId).Scan(&count)
	
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}
