package constants

import "database/sql"

type UserRead struct {
	Id    uint64
	Email string
	Name  string
	Admin bool
}

type UserRegister struct {
	Email    string
	Name     string
	Password string
	Admin    bool
}

type List struct {
	Id          uint64
	Name        string
	Description sql.NullString
	ShareCode   sql.NullString
}

type ListPostParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ItemInsert struct {
	Name     string
	ListId   uint64
	URL      sql.NullString
	Notes    sql.NullString
	Priority sql.NullInt64
}

type Item struct {
	Id       uint64
	Name     string
	URL      sql.NullString
	Notes    sql.NullString
	Priority sql.NullInt64
}
