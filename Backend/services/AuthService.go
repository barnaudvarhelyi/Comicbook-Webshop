package services

import (
	"encoding/json"
	"main/dtos"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user dtos.RegisterDto
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		SendResponse(w, 400, dtos.ErrorDto{Message: "Invalid request body"})
	}
	//TODO cont.
	// err = user.
}

func Login(w http.ResponseWriter, r *http.Request) {
}
