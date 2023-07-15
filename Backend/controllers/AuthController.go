package controllers

import (
	"main/services"
	"net/http"
)

func RegisterController(w http.ResponseWriter, r *http.Request) {
	services.Register(w, r)
}

func LoginController(w http.ResponseWriter, r *http.Request) {
	services.Login(w, r)
}

func VerifyEmailController(w http.ResponseWriter, r *http.Request) {
	services.EmailVerHandler(w, r)
}
