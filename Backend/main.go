package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"main/database"
)

func init() {
	database.InitDb()
}

func main() {
	mux := mux.NewRouter()

	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "X-Content-Type-Options"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	err := http.ListenAndServe(":8080", handlers.CORS(header, methods, origins)(mux))
	if err != nil {
		log.Fatal(err)
	}
}
