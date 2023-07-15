package services

import (
	"encoding/json"
	"fmt"
	"main/dtos"
	m "main/models"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user dtos.RegisterDto
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		SendResponse(w, 400, dtos.ErrorDto{Message: "Invalid request body"})
		return
	}

	err = user.ValidateUsername()
	if err != nil {
		SendResponse(w, 400, dtos.ErrorDto{Message: err.Error()})
		return
	}
	err = user.ValidatePassword()
	if err != nil {
		SendResponse(w, 400, dtos.ErrorDto{Message: err.Error()})
		return
	}
	var statusCode int
	statusCode, err = user.ValidateEmail()
	if err != nil {
		SendResponse(w, statusCode, dtos.ErrorDto{Message: err.Error()})
		return
	}
	err = user.UsernameExists()
	if err != nil {
		SendResponse(w, 400, dtos.ErrorDto{Message: err.Error()})
		return
	}
	err = user.EmailExists()
	if err != nil {
		SendResponse(w, 400, dtos.ErrorDto{Message: err.Error()})
		return
	}

	createdAt := time.Now().Local()
	rand.Seed(time.Now().UnixNano())

	timeout := time.Now().Local().AddDate(0, 0, 2)

	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		SendResponse(w, 500, dtos.ErrorDto{Message: err.Error()})
		return
	}

	emailVerPassword, emailVerPWhash, err := user.GenerateEmailVerPswAndHash()
	if err != nil {
		SendResponse(w, 500, dtos.ErrorDto{Message: err.Error()})
		return
	}

	verHash := string(emailVerPWhash)

	err = user.InsertIntoDb(user.Username, user.Email, string(hash), verHash, createdAt, timeout)

	if err != nil {
		SendResponse(w, 500, dtos.ErrorDto{Message: err.Error()})
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
		SendResponse(w, 500, dtos.ErrorDto{Message: err.Error()})
		return
	}
	SendResponse(w, 200, "200")
	return
}

func Login(w http.ResponseWriter, r *http.Request) {
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
		SendResponse(w, 400, dtos.ErrorDto{Message: "Please try link in verification email again"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.VerHash), []byte(linkVerPass))
	if err == nil {
		err = u.MakeActive()
		if err != nil {
			SendResponse(w, 400, dtos.ErrorDto{Message: "Please try email confirmation link again"})
			return
		}
		// session, _ := store.Get(r, "session")
		// session.Values["userId"] = u.ID
		// session.Save(r, w)
		fmt.Println("After register userID: ", u.ID)
		SendResponse(w, 200, "Account activated!")
		return
	}
	SendResponse(w, 401, dtos.ErrorDto{Message: "Unauthorized"})
}
