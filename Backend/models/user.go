package models

import (
	"database/sql"
	"errors"
	"fmt"
	db "main/database"
)

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"name"`
	Email     string `json:"email"`
	PswHash   string
	CreatedAt string `json:"created_at"`
	Active    bool
	VerHash   string
}

func (u *User) SelectById() error {
	stmt := "SELECT `id`, `Username`, `Email`, `PswHash`, `CreatedAt`, `Active`, `VerHash` FROM USERS WHERE `id`=?"
	row := db.Db.QueryRow(stmt, &u.ID)
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PswHash, &u.CreatedAt, &u.Active, &u.VerHash)
	if err != nil {
		return err
	}
	return err
}

func (u *User) GetUserByUsername() error {

	stmt := "SELECT u.id, u.Username, u.Email, u.PswHash, u.CreatedAt, u.Active, uvh.verHash " +
		"FROM users AS u INNER JOIN user_email_ver_hash AS uvh ON u.username = uvh.username " +
		"WHERE u.Username = ?"
	row := db.Db.QueryRow(stmt, u.Username)
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PswHash, &u.CreatedAt, &u.Active, &u.VerHash)
	if err != nil {
		fmt.Println("getUser() error selecting User, err: ", err)
		return err
	}
	return nil
}

func (u *User) MakeActive() error {
	stmt, err := db.Db.Prepare("UPDATE USERS SET `Active`=TRUE WHERE `id`=?")
	if err != nil {
		fmt.Println("error preparing statement to update Active")
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(u.ID)
	if err != nil {
		fmt.Println("error executing statemnt to update Active")
		return err
	}
	return nil
}

func (u *User) UpdateUser() error {
	var updateUserStmt *sql.Stmt
	updateUserStmt, err := db.Db.Prepare("UPDATE USERS SET `Username`=?, `Email`=?, `PswHash`=?, `Active`=? WHERE `id`=?;")
	if err != nil {
		fmt.Println("error preparring statement to update user in Db with Update, err:", err)
		return err
	}
	defer updateUserStmt.Close()
	var result sql.Result

	result, err = updateUserStmt.Exec(u.Username, u.Email, u.PswHash, u.Active, u.ID)

	rowsAff, _ := result.RowsAffected()

	if err != nil {
		fmt.Println("there was an erorr updating user in Update() err:", err)
		return errors.New("number of rows affected not equal to one")
	}
	if rowsAff != 1 {
		fmt.Println("rows affected not equal to one:", err)
		return errors.New("number of rows affected not equal to one")
	}
	return err
}
