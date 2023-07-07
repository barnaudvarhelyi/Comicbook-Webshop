package models

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"log"
	db "main/database"
	"net/http"
	"os"
	"unicode"

	emailverifier "github.com/AfterShip/email-verifier"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"name"`
	Email     string `json:"email"`
	PswHash   string
	password  string
	CreatedAt string `json:"created_at"`
	Active    bool
	VerHash   string
}

func (u *User) GetUserByUsername() error {
	stmt := "SELECT * FROM USERS WHERE `Username`=?"
	row := db.Db.QueryRow(stmt, u.Username)
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Active, &u.VerHash)
	if err != nil {
		fmt.Println("getUser() error selecting User, err: ", err)
		return err
	}
	return nil
}

func (u *User) ValidateUsername() error {
	for _, char := range u.Username {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return errors.New("Only alphanumeric characters allowed for username")
		}
	}
	if 5 <= len(u.Username) && len(u.Username) <= 50 {
		return nil
	}
	return errors.New("username lenght must be greater than 4 and less than 51 characters")
}

func (u *User) ValidatePassword() error {
	err := passwordvalidator.Validate(u.PswHash, 60)
	return err
}

func (u *User) ValidateEmail() (statusCode int, err error) {

	var Verifier = emailverifier.NewVerifier()

	Verifier = Verifier.EnableDomainSuggest()
	Verifier = Verifier.EnableSMTPCheck()
	dispEmailsDomains := MustDispEmailDom()
	Verifier = Verifier.AddDisposableDomains(dispEmailsDomains)

	res, err := Verifier.Verify(u.Email)
	if err != nil {
		fmt.Println("verify email address failed, error: ", err)
		return http.StatusInternalServerError, err
	}
	if !res.Syntax.Valid {
		err = errors.New("email address syntax is invalid")
		fmt.Println(err)
		return http.StatusBadRequest, err
	}
	if res.Disposable {
		err = errors.New("sorry, we do not accept disposable email address")
		return http.StatusBadRequest, err
	}
	if res.Suggestion != "" {
		err = errors.New("email address is not reachtable, looking for " + res.Suggestion + " instead?")
		return http.StatusBadRequest, err
	}
	if res.Reachable == "no" {
		err = errors.New("email address is not reachable")
		return http.StatusBadRequest, err
	}
	if !res.HasMxRecords {
		err = errors.New("domain entered not properly setup to recieve emails, MX record not found")
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func (u *User) UsernameExists() (exists bool) {
	exists = true
	stmt := "SELECT `UserId` FROM USERS WHERE `Username` = ?"
	row := db.Db.QueryRow(stmt, u.Username)
	var uID string
	err := row.Scan(&uID)
	if err == sql.ErrNoRows {
		return false
	}
	return exists
}

func (u *User) MakeActive() error {
	stmt, err := db.Db.Prepare("UPDATE USERS SET `Active`=TRUE WHERE `UserId`=?")
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

func (u *User) VerifyPswd() error {
	err := bcrypt.CompareHashAndPassword([]byte(u.PswHash), []byte(u.password))
	if err != nil {
		err = errors.New("Username and password do not match!")
		return err
	}

	if u.Active {
		err = errors.New("User email not verified yet!")
		return err
	}
	return nil
}

func (u *User) UpdateUser() error {
	var updateUserStmt *sql.Stmt
	updateUserStmt, err := db.Db.Prepare("UPDATE USERS SET `Username`=?, `Email`=?, `PswHash`=?, `Active`=?, `VerHash`=?, `Timeout`=? WHERE `id`=?;")
	if err != nil {
		fmt.Println("error preparring statement to update user in Db with Update, err:", err)
		return err
	}
	defer updateUserStmt.Close()
	var result sql.Result

	result, err = updateUserStmt.Exec(u.Username, u.Email, u.PswHash, u.Active, u.VerHash, u.ID)

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

func MustDispEmailDom() (dispEmailsDomains []string) {
	file, err := os.Open("../disposable_email_blocklist.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dispEmailsDomains = append(dispEmailsDomains, scanner.Text())
	}
	return dispEmailsDomains
}
