package controllers

import (
	"main/services"
	"net/http"
)

func TestController(w http.ResponseWriter, r *http.Request) {
	services.GetAllUserTest(w, r)
}
