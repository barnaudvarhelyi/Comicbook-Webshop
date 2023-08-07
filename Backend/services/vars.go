package services

import (
	"database/sql"
	database "main/database"
)

var db *sql.DB

func Init() {
	db = database.Db
}
