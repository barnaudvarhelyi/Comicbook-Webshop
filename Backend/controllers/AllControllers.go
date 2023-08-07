package controllers

import (
	"fmt"
	"log"
	"main/services"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func InitServer() {
	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "X-Content-Type-Options"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	mux := mux.NewRouter()
	mux.HandleFunc("/api/test", services.AuthMiddleware(TestController)).Methods("GET")

	mux.HandleFunc("/api/register", RegisterController).Methods("POST")
	mux.HandleFunc("/api/login", LoginController).Methods("POST")
	mux.HandleFunc("/api/emailver/{username}/{verPass}", VerifyEmailController)

	mux.HandleFunc("/api/comicbooks/all", GetAllComicbooksController).Methods("GET")
	mux.HandleFunc("/api/comicbook/{id}", GetComicbookByIdController).Methods("GET")

	mux.HandleFunc("/api/issue/{id}", GetIssueByIdController).Methods("GET")

	fmt.Println("--------------- App has been started! ---------------")
	err := http.ListenAndServe(":8080", handlers.CORS(header, methods, origins)(mux))
	if err != nil {
		log.Fatal(err)
	}
}
