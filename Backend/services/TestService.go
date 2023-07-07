package services

import (
	"fmt"
	"log"
	db "main/database"
	"main/models"
	"net/http"
)

func GetAllUserTest(w http.ResponseWriter, r *http.Request) {
	row, err := db.Db.Query("SELECT id, name, password, email FROM users")
	defer row.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	var users []models.User

	for row.Next() {
		var u models.User

		err = row.Scan(&u.ID, &u.Username, &u.PswHash, &u.Email)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}
	SendResponse(w, 200, users)
}
