package services

import (
	"database/sql"
	"fmt"
	"main/dtos"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAllComicbooks(w http.ResponseWriter, r *http.Request) {
	row, err := db.Query("SELECT `id`, `name`, `img`, `desc`, `publisher` FROM `volumes`")
	defer row.Close()
	if err != nil {
		fmt.Println("Error at getting all comikbooks, error: ", err)
		SendResponse(w, 500, dtos.ErrorResponseDto{Error: "Error at getting all comikbooks, error: " + err.Error()})
		return
	}

	var volumes []dtos.Volume
	for row.Next() {
		var v dtos.Volume
		err = row.Scan(&v.ID, &v.Name, &v.Image, &v.Description, &v.Publisher)
		if err != nil {
			fmt.Println("Error at scanning, error: ", err)
			SendResponse(w, 500, dtos.ErrorResponseDto{Error: "Error at getting all comikbooks, error: " + err.Error()})
		}
		volumes = append(volumes, v)
	}
	SendResponse(w, 200, volumes)
}

func GetComicbookById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	volumeId := vars["id"]

	var volume dtos.Volume

	err := db.QueryRow("SELECT `id`, `name`, `img`, `desc`, `publisher` FROM `volumes` WHERE `id` = ?", volumeId).Scan(&volume.ID, &volume.Name, &volume.Image, &volume.Description, &volume.Publisher)

	if err != nil {
		if err == sql.ErrNoRows {
			SendResponse(w, 404, dtos.ErrorResponseDto{Error: "Comicbook with id " + volumeId + " not found!"})
		} else {
			SendResponse(w, 500, dtos.ErrorResponseDto{Error: "Server error: " + err.Error()})
		}
		return
	}

	var issues []dtos.Issue
	row, err := db.Query("SELECT `id`, `name`, `issue_number`, `img`, `cover_date`, `date_added` FROM `issues` WHERE volume_id = ?", volumeId)
	defer row.Close()
	if err != nil {
		fmt.Println("Error at getting comikbook's issues, error: ", err)
		SendResponse(w, 500, dtos.ErrorResponseDto{Error: "Error at getting comikbook's issues, error: " + err.Error()})
		return
	}
	for row.Next() {
		var i dtos.Issue
		err = row.Scan(&i.ID, &i.Name, &i.IssueNumebr, &i.Image, &i.CoverDate, &i.DateAdded)
		if err != nil {
			fmt.Println("Error at scanning, error: ", err)
			SendResponse(w, 500, dtos.ErrorResponseDto{Error: "Error at getting all comikbooks, error: " + err.Error()})
		}
		issues = append(issues, i)
	}
	SendResponse(w, 200, dtos.VolumeById{Volume: volume, IssuesCount: len(issues), Issues: issues})
}

func GetIssueById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	issueId := vars["id"]
	var i dtos.Issue
	err := db.QueryRow("SELECT `id`, `name`, `issue_number`, `img`, `cover_date`, `date_added` FROM `issues` WHERE `id` = ?", issueId).Scan(&i.ID, &i.Name, &i.IssueNumebr, &i.Image, &i.CoverDate, &i.DateAdded)

	if err != nil {
		if err == sql.ErrNoRows {
			SendResponse(w, 404, dtos.ErrorResponseDto{Error: "Issue with id " + issueId + " not found!"})
		} else {
			SendResponse(w, 500, dtos.ErrorResponseDto{Error: "Server error: " + err.Error()})
		}
		return
	}
	SendResponse(w, 200, i)
}
