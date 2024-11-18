package constants

import "database/sql"

type UserRead struct {
	Id            uint64
	GroupId       uint64
	Email         string
	Name          string
	Admin         bool
	InstanceAdmin bool
}

type UserRegister struct {
	GroupId       int
	Email         string
	Name          string
	Password      string
	Admin         bool
	InstanceAdmin bool
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
	Id       uint64         `json:"id"`
	Name     string         `json:"name"`
	URL      sql.NullString `json:"url"`
	Priority sql.NullInt64  `json:"priority"`
	Notes    sql.NullString `json:"notes"`
}
