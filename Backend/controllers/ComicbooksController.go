package controllers

import (
	"main/services"
	"net/http"
)

func GetAllComicbooksController(w http.ResponseWriter, r *http.Request) {
	services.GetAllComicbooks(w, r)
}

func GetComicbookByIdController(w http.ResponseWriter, r *http.Request) {
	services.GetComicbookById(w, r)
}

func GetIssueByIdController(w http.ResponseWriter, r *http.Request) {
	services.GetIssueById(w, r)
}
