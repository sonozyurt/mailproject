package datamanagement

import (
	"database/sql"
	"time"
)

var db *sql.DB

const dbTimeout = time.Second * 3

func SetUpUsers(dbPool *sql.DB) User {
	db = dbPool

	return User{}
}
