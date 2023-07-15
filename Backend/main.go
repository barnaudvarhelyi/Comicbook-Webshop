package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"main/controllers"
	db "main/database"
)

func init() {
	db.InitDb()
}

func main() {

	mux := mux.NewRouter()

	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "X-Content-Type-Options"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	mux.HandleFunc("/api/test", controllers.TestController).Methods("GET")

	mux.HandleFunc("/api/register", controllers.RegisterController).Methods("POST")

	mux.HandleFunc("/api/emailver/{username}/{verPass}", controllers.VerifyEmailController)

	mux.HandleFunc("/api/login", controllers.LoginController).Methods("POST")

	fmt.Println("--------------- App has been started! ---------------")
	err := http.ListenAndServe(":8080", handlers.CORS(header, methods, origins)(mux))
	if err != nil {
		log.Fatal(err)
	}
}
