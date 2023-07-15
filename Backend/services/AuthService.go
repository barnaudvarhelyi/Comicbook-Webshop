package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	db "main/database"
	"main/dtos"
	m "main/models"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func Register(w http.ResponseWriter, r *http.Request) {

	var user dtos.RegisterDto
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		SendResponse(w, 400, dtos.ResponseDto{Message: "Invalid request body"})
		return
	}

	err = user.UsernameExists()
	if err != nil {
		SendResponse(w, 400, dtos.ResponseDto{Message: err.Error()})
		return
	}
	err = user.EmailExists()
	if err != nil {
		SendResponse(w, 400, dtos.ResponseDto{Message: err.Error()})
		return
	}

	err = user.ValidateUsername()
	if err != nil {
		SendResponse(w, 400, dtos.ResponseDto{Message: err.Error()})
		return
	}
	err = user.ValidatePassword()
	if err != nil {
		SendResponse(w, 400, dtos.ResponseDto{Message: err.Error()})
		return
	}
	var statusCode int
	statusCode, err = user.ValidateEmail()
	if err != nil {
		SendResponse(w, statusCode, dtos.ResponseDto{Message: err.Error()})
		return
	}

	createdAt := time.Now().Local()
	rand.Seed(time.Now().UnixNano())

	timeout := time.Now().Local().AddDate(0, 0, 2)

	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		SendResponse(w, 500, dtos.ResponseDto{Message: err.Error()})
		return
	}

	emailVerPassword, emailVerPWhash, err := user.GenerateEmailVerPswAndHash()
	if err != nil {
		SendResponse(w, 500, dtos.ResponseDto{Message: err.Error()})
		return
	}

	verHash := string(emailVerPWhash)

	err = insertIntoDb(user.Username, user.Email, string(hash), verHash, createdAt, timeout)

	if err != nil {
		SendResponse(w, 500, dtos.ResponseDto{Message: err.Error()})
		return
	}

	domName := "http://localhost:8080"
	subject := "Email Verification"
	HTMLbody :=
		`<html>
			<h1>Click Link to Verify Email</h1>
			<a href="` + domName + `/api/emailver/` + user.Username + `/` + emailVerPassword + `">Click to verify email</a>
		</html>`

	err = SendEmail(user.Email, subject, HTMLbody)

	if err != nil {
		SendResponse(w, 500, dtos.ResponseDto{Message: err.Error()})
		return
	}
	SendResponse(w, 200, dtos.ResponseDto{Message: "200"})
	return
}

func Login(w http.ResponseWriter, r *http.Request) {
	var login dtos.LoginDto
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil || login.Username == "" || login.Password == "" {
		SendResponse(w, 400, dtos.ResponseDto{Message: "Invalid request body"})
		return
	}
	var userId int
	var hash string
	var active bool
	stmt := "SELECT `id`, `pswHash`, `active` FROM users WHERE `username` = ?"
	row := db.Db.QueryRow(stmt, login.Username)

	err = row.Scan(&userId, &hash, &active)
	if err != nil {
		SendResponse(w, 400, dtos.ResponseDto{Message: "Invalid username or password!"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(login.Password))

	if err != nil {
		SendResponse(w, 400, dtos.ResponseDto{Message: "Invalid username or password!"})
		return
	}
	if !active {
		SendResponse(w, 400, dtos.ResponseDto{Message: "User email not verified yet!"})
		return
	}
	session, _ := store.Get(r, "session")
	session.Values["userID"] = userId
	session.Save(r, w)
	SendResponse(w, 200, dtos.ResponseDto{Message: "Successfully loged in!"})
	return
}

func EmailVerHandler(w http.ResponseWriter, r *http.Request) {
	var u m.User
	var linkVerPass string

	vars := mux.Vars(r)

	u.Username, _ = vars["username"]
	linkVerPass, _ = vars["verPass"]

	err := u.GetUserByUsername()
	if err != nil {
		fmt.Println("error selecting verHash in DB by username, err: ", err)
		SendResponse(w, 400, dtos.ResponseDto{Message: "Please try link in verification email again"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.VerHash), []byte(linkVerPass))
	if err == nil {
		err = u.MakeActive()
		if err != nil {
			SendResponse(w, 400, dtos.ResponseDto{Message: "Please try email confirmation link again"})
			return
		}
		session, _ := store.Get(r, "session")
		session.Values["userId"] = u.ID
		session.Save(r, w)
		SendResponse(w, 200, dtos.ResponseDto{Message: "Account activated!"})
		return
	}
	SendResponse(w, 401, dtos.ResponseDto{Message: "Unauthorized"})
}

func AuthMiddleware(hf http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessions, _ := store.Get(r, "session")
		_, ok := sessions.Values["userId"]
		if !ok {
			SendResponse(w, 401, dtos.ResponseDto{Message: "Authorization required!"})
			return
		}
		hf.ServeHTTP(w, r)
	}
}

func insertIntoDb(username, email, hash, verHash string, createdAt, timeout time.Time) error {
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
	insertStmt, err = tx.Prepare("INSERT INTO users (`Username`, `Email`, `PswHash`,`CreatedAt`, `Active`) VALUES (?, ?, ?, ?, ?)")

	if err != nil {
		fmt.Println("error preparing statement: ", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr: ", rollbackErr)
			return err
		}
	}
	defer insertStmt.Close()

	var result sql.Result
	result, err = insertStmt.Exec(username, email, hash, createdAt, 0)

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
	insertStmt2, err = tx.Prepare("INSERT INTO user_email_ver_hash (`Username`, `VerHash`, `Timeout`) VALUES (?, ?, ?)")

	if err != nil {
		fmt.Println("error preparing statement: ", err)
		if rollbackErr := tx2.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr: ", rollbackErr)
			return err
		}
	}
	defer insertStmt2.Close()

	var result2 sql.Result
	result2, err = insertStmt2.Exec(username, verHash, timeout)

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
