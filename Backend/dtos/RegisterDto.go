package dtos

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
	"unicode"

	db "main/database"

	emailverifier "github.com/AfterShip/email-verifier"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

type RegisterDto struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	ConfPassword string `json:"conf_password"`
	Email        string `json:"email"`
}

var emailVerPassword string

func (u *RegisterDto) ValidateUsername() error {
	for _, char := range u.Username {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return errors.New("Only alphanumeric characters allowed for username")
		}
	}
	if 5 <= len(u.Username) && len(u.Username) <= 50 {
		return nil
	}
	return errors.New("Username lenght must be greater than 4 and less than 50 characters")
}

func (u *RegisterDto) ValidatePassword() error {
	err := passwordvalidator.Validate(u.Password, 60)
	return err
}

func (u *RegisterDto) ValidateEmail() (statusCode int, err error) {

	var Verifier = emailverifier.NewVerifier()

	Verifier = Verifier.EnableDomainSuggest()
	Verifier = Verifier.EnableSMTPCheck()
	dispEmailsDomains := mustDispEmailDom()
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

func (u *RegisterDto) UsernameExists() (exists error) {
	stmt := "SELECT `id` FROM users WHERE `username` = ?"
	row := db.Db.QueryRow(stmt, u.Username)

	var id string
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return nil
	}
	return errors.New("Username already exists!")
}

func (u *RegisterDto) EmailExists() (exists error) {
	stmt := "SELECT `id` FROM users WHERE `email` = ?"
	row := db.Db.QueryRow(stmt, u.Email)

	var id string
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return nil
	}
	return errors.New("Email already exists!")
}

func mustDispEmailDom() (dispEmailsDomains []string) {
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

func (u *RegisterDto) GenerateEmailVerPswAndHash() (string, []byte, error) {
	var err error
	var alphaNumRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQESTUVWXYZ1234567890")
	emailVerRandRune := make([]rune, 64)
	for i := 0; i < 64; i++ {
		emailVerRandRune[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes)-1)]
	}
	emailVerPassword = string(emailVerRandRune)

	var emailVerPswHash []byte
	emailVerPswHash, err = bcrypt.GenerateFromPassword([]byte(emailVerPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}

	return emailVerPassword, emailVerPswHash, nil
}

func (u *RegisterDto) InsertIntoDb(username, email, hash, verHash string, createdAt, timeout time.Time) error {
	tx, err := db.Db.Begin()
	if err != nil {
		fmt.Println("failed to begin transaction, err", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr: ", rollbackErr)
		}
		return err
	}
	defer tx.Rollback()

	var insertStmt *sql.Stmt
	insertStmt, err = tx.Prepare("INSERT INTO users (`Username`, `Email`, `PswHash`,`CreatedAt`, `Active`, `VerHash`) VALUES (?, ?, ?, ?, ?, ?)")

	if err != nil {
		fmt.Println("error preparing statement: ", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr: ", rollbackErr)
			return err
		}
	}
	defer insertStmt.Close()

	var result sql.Result
	result, err = insertStmt.Exec(username, email, hash, createdAt, 0, verHash)

	aff, err := result.RowsAffected()
	if aff == 0 {
		fmt.Println("error at inserting: ", err)
		return err
	}

	var tx2 *sql.Tx
	tx2, err = db.Db.Begin()
	if err != nil {
		fmt.Println("failed to begin transaction, err", err)
		if rollbackErr := tx2.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr: ", rollbackErr)
		}
		return err
	}
	defer tx2.Rollback()

	var insertStmt2 *sql.Stmt
	insertStmt2, err = tx.Prepare("INSERT INTO user_email_ver_hash (`Username`, `Email`, `VerHash`, `Timeout`) VALUES (?, ?, ?, ?)")

	if err != nil {
		fmt.Println("error preparing statement: ", err)
		if rollbackErr := tx2.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr: ", rollbackErr)
			return err
		}
	}
	defer insertStmt2.Close()

	var result2 sql.Result
	result2, err = insertStmt2.Exec(username, email, verHash, timeout)

	aff, err = result2.RowsAffected()
	if aff == 0 {
		fmt.Println("Error at inserting: ", err)
		return err
	}

	if err != nil {
		if rollbackErr := tx2.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr: ", rollbackErr)
			return err
		}
		return err
	}
	if commitErr := tx.Commit(); commitErr != nil {
		fmt.Println("error commiting changes, err: ", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr: ", rollbackErr)
		}
	}
	if commitErr := tx2.Commit(); commitErr != nil {
		fmt.Println("error commiting changes, err: ", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr: ", rollbackErr)
		}
	}
	return nil
}
