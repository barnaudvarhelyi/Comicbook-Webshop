package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func InitDb() {
	var envs map[string]string

	envs, err := godotenv.Read(`D:\Dolgaim\Programozás\Golang Learning\Comicbook-Webshop\Backend\.env`)
	if err != nil {
		log.Fatal(err)
	}

	dsn := envs["MySQLUsername"] + ":" + envs["MySQLPassword"] + "@tcp(" + envs["MyAddress"] + ":" + envs["MyPort"] + ")/comicbooks"

	db, err = sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err)
		return
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Connected to MySQL")
}
