package main

import (
	"main/controllers"
	db "main/database"
	"main/services"
)

func init() {
	db.InitDb()
	services.Init()
}

func main() {
	controllers.InitServer()
}
