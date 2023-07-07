package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var Dsn string
var Db *sql.DB

func InitDb() {
	var envs map[string]string

	envs, err := godotenv.Read(`D:\Dolgaim\Programozás\Golang Learning\Comicbook-Webshop\Backend\.env`)
	if err != nil {
		log.Fatal(err)
	}

	Dsn = envs["MySQLUsername"] + ":" + envs["MySQLPassword"] + "@tcp(" + envs["MyAddress"] + ":" + envs["MyPort"] + ")/comicbooks"

	Db, err = sql.Open("mysql", Dsn)

	if err != nil {
		log.Fatal(err)
		return
	}

	err = Db.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Connected to MySQL")
	return
}
